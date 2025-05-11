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

func getVersion(ctx context.Context, ctr *dagger.Container, pkg string) (string, error) {
	return ctr.
		WithWorkdir(fmt.Sprintf("/app/packages/%s", pkg)).
		WithExec([]string{"node", "-p", "require('./package.json').version"}).
		Stdout(ctx)
}

func getTags(ctx context.Context, base *dagger.Container, pkg string) ([]string, error) {
	version, err := getVersion(ctx, base, pkg)
	if err != nil {
		return nil, err
	}

	return []string{"latest", strings.TrimSpace(version)}, nil
}

func PublishImage(
	ctx context.Context,
	base *dagger.Container,
	image *dagger.Container,
	pkg string,
	sec *dagger.Secret,
) ([]string, error) {
	tags, err := getTags(ctx, base, pkg)
	if err != nil {
		return tags, err
	}

	var out []string
	for _, t := range tags {
		imageRef, err := image.WithRegistryAuth("ghcr.io", "vieitesss", sec).
			Publish(ctx, fmt.Sprintf("ghcr.io/vieites-tfg/zoo-%s:%s", pkg, t))
		if err != nil {
			return out, err
		}
		out = append(out, imageRef)
	}

	return out, nil
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
	config *dagger.File,
	pkg string,
	pat *dagger.Secret,
) (string, error) {
	return base.
		WithSecretVariable("CR_PAT", pat).
		WithFile("/app/.npmrc", config).
		WithExec([]string{"yarn", "publish", "--access", "restricted", fmt.Sprintf("/app/packages/%s", pkg), "--non-interactive"}).
		Stdout(ctx)
}
