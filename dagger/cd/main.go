package main

import (
	"context"
	"dagger/cd/internal/dagger"
	"fmt"
	"strings"
)

type Cd struct {
	// Docker socket to connect to the external Docker Engine. Please carefully
	// use this option it can expose your host to the container.
	// E.g. /var/run/docker.sock
	// +required
	Socket *dagger.Socket

	// It should be the tcp://127.0.0.1 followed by any port.
	// E.g. tcp://127.0.0.1:3000
	// +required
	KindSvc *dagger.Service

	// The name for the cluster.
	// +optional
	ClusterName string

	// The name for the cluster.
	// +optional
	ConfigFile *dagger.File
}

type Envs string

const (
	DEV Envs = "dev"
	PRE Envs = "pre"
	PRO Envs = "pro"
)

func New(
	socket *dagger.Socket,
	kindSvc *dagger.Service,
	// +optional
	clusterName string,
	// +optional
	configFile *dagger.File,
) *Cd {
	if clusterName == "" {
		clusterName = "zoo-cluster"
	}

	return &Cd{
		Socket:      socket,
		KindSvc:     kindSvc,
		ClusterName: clusterName,
		ConfigFile:  configFile,
	}
}

func (m *Cd) Base() *dagger.Container {
	return dag.Container().From("alpine:3.22").
		WithExec([]string{"apk", "update"}).
		WithExec([]string{"apk", "add", "--no-cache", "wget"}).
		WithExec([]string{"sh", "-c", `
      wget https://get.helm.sh/helm-v3.15.2-linux-amd64.tar.gz && \
      tar -zxvf helm-v3.15.2-linux-amd64.tar.gz && \
      mv linux-amd64/helm /usr/local/bin/
    `}).
		WithExec([]string{"sh", "-c", `
      wget https://github.com/helmfile/helmfile/releases/download/v0.165.0/helmfile_0.165.0_linux_amd64.tar.gz && \
      tar -zxvf helmfile_0.165.0_linux_amd64.tar.gz && \
      mv helmfile /usr/local/bin/
    `}).
		WithExec([]string{"sh", "-c", `
      wget -qO /usr/local/bin/yq https://github.com/mikefarah/yq/releases/latest/download/yq_linux_amd64 && \
      chmod a+x /usr/local/bin/yq
    `}).
		WithExec([]string{"sh", "-c", `
      wget https://github.com/getsops/sops/releases/download/v3.8.1/sops-v3.8.1.linux.amd64 -O /usr/local/bin/sops && \
      chmod a+x /usr/local/bin/sops
    `})
}
func (m *Cd) Cluster(ctx context.Context) *dagger.Container {
	base := m.Base()

	kindClient := dag.
		Kind(m.Socket, m.KindSvc, dagger.KindOpts{
			ClusterName: m.ClusterName,
			ConfigFile:  m.ConfigFile,
		}).
		Container()

	helmBinary := base.File("/usr/local/bin/helm")
	helmfileBinary := base.File("/usr/local/bin/helmfile")
	yqBinary := base.File("/usr/local/bin/yq")
	sopsBinary := base.File("/usr/local/bin/sops")

	clusterClientWithTools := kindClient.
		WithFile("/usr/local/bin/helm", helmBinary).
		WithFile("/usr/local/bin/helmfile", helmfileBinary).
		WithFile("/usr/local/bin/yq", yqBinary).
		WithFile("/usr/local/bin/sops", sopsBinary).
		WithExec([]string{"apk", "add", "--no-cache", "git"})

	return clusterClientWithTools.
		WithExec([]string{"mkdir", "-p", "/app"})
}

