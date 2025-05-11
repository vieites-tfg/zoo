package main

import (
	"context"
	"dagger/dagger/internal/dagger"
)

type Frontend struct {
	Name      string
	Base      *dagger.Container
	Secrets   SecMap
}

func (m *Frontend) Build(ctx context.Context) *dagger.Container {
	build := m.Base.
		WithWorkdir("/app").
		WithExec([]string{"lerna", "run", "--scope", "@vieites-tfg/zoo-frontend", "build"})

	return build
}

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

func (m *Frontend) Service(ctx context.Context) *dagger.Service {
	return m.Ctr(ctx).AsService().WithHostname("zoo-frontend")
}

func (m *Frontend) Lint(ctx context.Context) (string, error) {
	return Lint(ctx, m.Base, m.Name)
}

// Run the tests for the frontend.
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

func (m *Frontend) PublishImage(ctx context.Context) ([]string, error) {
	return PublishImage(ctx, m.Base, m.Ctr(ctx), m.Name, m.Secrets.Get("CR_PAT"))
}

func Cypress(src *dagger.Directory) *dagger.Container {
	return dagger.Connect().
		Container().
		From("cypress/browsers").
		WithMountedDirectory("/e2e", src).
		WithWorkdir("/e2e").
		WithExec([]string{"npx", "cypress", "install"}).
		WithExec([]string{"yarn", "add", "lerna@8.2.1", "-W"})
}
