package main

import (
	"context"
	"dagger/dagger/internal/dagger"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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

func MakeSecrets(ctx context.Context, vars map[string]string) error {
	secrets = make(Secrets)
	client := dagger.Connect()
	for key, value := range vars {
		secrets[key] = client.SetSecret(key, value)
	}

	_, ok := secrets["MONGO_PORT"]
	if !ok {
		secrets["MONGO_PORT"] = client.SetSecret("MONGO_PORT", "27017")
	}

	mongoUri, err := GetMongoUri(ctx)
	if err != nil {
		return err
	}
	secrets["MONGODB_URI"] = client.SetSecret("MONGODB_URI", mongoUri)

	return nil
}

func GetMongoPort(ctx context.Context) (int, error) {
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

func GetMongoUri(ctx context.Context) (string, error) {
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
func ParseEnvFile(content string) map[string]string {
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
