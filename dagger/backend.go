package main

import (
	"context"
	"dagger/dagger/internal/dagger"
	"fmt"
	"strconv"
)

type Backend struct {
	// The name of the package.
	Name string

	// The base container from which the functions will be run.
	Base *dagger.Container

	// The secrets needed to launch the package.
	Secrets SecMap
}

// Builds the backend package, generating only one executable file and returns the container.
func (m *Backend) Build(ctx context.Context) *dagger.Container {
	build := m.Base.
		WithWorkdir("/app").
		WithExec([]string{"lerna", "run", "--scope", "@vieites-tfg/zoo-backend", "build"}).
		WithExec([]string{"ncc", "build", "./packages/backend/dist/index.js", "-o", "./dist"})

	return build
}

// Based on the build stage, gets the executable file and creates a ready to run container. Since it needs a Mongo database, it does not work by itself as a service.
func (m *Backend) Ctr(ctx context.Context) *dagger.Container {
	build := m.Build(ctx)

	compiled := build.File("/app/dist/index.js")
	pkgJson := build.File("/app/packages/backend/package.json")

	dag := dagger.Connect()
	back := dag.
		Container().From("node:20-alpine").
		WithExposedPort(3000).
		WithEnvVariable("NODE_ENV", "production").
		WithEnvVariable("YARN_CACHE_FOLDER", "/.yarn/cache").
		WithMountedCache("/.yarn/cache", dag.CacheVolume("yarn-cache")).
		WithWorkdir("/app").
		WithFile("/app/package.json", pkgJson).
		WithFile("/app/index.js", compiled).
		WithExec([]string{"yarn", "install", "--production"}).
		WithEntrypoint([]string{"node", "index.js"})

	return back
}

// Creates a Mongo database and, based on the ready-to-run container, binds it to the backend using the environment variable. It returns the backend service with the 3000 port exported.
func (m *Backend) Service(
	ctx context.Context,
	// +defaultPath="/"
	src *dagger.Directory,
) (*dagger.Service, error) {
	mongoPort, err := getMongoPort(ctx, m.Secrets.Get("MONGO_PORT"))
	if err != nil {
		return nil, err
	}

	mongoInit := src.Directory("mongo-init")
	mongo := dagger.Connect().
		Container().
		From("mongo:7.0").
		WithSecretVariable("MONGO_INITDB_DATABASE", m.Secrets.Get("MONGO_DATABASE")).
		WithSecretVariable("MONGO_INITDB_ROOT_USERNAME", m.Secrets.Get("MONGO_ROOT")).
		WithSecretVariable("MONGO_INITDB_ROOT_PASSWORD", m.Secrets.Get("MONGO_ROOT_PASS")).
		WithExposedPort(mongoPort).
		WithMountedDirectory("/docker-entrypoint-initdb.d", mongoInit).
		AsService().
		WithHostname("mongodb")

	_, err = mongo.Start(ctx)
	if err != nil {
		return nil, err
	}

	mongoUri, err := createMongoUri(ctx, m.Secrets)

	back := m.Ctr(ctx).
		WithSecretVariable("MONGODB_URI", mongoUri).
		AsService().WithHostname("zoo-bakend")

	svc, err := back.Start(ctx)
	if err != nil {
		return nil, err
	}

	return svc, nil
}

// Run the test for the package.
func (m *Backend) Test(ctx context.Context) (string, error) {
	return m.Base.
		WithExec([]string{"lerna", "run", "test", "--scope", "@vieites-tfg/zoo-backend"}).
		Stdout(ctx)
}

// Runs the linter for the package.
func (m *Backend) Lint(ctx context.Context) (string, error) {
	return Lint(ctx, m.Base, m.Name)
}

// Publish the Docker image of the package with the "latest" and the npm package (inside the 'package.json') versions.
func (m *Backend) PublishImage(ctx context.Context) ([]string, error) {
	return PublishImage(ctx, m.Base, m.Ctr(ctx), m.Name, m.Secrets.Get("CR_PAT"))
}

func getMongoPort(ctx context.Context, port *dagger.Secret) (int, error) {
	mongo_portStr, err := port.Plaintext(ctx)
	if err != nil {
		return 0, err
	}

	mongo_port, err := strconv.Atoi(mongo_portStr)
	if err != nil {
		return 0, err
	}

	return mongo_port, nil
}

// Publish the npm package.
func (m *Backend) PublishPkg(
	ctx context.Context,
	// +defaultPath="/"
	src *dagger.Directory,
) (string, error) {
	var err error

	_, err = m.Lint(ctx)
	if err != nil {
		return "", err
	}

	_, err = m.Test(ctx)
	if err != nil {
		return "", err
	}

	return PublishPkg(ctx, m.Base, m.Name, m.Secrets.Get("CR_PAT"))
}

func createMongoUri(ctx context.Context, secrets SecMap) (*dagger.Secret, error) {
	var (
		err       error
		root      string
		rootPass  string
		mongoPort string
		db        string
	)

	root, err = secrets.Get("MONGO_ROOT").Plaintext(ctx)
	if err != nil {
		return nil, err
	}

	rootPass, err = secrets.Get("MONGO_ROOT_PASS").Plaintext(ctx)
	if err != nil {
		return nil, err
	}

	mongoPort, err = secrets.Get("MONGO_PORT").Plaintext(ctx)
	if err != nil {
		return nil, err
	}

	db, err = secrets.Get("MONGO_DATABASE").Plaintext(ctx)
	if err != nil {
		return nil, err
	}

	mongoUri := fmt.Sprintf("mongodb://%s:%s@mongodb:%s/%s?authSource=admin",
		root,
		rootPass,
		mongoPort,
		db,
	)

	return dagger.Connect().SetSecret("mongoUri", mongoUri), nil
}
