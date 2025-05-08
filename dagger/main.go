package main

import (
	"context"
	"dagger/dagger/internal/dagger"
)

type Dagger struct{
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

	return src.DockerBuild(dagger.DirectoryDockerBuildOpts{Target: "base"}).
		WithMountedDirectory("/app", src).
		WithWorkdir("/app"), nil
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

// Returns the backend service with the Mongo database.
func (m *Dagger) Backend(
	ctx context.Context,
	// +defaultPath="/"
	src *dagger.Directory,
) (*dagger.Service, error) {
	ctr, err := m.Init(ctx, src)
	if err != nil {
		return nil, err
	}

	mongoPort, err := GetMongoPort(ctx)
	if err != nil {
		return nil, err
	}

	mongoInit := ctr.Directory("/app/mongo-init")
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

	back := dag.
		Container().
		From("ghcr.io/vieites-tfg/zoo-backend").
		WithExposedPort(3000).
		WithEnvVariable("NODE_ENV", "production").
		WithEnvVariable("YARN_CACHE_FOLDER", "/.yarn/cache").
		WithMountedCache("/.yarn/cache", dag.CacheVolume("yarn-cache")).
		WithSecretVariable("MONGODB_URI", secrets["MONGODB_URI"]).
		AsService().
		WithHostname("zoo-bakend")

	svc, err := back.Start(ctx)
	if err != nil {
		return nil, err
	}

	return svc, nil
}

// Returns the frontend service.
func (m *Dagger) Frontend(
	ctx context.Context,
	// +defaultPath="/"
	src *dagger.Directory,
) (*dagger.Service, error) {
	_, err := m.Init(ctx, src)
	if err != nil {
		return nil, err
	}

	front := dag.
		Container().
		From("ghcr.io/vieites-tfg/zoo-frontend").
		WithExposedPort(80).
		WithEnvVariable("NODE_ENV", "production").
		WithEnvVariable("YARN_CACHE_FOLDER", "/.yarn/cache").
		WithMountedCache("/.yarn/cache", dag.CacheVolume("yarn-cache")).
		AsService().
		WithHostname("zoo-frontend")

	return front, nil
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
