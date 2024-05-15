package github

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v61/github"
)

type Github struct {
	Client         *github.Client
	Owner          string
	Repo           string
	Branch         string
	SubDir         string
	Files          []string
	CommitterName  string
	CommitterEmail string
	CommitMessage  string
}

func New(
	token string,
	owner string,
	repo string,
	branch string,
	subDir string,
	files []string,
	committerName string,
	committerEmail string,
	commitMessage string,
) *Github {
	return &Github{
		Client:         github.NewClient(nil).WithAuthToken(token),
		Owner:          owner,
		Repo:           repo,
		Branch:         branch,
		SubDir:         subDir,
		Files:          files,
		CommitterName:  committerName,
		CommitterEmail: committerEmail,
		CommitMessage:  commitMessage,
	}
}

// getFileContent loads the local content of a file and return the target name
// of the file in the target repository and its contents.
func getFileContent(fileArg string) (targetName string, b []byte, err error) {
	var localFile string
	files := strings.Split(fileArg, ":")
	switch {
	case len(files) < 1:
		return "", nil, errors.New("empty `-files` parameter")
	case len(files) == 1:
		localFile = files[0]
		targetName = files[0]
	default:
		localFile = files[0]
		targetName = files[1]
	}

	b, err = os.ReadFile(localFile)
	return targetName, b, err
}

func (g *Github) Commit() (ref *github.Reference, statusCode int, err error) {
	ctx := context.Background()

	// first we are getting the latest commit on the requested branch
	ref, resp, err := g.Client.Git.GetRef(ctx, g.Owner, g.Repo, "refs/heads/"+g.Branch)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	//  here we are	creating a new git tree, based off of the latest commit on the requested branch
	//  and appending your newly edited file that you want to commit
	entries := []*github.TreeEntry{}
	//  so this is supposed to support committing multiple files
	for _, f := range g.Files {
		file, content, err := getFileContent(f)
		if err != nil {
			return nil, 0, err
		}
		entries = append(
			entries,
			&github.TreeEntry{
				Path:    github.String(fmt.Sprintf("%s/%s", g.SubDir, file)),
				Type:    github.String("blob"),
				Content: github.String(string(content)),
				Mode:    github.String("100644"),
			},
		)
	}
	tree, resp, err := g.Client.Git.CreateTree(ctx, g.Owner, g.Repo, *ref.Object.SHA, entries)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	// fetching the parent commit sha to ensure that our new commit is \
	// attached properly to the branch's latest
	parent, resp, err := g.Client.Repositories.GetCommit(ctx, g.Owner, g.Repo, *ref.Object.SHA, nil)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	parent.Commit.SHA = parent.SHA

	// we are in commit creation territory now
	date := time.Now()
	author := &github.CommitAuthor{
		Date:  &github.Timestamp{Time: date},
		Name:  github.String(g.CommitterName),
		Email: github.String(g.CommitterEmail),
	}
	commit := &github.Commit{
		Author:  author,
		Message: github.String(g.CommitMessage),
		Tree:    tree,
		Parents: []*github.Commit{parent.Commit},
	}
	opts := github.CreateCommitOptions{}
	newCommit, resp, err := g.Client.Git.CreateCommit(ctx, g.Owner, g.Repo, commit, &opts)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	// once the commit has been created, we need to essentially push it
	//  so we update the orgin's ref to our new commit sha
	ref.Object.SHA = newCommit.SHA // ref here is the original ref to the branch before our new commit
	updatedRef, resp, err := g.Client.Git.UpdateRef(ctx, g.Owner, g.Repo, ref, false)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	return updatedRef, resp.StatusCode, nil
}
