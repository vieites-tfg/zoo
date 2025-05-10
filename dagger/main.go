package main

import (
	"context"
	"dagger/dagger/internal/dagger"
	"fmt"
	"slices"
)

type Dagger struct {
	// +optional
	Sec *dagger.Secret
}

var secrets Secrets

// Builds the base image from the Dockerfile.
func (m *Dagger) Base(ctx context.Context, src *dagger.Directory) (*dagger.Container, error) {
	if m.Sec != (&dagger.Secret{}) {
		contents, err := src.File(".env").Contents(ctx)
		if err != nil {
			return nil, err
		}

		m.Sec = dagger.Connect().SetSecret("secrets", contents)
	}

	ctr := dag.
		Container().
		From("node:20").
		WithWorkdir("/app").
		WithFile("package.json", src.File("package.json")).
		WithFile("lerna.json", src.File("lerna.json")).
		WithFile("yarn.lock", src.File("yarn.lock")).
		WithDirectory("packages", src.Directory("packages")).
		WithExec([]string{"yarn", "install"}).
		WithExec([]string{"yarn", "global", "add", "lerna@8.2.1"}).
		WithExec([]string{"yarn", "global", "add", "@vercel/ncc"})

	return ctr, nil
}

// Init configures the content with the .env environment variables
func (m *Dagger) Init(
	ctx context.Context,
	// +defaultPath="/"
	src *dagger.Directory,
) (*dagger.Container, error) {
	ctr, err := m.Base(ctx, src)
	if err != nil {
		return nil, err
	}

	content, err := m.Sec.Plaintext(ctx)
	if err != nil {
		return nil, err
	}

	vars := ParseEnvFile(content)
	err = MakeSecrets(ctx, vars)
	if err != nil {
		return nil, err
	}

	for key, value := range secrets {
		ctr = ctr.WithSecretVariable(key, value)
	}

	return ctr, nil
}

// Returns the backend container as a service with the Mongo database.
func (m *Dagger) Backend(
	ctx context.Context,
	// +defaultPath="/"
	src *dagger.Directory,
) (*dagger.Service, error) {
	_, err := m.Init(ctx, src)
	if err != nil {
		return nil, err
	}

	mongoPort, err := GetMongoPort(ctx)
	if err != nil {
		return nil, err
	}

	mongoInit := src.Directory("mongo-init")
	mongo := dag.
		Container().
		From("mongo:7.0").
		WithSecretVariable("MONGO_INITDB_DATABASE", secrets["MONGO_DATABASE"]).
		WithSecretVariable("MONGO_INITDB_ROOT_USERNAME", secrets["MONGO_ROOT"]).
		WithSecretVariable("MONGO_INITDB_ROOT_PASSWORD", secrets["MONGO_ROOT_PASS"]).
		WithExposedPort(mongoPort).
		WithMountedDirectory("/docker-entrypoint-initdb.d", mongoInit).
		AsService().
		WithHostname("mongodb")

	_, err = mongo.Start(ctx)
	if err != nil {
		return nil, err
	}

	build, err := m.BackendBuild(ctx, src)
	if err != nil {
		return nil, err
	}

	compiled := build.File("/app/dist/index.js")
	pkgJson := build.File("/app/packages/backend/package.json")

	back := dag.
		Container().From("node:20-alpine").
		WithExposedPort(3000).
		WithEnvVariable("NODE_ENV", "production").
		WithEnvVariable("YARN_CACHE_FOLDER", "/.yarn/cache").
		WithMountedCache("/.yarn/cache", dag.CacheVolume("yarn-cache")).
		WithSecretVariable("MONGODB_URI", secrets["MONGODB_URI"]).
		WithWorkdir("/app").
		WithFile("/app/package.json", pkgJson).
		WithFile("/app/index.js", compiled).
		WithExec([]string{"yarn", "install", "--production"}).
		WithEntrypoint([]string{"node", "index.js"}).
		AsService().
		WithHostname("zoo-bakend")

	svc, err := back.Start(ctx)
	if err != nil {
		return nil, err
	}

	return svc, nil
}