func (m *Cd) Launch(
	ctx context.Context,

	// `.env` file with the credentials to use the private images and the variables
	// related to the mongo database.
	// +required
	secEnv *dagger.File,

	// environment in which the application will be deployed
	// +required
	env Envs,

	// AGE private key file (e.g., age.agekey)
	// +required
	ageKey *dagger.File,

	// .sops.yaml configuration file for SOPS
	// +required
	sopsConfig *dagger.File,

	// `helmfile.yaml` necessary to launch the chart.
	// +optional
	helmfile *dagger.File,
) (*dagger.Directory, error) {
	ctr := m.Cluster(ctx).
		WithExec([]string{"git", "clone", "https://github.com/vieites-tfg/state.git", "/app/state"})

	if helmfile != nil {
		ctr = ctr.WithFile("/app/state/helmfile.yaml.gotmpl", helmfile)
	}

	ctr, err := setEnvVariables(ctx, ctr, secEnv)
	if err != nil {
		return nil, err
	}

	secretGeneratorFile := `apiVersion: viaduct.ai/v1
kind: ksops
metadata:
  name: secret-generator
  annotations:
    config.kubernetes.io/function: |
      exec:
        path: ksops
files:
  - secrets.yaml # The name of the encrypted files
`

	kustomizationFile := `apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - non-secrets.yaml
generators:
  - secret_generator.yaml
`

	processingScript := fmt.Sprintf(`
		set -euxo pipefail

		echo "--- Templating all resources for environment '%s' ---"
		helmfile -e %s template > /app/all-objects.yaml

		echo "--- Separating secrets from other resources ---"
		yq 'select(.kind == "Secret")' /app/all-objects.yaml > /app/secrets.yaml
		yq 'select(.kind != "Secret")' /app/all-objects.yaml > /app/non-secrets.yaml

		if [ -s /app/secrets.yaml ]; then
			echo "--- Encrypting secrets with SOPS ---"
			sops --encrypt --in-place /app/secrets.yaml
		else
			echo "--- No secrets found to encrypt ---"
		fi

		rm /app/all-objects.yaml
	`, string(env), string(env))

	ctr = ctr.
		WithWorkdir("/app/state").
		// Config and key
		WithFile("/app/state/.sops.yaml", sopsConfig).
		WithExec([]string{"mkdir", "-p", "/root/.config/sops/age/"}).
		WithFile("/root/.config/sops/age/keys.txt", ageKey).
		WithEnvVariable("XDG_CONFIG_HOME", "/root/.config").
		WithNewFile("/app/kustomization.yaml", kustomizationFile).
		WithNewFile("/app/secret_generator.yaml", secretGeneratorFile).
		WithExec([]string{"sh", "-c", processingScript})

	commitScript := fmt.Sprintf(`
		set -euxo pipefail

		NAMESPACE=%s
		git config --global user.email "dvieitest@gmail.com"
		git config --global user.name "Dagger CD Bot"
		
		echo "--- Cloning state repository ---"
		REPO_URL="https://$STATE_REPO@github.com/vieites-tfg/state.git"
		git clone --depth 1 --branch deploy "$REPO_URL" /deploy
		
		cd /deploy
		mkdir -p "$NAMESPACE"
		rm -rf "$NAMESPACE"/* || true
		
		echo "--- Committing template to state repo ---"
		mv /app/non-secrets.yaml "$NAMESPACE"/non-secrets.yaml
		if [ -s /app/secrets.yaml ]; then
			mv /app/secrets.yaml "$NAMESPACE"/secrets.yaml
			mv /app/kustomization.yaml "$NAMESPACE"/kustomization.yaml
			mv /app/secret_generator.yaml "$NAMESPACE"/secret_generator.yaml
		fi
		
		git add .
		if git diff --staged --quiet; then
			echo "No changes to commit."
		else
			git commit -m "Update manifests for $NAMESPACE"
			git push origin deploy
		fi
	`, string(env))

	ctr = ctr.WithExec([]string{"sh", "-c", commitScript})

	finalState := ctr.Directory(fmt.Sprintf("/deploy/%s", string(env)))

	return finalState, nil
}

func setEnvVariables(
	ctx context.Context,
	ctr *dagger.Container,
	env *dagger.File,
) (*dagger.Container, error) {
	envContents, err := env.Contents(ctx)

	if err != nil {
		return nil, err
	}

	lines := strings.Split(envContents, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := parts[0]
			value := parts[1]
			secretValue := dag.SetSecret(key, value)
			ctr = ctr.WithSecretVariable(key, secretValue)
		}
	}
	return ctr, nil
}
