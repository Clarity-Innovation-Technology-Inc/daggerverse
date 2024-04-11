package main

import (
	"dagger/github/internal/dagger"
	"fmt"
)

type Repo struct {
	URL    string
	Branch string
}

func New() *Repo {
	return &Repo{}
}

func (r *Repo) WithUrl(url string) *Repo {
	r.URL = url
	return r
}

func (r *Repo) WithBranch(branch string) *Repo {
	r.Branch = branch
	return r
}

func (r *Repo) Container(sshSocket *Socket) (*Container, error) {
	repo := dag.Git(
		r.URL,
		dagger.GitOpts{
			SSHAuthSocket: sshSocket,
		}).
		Branch(r.Branch).
		Tree()
	if repo == nil {
		return nil, fmt.Errorf("invalid Git repository or branch: %s/%s", r.URL, r.Branch)
	}

	cntr := dag.Container().
		From("alpine:latest").
		WithDirectory("/src", repo, dagger.ContainerWithDirectoryOpts{})

	return cntr, nil
}