// Returns the backend container.
func (m *Dagger) BackendBuild(
	ctx context.Context,
	// +defaultPath="/"
	src *dagger.Directory,
) (*dagger.Container, error) {
	base, err := m.Init(ctx, src)
	if err != nil {
		return nil, err
	}

	build := base.
		WithWorkdir("/app").
		WithExec([]string{"lerna", "run", "--scope", "@vieites-tfg/zoo-backend", "build"}).
		WithExec([]string{"ncc", "build", "./packages/backend/dist/index.js", "-o", "./dist"})

	return build, nil
}

// Returns the frontend container.
func (m *Dagger) FrontendBuild(
	ctx context.Context,
	// +defaultPath="/"
	src *dagger.Directory,
) (*dagger.Container, error) {
	base, err := m.Init(ctx, src)
	if err != nil {
		return nil, err
	}

	front := base.
		WithWorkdir("/app").
		WithExec([]string{"lerna", "run", "--scope", "@vieites-tfg/zoo-frontend", "build"})

	return front, nil
}

// Returns the frontend container as a service.
func (m *Dagger) Frontend(
	ctx context.Context,
	// +defaultPath="/"
	src *dagger.Directory,
) (*dagger.Service, error) {
	front, err := m.FrontendBuild(ctx, src)
	if err != nil {
		return nil, err
	}

	return front.AsService().WithHostname("zoo-frontend"), nil
}

func (m *Dagger) Lint(
	ctx context.Context,
	// +defaultPath="/"
	src *dagger.Directory,
) (string, error) {
	ctr, err := m.Init(ctx, src)
	if err != nil {
		return "", err
	}

	return ctr.WithMountedDirectory("/app", src).
		WithWorkdir("/app").
		WithExec([]string{"yarn", "lint"}).
		Stdout(ctx)
}

// Run the tests for the backend.
func (m *Dagger) TestBackend(
	ctx context.Context,
	// +defaultPath="/"
	src *dagger.Directory,
) (string, error) {
	ctr, err := m.Init(ctx, src)
	if err != nil {
		return "", err
	}

	return ctr.
		WithExec([]string{"lerna", "run", "test", "--scope", "@vieites-tfg/zoo-backend"}).
		Stdout(ctx)
}

// Returns the Cypress container used to run the e2e tests for the frontend.
func (m *Dagger) Cypress(src *dagger.Directory) *dagger.Container {
	return dag.
		Container().
		From("cypress/browsers").
		WithMountedDirectory("/e2e", src).
		WithWorkdir("/e2e").
		WithExec([]string{"npx", "cypress", "install"}).
		WithExec([]string{"yarn", "add", "lerna@8.2.1", "-W"})
}

// Run the tests for the frontend.
func (m *Dagger) TestFrontend(
	ctx context.Context,
	// +defaultPath="/"
	src *dagger.Directory,
	front *dagger.Service,
) (string, error) {
	_, err := m.Init(ctx, src)
	if err != nil {
		return "", err
	}

	return m.Cypress(src).
		WithServiceBinding("zoo-frontend", front).
		WithExec([]string{"yarn", "run", "e2e"}).
		Stdout(ctx)
}

// Builds the package Docker image indicated and returns its container.
// package is optional can be one of: ["backend", "frontend"].
// If package is empty, both will be built.
func (m *Dagger) BuildImage(
	ctx context.Context,
	// +defaultPath="/"
	src *dagger.Directory,
	// +default="all"
	pkg string,
) ([]*dagger.Container, error) {
	pkgs := []string{"backend", "frontend"}
	if slices.Contains(pkgs, pkg) {
		pkgs = []string{pkg}
	} else if pkg != "all" {
		return nil, fmt.Errorf("Not a valid option: '%s' is not one of ['backend', 'frontend'].", pkg)
	}

	var (
		containers []*dagger.Container
		ctr        *dagger.Container
		err        error
	)
	for _, p := range pkgs {
		if p == "backend" {
			ctr, err = m.BackendBuild(ctx, src)
		} else if p == "frontend" {
			ctr, err = m.FrontendBuild(ctx, src)
		}
		if err != nil {
			return nil, err
		}
		containers = append(containers, ctr)
	}

	return containers, nil
}
