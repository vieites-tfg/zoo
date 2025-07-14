package main

import (
	"context"
	"dagger/dagger/internal/dagger"
	"fmt"
	"os"
	"strings"
)

type secrets map[string]*dagger.Secret

func MakeSecrets(ctx context.Context, envContent string) (secrets, error) {
	secrets := make(secrets)

	vars := parseEnvFile(envContent)
	client := dagger.Connect()
	for key, value := range vars {
		secrets[key] = client.SetSecret(key, value)
	}

	_, ok := secrets["MONGO_PORT"]
	if !ok {
		secrets["MONGO_PORT"] = client.SetSecret("MONGO_PORT", "27017")
	}

	return secrets, nil
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

func PublishImage(
	ctx context.Context,
	image *dagger.Container,
	pkg string,
	sec *dagger.Secret,
	tag string,
) (string, error) {
	return image.
		WithRegistryAuth("ghcr.io", "vieitesss", sec).
		Publish(ctx, fmt.Sprintf("ghcr.io/vieites-tfg/zoo-%s:%s", pkg, tag))
}

func Lint(ctx context.Context, base *dagger.Container, pkg string) (string, error) {
	return base.
		WithWorkdir(fmt.Sprintf("/app/packages/%s", pkg)).
		WithExec([]string{"yarn", "lint"}).
		Stdout(ctx)
}

func PublishPkg(
	ctx context.Context,
	base *dagger.Container,
	pkg string,
	pat *dagger.Secret,
) (string, error) {
	return base.
		WithSecretVariable("CR_PAT", pat).
		WithNewFile("/app/.npmrc", "//npm.pkg.github.com/:_authToken=${CR_PAT}\n").
		WithExec([]string{"yarn", "publish", "--access", "restricted", fmt.Sprintf("/app/packages/%s", pkg), "--non-interactive"}).
		Stdout(ctx)
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
