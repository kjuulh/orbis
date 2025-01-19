// A generated module for Ci functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
	"dagger/ci/internal/dagger"
)

type Ci struct{}

func (m *Ci) Main(ctx context.Context, service string, source *dagger.Directory, version string) (string, error) {
	golangBase := dag.
		Container().
		From("golang:1.23.5-alpine")

	serviceBuild := golangBase.
		WithWorkdir("/mnt/src").
		WithDirectory(".", source, dagger.ContainerWithDirectoryOpts{
			Include: []string{
				"**/go.mod",
				"**/go.sum",
				"**/go.work",
				"**/go.work.sum",
			},
		}).
		WithExec([]string{"go", "mod", "download"}).
		WithDirectory(".", source, dagger.ContainerWithDirectoryOpts{
			Exclude: []string{"**/*.git/", "**/*.jj/", "**/*.cuddle/"},
		}).
		WithEnvVariable("CGO_ENABLED", "0").
		WithExec([]string{"go", "build", "-o", "./dist/" + service, "./cmd/" + service})

	binImage := dag.
		Container().
		WithFile("/usr/local/bin/"+service, serviceBuild.File("dist/"+service))

	published, err := binImage.Publish(ctx, "docker.io/kasperhermansen/"+service+":"+version)
	if err != nil {
		return "", err
	}

	return published, nil
}
