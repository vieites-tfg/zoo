package main

import (
	"context"
	"dagger/cd/internal/dagger"
)

type Cd struct{
	// Docker socket to connect to the external Docker Engine. Please carefully
	// use this option it can expose your host to the container.
	// E.g. /var/run/docker.sock
	// +required
	Socket *dagger.Socket

	// `.env` file with the credentials to use the private images and the variables
	// related to the mongo database.
	// +required
	SecEnv *dagger.File

	// `helmfile.yaml` necessary to launch the chart.
	// +optional
	Helmfile *dagger.File

	// The name for the cluster.
	// +required
	clusterName string
}

func New(
	socket *dagger.Socket,
	secEnv *dagger.File,
	//+optional
	helmfile *dagger.File) *Cd {
	return &Cd{
		Socket: socket,
		Helmfile: helmfile,
		SecEnv: secEnv,
		clusterName: "zoo-cluster",
	}
}

func (m *Cd) Cluster(
	ctx context.Context,

	// An executable file to create the kind cluster
	// +required
	script *dagger.File,
) *dagger.Container {
	ctr := dag.
		Kind(m.Socket).
		Container().
		WithExec([]string{"apk", "update"}).
		WithExec([]string{"apk", "add", "--no-cache", "git"}).
		WithExec([]string{"git", "clone", "https://github.com/vieites-tfg/values.git", "/app/values"}).
		WithExec([]string{"wget", "https://github.com/helmfile/helmfile/releases/download/v1.1.0/helmfile_1.1.0_linux_386.tar.gz"}).
		WithExec([]string{"tar", "xzvf", "helmfile_1.1.0_linux_386.tar.gz"}).
		WithExec([]string{"mv", "./helmfile", "/usr/local/bin"}).
		WithExec([]string{"mkdir", "-p", "/app"}).
		WithFile("/app/.env", m.SecEnv).
		WithFile("/app/create_cluster.sh", script, dagger.ContainerWithFileOpts{Permissions: 0555})
	
		if m.Helmfile == (&dagger.File{}) {
			ctr = ctr.WithFile("/app/values/helmfile.yaml.gotmpl", m.Helmfile)
		}

	return ctr.WithExec([]string{"/app/create_cluster.sh"})
}

// Should call Cd.Cluster and install the values in order to launch everything
func (m *Cd) Launch() (string, error) {
	return "", nil
}
