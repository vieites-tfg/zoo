package main

import (
	"context"
	"dagger/dagger/internal/dagger"
)

type Frontend struct {
	// The name of the package.
	Name      string

	// The base container from which the functions will be run.
	Base      *dagger.Container

	// The secrets needed to launch the package.
	Secrets   SecMap

	// The main object.
	Ci *Dagger
}

// Builds the frontend package, generating only one executable file and returns the container.
func (m *Frontend) Build(ctx context.Context) *dagger.Container {
	build := m.Base.
		WithWorkdir("/app").
		WithExec([]string{"lerna", "run", "--scope", "@vieites-tfg/zoo-frontend", "build"})

	return build
}

// Based on the build stage, gets the executable file and creates a ready to run container with Nginx and the port 80 exported.
func (m *Frontend) Ctr(ctx context.Context) *dagger.Container {
	build := m.Build(ctx)

	dist := build.Directory("/app/packages/frontend/dist")

	ctr := dag.
		Container().
		From("nginx:alpine").
		WithWorkdir("/usr/share/nginx/html").
		WithDirectory(".", dist).
		WithExposedPort(80).
		WithEntrypoint([]string{"nginx", "-g", "daemon off;"})

	return ctr
}

// Returns the ready-to-use container as a service.
func (m *Frontend) Service(ctx context.Context) *dagger.Service {
	return m.Ctr(ctx).AsService().WithHostname("zoo-frontend")
}

// Runs the linter for the package.
func (m *Frontend) Lint(ctx context.Context) (string, error) {
	return Lint(ctx, m.Base, m.Name)
}

// Run the e2e test. This requires to have both the backend and frontend services up and running correctly. You have to pass the frontend service as an argument to the function.
func (m *Frontend) Test(
	ctx context.Context,
	// +defaultPath="/"
	src *dagger.Directory,
	front *dagger.Service,
) (string, error) {
	return Cypress(src).
		WithServiceBinding("zoo-frontend", front).
		WithExec([]string{"yarn", "run", "e2e"}).
		Stdout(ctx)
}

// Publish the Docker image of the package with the "latest" and the npm package (inside the 'package.json') versions.
func (m *Frontend) PublishImage(
	ctx context.Context,
	// +defaultPath="/"
	src *dagger.Directory,
) ([]string, error) {
	var err error

	_, err = m.Lint(ctx)
	if err != nil {
		return []string{}, err
	}

	_, err = m.Ci.Endtoend(ctx, src)
	if err != nil {
		return []string{}, err
	}

	return PublishImage(ctx, m.Base, m.Ctr(ctx), m.Name, m.Secrets.Get("CR_PAT"))
}

// Publish the npm package.
func (m *Frontend) PublishPkg(
	ctx context.Context,
	// +defaultPath="/"
	src *dagger.Directory,
) (string, error) {
	var err error

	_, err = m.Lint(ctx)
	if err != nil {
		return "", err
	}

	_, err = m.Ci.Endtoend(ctx, src)
	if err != nil {
		return "", err
	}

	return PublishPkg(ctx, m.Base, m.Name, m.Secrets.Get("CR_PAT"))
}
