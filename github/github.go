package main

import (
	"github.com/Clarity-Innovation-Technology-Inc/daggerverse/github/internal/dagger"
	"fmt"
)

type Github struct {
	URL    string
	Branch string
}

func New(
	// The url for the github repo to fetch.
	//
	// +optional
	url string,

	// The branch of the Github repo to fetch.
	//
	//+optional
	//+default="main"
	branch string,
) *Github {
	return &Github{
		URL:    url,
		Branch: branch,
	}
}

func (g *Github) WithUrl(url string) *Github {
	g.URL = url
	return g
}

func (g *Github) WithBranch(branch string) *Github {
	g.Branch = branch
	return g
}

func (g *Github) Container(sshSocket *Socket) (*Container, error) {
	repo := dag.Git(
		g.URL,
		dagger.GitOpts{
			SSHAuthSocket: sshSocket,
		}).
		Branch(g.Branch).
		Tree()
	if repo == nil {
		return nil, fmt.Errorf("invalid Git repository or branch: %s/%s", g.URL, g.Branch)
	}

	cntr := dag.Container().
		From("alpine:latest").
		WithDirectory("/src", repo, dagger.ContainerWithDirectoryOpts{})

	return cntr, nil
}
