package main

import (
	"context"
	// "dagger/github/internal/dagger"
	// "fmt"
	// "os"
)

type Github struct {
	Repo   *Directory
	Branch string
	Cntr   *Container
}

func New(
	// The ssh url for the github repo to fetch.
	//
	// +optional
	repo *Directory,

	// The branch of the Github repo to fetch.
	//
	//+optional
	//+default="main"
	branch string,
) *Github {
	return &Github{
		Repo:   repo,
		Branch: branch,
	}
}

func (g *Github) WithRepo(repo *Directory) *Github {
	g.Repo = repo
	return g
}

func (g *Github) WithBranch(branch string) *Github {
	g.Branch = branch
	return g
}

// When you have a Directory argument and you pass a git URL through the CLI,
// it adds your ssh socket if you have SSH_AUTH_SOCK set
func (g *Github) Container(ctx context.Context) (*Container, error) {
	cntr := dag.Container().
		From("alpine:latest").
		WithDirectory("/src", g.Repo)

	g.Cntr = cntr

	return cntr, nil
}
