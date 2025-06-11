package main

import (
	"context"
	"dagger/cd/internal/dagger"
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
	return dag.Container().From("alpine:3.18").
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
		`})
}

func (m *Cd) Cluster(ctx context.Context) *dagger.Container {
	base := m.Base()

	kindClient := dag.
		Kind(m.Socket, m.KindSvc, dagger.KindOpts{
			ClusterName: m.ClusterName,
			ConfigFile: m.ConfigFile,
		}).
		Container()

	helmBinary := base.File("/usr/local/bin/helm")
	helmfileBinary := base.File("/usr/local/bin/helmfile")

	clusterClientWithTools := kindClient.
		WithFile("/usr/local/bin/helm", helmBinary).
		WithFile("/usr/local/bin/helmfile", helmfileBinary).
		WithExec([]string{"apk", "add", "--no-cache", "git"})

	applyIngress := []string{"sh", "-c", `
		kubectl apply -f \
		https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml &&
		sleep 0.5 &&
		kubectl wait --namespace ingress-nginx \
		--for=condition=ready pod \
		--selector=app.kubernetes.io/component=controller \
		--timeout=90s
	`}

	return clusterClientWithTools.
		WithExec(applyIngress).
		WithExec([]string{"kubectl", "create", "namespace", "dev"}).
		WithExec([]string{"kubectl", "create", "namespace", "pre"}).
		WithExec([]string{"kubectl", "create", "namespace", "pro"}).
		WithExec([]string{"mkdir", "-p", "/app"})
}

func (m *Cd) Launch(
	ctx context.Context,

	// `.env` file with the credentials to use the private images and the variables
	// related to the mongo database.
	// +required
	secEnv *dagger.File,

	// `helmfile.yaml` necessary to launch the chart.
	// +optional
	helmfile *dagger.File,
) (*dagger.Container, error) {
	ctr := m.Cluster(ctx).
		WithExec([]string{"git", "clone", "https://github.com/vieites-tfg/values.git", "/app/values"})

	if helmfile != nil {
		ctr = ctr.WithFile("/app/values/helmfile.yaml.gotmpl", helmfile)
	}

	ctr, err := setEnvVariables(ctx, ctr, secEnv)
	if err != nil {
		return nil, err
	}

	ctr = ctr.
		WithWorkdir("/app/values").
		WithExec([]string{"helmfile", "-e", "dev", "sync"})
	
	return ctr, nil
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
