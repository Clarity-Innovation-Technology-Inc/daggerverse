package main

type Github struct {
	URL    string
	Branch string
}

func New(
	// The https url for the github repo to fetch.
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

func (g *Github) WithURL(url string) *Github {
	g.URL = url
	return g
}

func (g *Github) WithBranch(branch string) *Github {
	g.Branch = branch
	return g
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
