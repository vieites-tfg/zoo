package main

import (
	"context"
	"dagger/dagger/internal/dagger"
)

type Dagger struct {
	// This is the '.env' file with the environment variables needed to launch the application.
	// +required
	SecEnv  *dagger.Secret

	secrets secrets
}

func New(secEnv *dagger.Secret) *Dagger {
	return &Dagger{
		SecEnv: secEnv,
	}
}

// Builds the base image from the Dockerfile.
func (m *Dagger) Base(ctx context.Context, src *dagger.Directory) (*dagger.Container, error) {
	ctr := dag.
		Container().
		From("node:20").
		WithWorkdir("/app").
		WithFile("package.json", src.File("package.json")).
		WithFile("lerna.json", src.File("lerna.json")).
		WithFile("yarn.lock", src.File("yarn.lock")).
		WithDirectory("packages", src.Directory("packages")).
		WithDirectory(".git", src.Directory(".git")).
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

	content, err := m.SecEnv.Plaintext(ctx)
	if err != nil {
		return nil, err
	}

	m.secrets, err = MakeSecrets(ctx, content)
	if err != nil {
		return nil, err
	}

	return ctr, nil
}

// Functions related to the backend package.
func (m *Dagger) Backend(
	ctx context.Context,
	// +defaultPath="/"
	src *dagger.Directory,
) (*Backend, error) {
	base, err := m.Init(ctx, src)
	if err != nil {
		return nil, err
	}

	keys := []string{"MONGO_PORT", "MONGO_DATABASE", "MONGO_ROOT", "MONGO_ROOT_PASS", "CR_PAT"}
	values := []*dagger.Secret{
		m.secrets["MONGO_PORT"],
		m.secrets["MONGO_DATABASE"],
		m.secrets["MONGO_ROOT"],
		m.secrets["MONGO_ROOT_PASS"],
		m.secrets["CR_PAT"],
	}

	return &Backend{
		Name:    "backend",
		Base:    base,
		Secrets: SecMap{Keys: keys, Values: values},
	}, nil
}

// Functions related to the backend package.
func (m *Dagger) Frontend(
	ctx context.Context,
	// +defaultPath="/"
	src *dagger.Directory,
) (*Frontend, error) {
	base, err := m.Init(ctx, src)
	if err != nil {
		return nil, err
	}

	keys := []string{"CR_PAT"}
	values := []*dagger.Secret{m.secrets["CR_PAT"]}

	return &Frontend{
		Name:    "frontend",
		Base:    base,
		Secrets: SecMap{Keys: keys, Values: values},
		Ci:      m,
	}, nil
}

func (m *Dagger) Endtoend(
	ctx context.Context,
	// +defaultPath="/"
	src *dagger.Directory,
) (string, error) {
	back, err := m.Backend(ctx, src)
	if err != nil {
		return "", err
	}

	front, err := m.Frontend(ctx, src)
	if err != nil {
		return "", err
	}

	// Linter
	_, err = back.Lint(ctx)
	if err != nil {
		return "", err
	}

	_, err = front.Lint(ctx)
	if err != nil {
		return "", err
	}

	// Backend tests
	_, err = back.Test(ctx)
	if err != nil {
		return "", err
	}

	// Frontend tests
	backSvc, err := back.Service(ctx, src)
	if err != nil {
		return "", err
	}

	backendSvc, err := backSvc.Start(ctx)
	if err != nil {
		return "", err
	}

	frontendSvc, err := front.Service(ctx).Start(ctx)
	if err != nil {
		return "", err
	}

	test := Cypress(src).
		WithServiceBinding("zoo-frontend", frontendSvc).
		WithServiceBinding("zoo-backend", backendSvc).
		WithEnvVariable("BASE_URL", "http://zoo-frontend").
		WithExec([]string{"yarn", "run", "e2e"})

	return test.Stdout(ctx)
}
