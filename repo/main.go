package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

type Repo struct {
	URL          string
	Branch       string
	SSHAgentPath string // SSH_AUTH_SOCK
}

func NewRepo(
	// +optional
	// +default="SSH_AUTH_SOCK"
	sshAgentPath string,
) *Repo {
	return &Repo{
		SSHAgentPath: sshAgentPath,
	}
}

func (r *Repo) WithUrl(url string) *Repo {
	r.URL = url
	return r
}

func (r *Repo) WithBranch(branch string) *Repo {
	r.Branch = branch
	return r
}

func (r *Repo) WithSSHAgentPath(sshAgentPath string) *Repo {
	r.SSHAgentPath = sshAgentPath
	return r
}

func (r *Repo) Container() (*dagger.Container, error) {
	ctx := context.Background()
	cli, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		return nil, err
	}

	sshAgentPath := os.Getenv(r.SSHAgentPath)
	repo := cli.Git(
		r.URL,
		dagger.GitOpts{
			SSHAuthSocket: cli.Host().UnixSocket(sshAgentPath),
		}).
		Branch(r.Branch).
		Tree()
	if repo == nil {
		return nil, fmt.Errorf("invalid Git repository or branch: %s/%s", r.URL, r.Branch)
	}

	cont := cli.Container().
		From("alpine:latest").
		WithDirectory("/src", repo, dagger.ContainerWithDirectoryOpts{})

	return cont, nil
}
