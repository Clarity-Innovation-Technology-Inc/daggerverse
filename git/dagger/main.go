// A generated module for Git functions
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
	auth "dagger/git/internal/pkg"
	"fmt"
)

type Git struct{}

// Returns a container that echoes whatever string argument is provided
func (m *Git) ContainerEcho(stringArg string) *Container {
	return dag.Container().From("alpine:latest").WithExec([]string{"echo", stringArg})
}

// Returns lines that match a pattern in the files of the provided Directory
func (m *Git) GrepDir(ctx context.Context, directoryArg *Directory, pattern string) (string, error) {
	return dag.Container().
		From("alpine:latest").
		WithMountedDirectory("/mnt", directoryArg).
		WithWorkdir("/mnt").
		WithExec([]string{"grep", "-R", pattern, "."}).
		Stdout(ctx)
}

func Github(url string, branch string) {

	ctx := context.Background()
	client, err := auth.NewClient(ctx)
	if err != nil {
		fmt.Println(err) // SSHTOEKEN + FUNCTION WORKS.
	}
	// git@github.com:Clarity-Innovation-Technology-Inc/core-devops-k8s-apps.git
	container, err := client.LoadGithubRepo(ctx, url, branch)
	if err != nil {
		fmt.Println(err)
	}

	dir, err := container.Directory("/src").Entries(ctx)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(dir)
	// return dir, nil
}
