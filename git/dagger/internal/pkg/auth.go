package auth

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

type Client struct {
	*dagger.Client
}

func NewClient(ctx context.Context) (*Client, error) {
	// initialize Dagger client
	cli, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		return nil, err
	}

	// Wrap the dagger.Client in your local Client type
	return &Client{cli}, nil
}

func (c *Client) Close() error {
	return c.Client.Close()
}

func (c *Client) LoadGithubRepo(ctx context.Context, repoURL, branch string) (*dagger.Container, error) {
	// Retrieve path of authentication agent socket from host
	sshAgentPath := os.Getenv("SSH_AUTH_SOCK") // how do we get the ssh key without exposing it.
	// Write a module ls -la. with the org private directory.
	//
	repo := c.Git(repoURL,
		dagger.GitOpts{
			SSHAuthSocket: c.Host().UnixSocket(sshAgentPath),
		}).
		Branch(branch).
		Tree()
	if repo == nil {
		return nil, fmt.Errorf("invalid Git repository or branch: %s/%s", repoURL, branch)
	}

	// Clone the Git repository into the container
	container := c.Container().
		From("alpine:latest").
		WithDirectory("/src", repo, dagger.ContainerWithDirectoryOpts{})

	return container, nil
}

func (c *Client) GetFileBytes(ctx context.Context, container *dagger.Container, workingDir, fileName string) ([]byte, error) {
	// Access the specified values file within the container
	valuesFile := container.File(
		fmt.Sprintf("/src/%s/%s", workingDir, fileName),
	)

	fileContent, err := valuesFile.Contents(ctx)
	if err != nil {
		return nil, err
	}

	return []byte(fileContent), nil
}
