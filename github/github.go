package main

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
)

const (
	ValidURLScheme = "https"
	ValidURLHost   = "github.com"
)

type Github struct {
	URL    string
	Branch string
}

func New(
	// The url for the Github repo to fetch. Only supports git cloning via https!
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

func (g *Github) WithURL(addr string) (*Github, error) {
	// parsing the url for validation purposes
	parsed, err := url.Parse(addr)
	if err != nil {
		return &Github{}, err
	}

	// sanitization to ensure that the url uses https
	if parsed.Scheme != ValidURLScheme {
		return &Github{}, fmt.Errorf("url scheme must be %s, got %s", ValidURLScheme, parsed.Scheme)
	}

	// sanitization to ensure that the url host is in fact github.com
	if parsed.Host != ValidURLHost {
		return &Github{}, fmt.Errorf("host must be %s, got %s", ValidURLHost, parsed.Host)
	}

	g.URL = addr
	return g, nil
}

func (g *Github) WithBranch(branch string) (*Github, error) {
	// limit branch string char types and length
	branchRegex := regexp.MustCompile(`^[A-Za-z0-9-]{1,45}$`)
	if !branchRegex.MatchString(branch) {
		return &Github{}, errors.New("invalid branch name, strict regex match required")
	}

	g.Branch = branch
	return g, nil
}

func (g *Github) Repo(path string, token *Secret) *Container {
	repo := dag.Git(g.URL).
		WithAuthToken(token).
		Branch(g.Branch).
		Tree()

	return dag.Container().
		From("alpine:latest").
		WithDirectory(path, repo)
}
