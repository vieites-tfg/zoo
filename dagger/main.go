package main

import (
	"context"
	"dagger/dagger/internal/dagger"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Dagger struct{}

var secrets Secrets

// Builds the base image from the Dockerfile.
func (m *Dagger) Base(ctx context.Context, src *dagger.Directory) *dagger.Container {
	return src.DockerBuild(dagger.DirectoryDockerBuildOpts{Target: "base"}).
		WithMountedDirectory("/app", src).
		WithWorkdir("/app")
}

// Init configures the content with the .env environment variables
func (m *Dagger) Init(
	ctx context.Context,
	// +defaultPath="/"
	src *dagger.Directory,
	sec *dagger.Secret,
) (*dagger.Container, error) {
	ctr := m.Base(ctx, src)

	content, err := sec.Plaintext(ctx)
	if err != nil {
		return nil, err
	}

	vars := parseEnvFile(content)
	err = makeSecrets(ctx, vars)
	if err != nil {
		return nil, err
	}

	for key, value := range secrets {
		ctr = ctr.WithSecretVariable(key, value)
	}

	return ctr, nil
}

// Returns the backend service.
func (m *Dagger) LaunchBackend(
	ctx context.Context,
	// +defaultPath="/"
	src *dagger.Directory,
	sec *dagger.Secret,
) (*dagger.Service, error) {
	ctr, err := m.Init(ctx, src, sec)
	if err != nil {
		return nil, err
	}

	mongoPort, err := getMongoPort(ctx)
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
		WithEnvVariable("YARN_CACHE_FOLDER", ".cache").
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
func (m *Dagger) LaunchFrontend(
	ctx context.Context,
	// +defaultPath="/"
	src *dagger.Directory,
	sec *dagger.Secret,
) (*dagger.Service, error) {
	_, err := m.Init(ctx, src, sec)
	if err != nil {
		return nil, err
	}

	front := dag.
		Container().
		From("ghcr.io/vieites-tfg/zoo-frontend").
		WithExposedPort(80).
		WithEnvVariable("NODE_ENV", "production").
		WithEnvVariable("YARN_CACHE_FOLDER", ".cache").
		AsService().
		WithHostname("zoo-frontend")

	return front, nil
}

func (m *Dagger) Lint(
	ctx context.Context,
	// +defaultPath="/"
	src *dagger.Directory,
	sec *dagger.Secret,
) (string, error) {
	ctr, err := m.Init(ctx, src, sec)
	if err != nil {
		return "", err
	}

	return ctr.WithMountedDirectory("/app", src).
		WithWorkdir("/app").
		WithExec([]string{"yarn", "lint"}).
		Stdout(ctx)
}

func (m *Dagger) TestBackend(
	ctx context.Context,
	// +defaultPath="/"
	src *dagger.Directory,
	sec *dagger.Secret,
) (string, error) {
	ctr, err := m.Init(ctx, src, sec)
	if err != nil {
		return "", err
	}

	return ctr.
		WithExec([]string{"lerna", "run", "test", "--scope", "@vieites-tfg/zoo-backend"}).
		Stdout(ctx)
}

func (m *Dagger) Cypress(ctx context.Context, src *dagger.Directory) *dagger.Container {
	return dag.
		Container().
		From("cypress/browsers").
		WithMountedDirectory("/e2e", src).
		WithWorkdir("/e2e").
		WithExec([]string{"npx", "cypress", "install"}).
		WithExec([]string{"yarn", "add", "lerna@8.2.1", "-W"})
}

func (m *Dagger) TestFrontend(
	ctx context.Context,
	// +defaultPath="/"
	src *dagger.Directory,
	sec *dagger.Secret,
	front *dagger.Service,
) (string, error) {
	_, err := m.Init(ctx, src, sec)
	if err != nil {
		return "", err
	}

	return m.Cypress(ctx, src).
		WithServiceBinding("zoo-frontend", front).
		Terminal().
		WithExec([]string{"yarn", "run", "e2e"}).
		Stdout(ctx)
}

func makeSecrets(ctx context.Context, vars map[string]string) error {
	secrets = make(Secrets)
	client := dagger.Connect()
	for key, value := range vars {
		secrets[key] = client.SetSecret(key, value)
	}

	_, ok := secrets["MONGO_PORT"]
	if !ok {
		secrets["MONGO_PORT"] = client.SetSecret("MONGO_PORT", "27017")
	}

	mongoUri, err := getMongoUri(ctx)
	if err != nil {
		return err
	}
	secrets["MONGODB_URI"] = client.SetSecret("MONGODB_URI", mongoUri)

	return nil
}

func getMongoPort(ctx context.Context) (int, error) {
	mongo_portStr, err := secrets["MONGO_PORT"].Plaintext(ctx)
	if err != nil {
		return 0, err
	}

	mongo_port, err := strconv.Atoi(mongo_portStr)
	if err != nil {
		return 0, err
	}

	return mongo_port, nil
}

func getMongoUri(ctx context.Context) (string, error) {
	var (
		err       error
		root      string
		rootPass  string
		mongoPort string
		db        string
	)

	root, err = secrets.Plaintext(ctx, "MONGO_ROOT")
	if err != nil {
		return "", err
	}

	rootPass, err = secrets.Plaintext(ctx, "MONGO_ROOT_PASS")
	if err != nil {
		return "", err
	}

	mongoPort, err = secrets.Plaintext(ctx, "MONGO_PORT")
	if err != nil {
		return "", err
	}

	db, err = secrets.Plaintext(ctx, "MONGO_DATABASE")
	if err != nil {
		return "", err
	}

	mongoUri := fmt.Sprintf("mongodb://%s:%s@mongodb:%s/%s?authSource=admin",
		root,
		rootPass,
		mongoPort,
		db,
	)

	return mongoUri, nil
}

// parseEnvFile process the content of the .env and returns a variables map
func parseEnvFile(content string) map[string]string {
	envVars := make(map[string]string)
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			value = strings.Trim(value, "\"'")
			envVars[key] = value
		}
	}

	dct := os.Getenv("DAGGER_CLOUD_TOKEN")
	if dct != "" {
		envVars["DAGGER_CLOUD_TOKEN"] = dct
	}

	return envVars
}

type Secrets map[string]*dagger.Secret

func (s Secrets) Get(key string) (*dagger.Secret, error) {
	var (
		value *dagger.Secret
		ok    bool
	)
	if value, ok = s[key]; !ok {
		return nil, fmt.Errorf("%s key not set.", key)
	}

	return value, nil
}

func (s Secrets) Plaintext(ctx context.Context, key string) (string, error) {
	var (
		sec   *dagger.Secret
		value string
		err   error
	)
	if sec, err = s.Get(key); err != nil {
		return "", err
	}

	if value, err = sec.Plaintext(ctx); err != nil {
		return "", err
	}

	return value, nil
}
