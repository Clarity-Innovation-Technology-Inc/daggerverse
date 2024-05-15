package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"dagger/cue-file-editor/internal/cue"
	"dagger/cue-file-editor/internal/github"
)

type CueFileEditor struct {
	GithubOwner         string
	GithubRepo          string
	GithubBranch        string
	GithubCommitterName  string
	GithubCommitterEmail string
	GithubCommitMessage string
	GithubRepoPath      string
	CueFileName         string
	CuePath             string
	NewStringValue      string
}

func New(
	// the owner of the github org that stores the cue file you want to edit
	//
	githubOwner string,

	// the name of the repo that stores the cue file you want to edit
	//
	githubRepo string,

	// the branch of the repo that stores the cue file you want to edit
	//
	//+optional
	//+default="main"
	githubBranch string,

	// the name of the github user that will commit the changes to your cue file
	//
	githubCommiterName string,

	// the email of the github user that will commit the changes to your cue file
	//
	githubCommiterEmail string,

	// the commit message that will be used when committing the cue file changes
	//
	//+optional
	//+default="update cue value"
	githubCommitMsg string,

	// the subdirectory within the github repo that contains the cue file
	//
	//+optional
	//+default=""
	githubRepoPath string,

	// the name of the cue config file you want to update
	//
	//+optional
	//+default="values.cue"
	cueFile string,

	// a cue path (dot seperated) to a specific key within a cue configuration file
	//
	cuePath string,

	// a concrete string value you want to use to update your cue file
	//
	newStringVal string,

) *CueFileEditor {
	return &CueFileEditor{
		GithubOwner:         githubOwner,
		GithubRepo:          githubRepo,
		GithubBranch:        githubBranch,
		GithubCommitterName:  githubCommiterName,
		GithubCommitterEmail: githubCommiterEmail,
		GithubCommitMessage: githubCommitMsg,
		GithubRepoPath:      githubRepoPath,
		CueFileName:         cueFile,
		CuePath:             cuePath,
		NewStringValue:      newStringVal,
	}
}

func (c *CueFileEditor) WithGithubBranch(branch string) *CueFileEditor {
	c.GithubBranch = branch
	return c
}
func (c *CueFileEditor) WithGithubCommitMessage(msg string) *CueFileEditor {
	c.GithubCommitMessage = msg
	return c
}
func (c *CueFileEditor) WithGithubRepoPath(path string) *CueFileEditor {
	c.GithubRepoPath = path
	return c
}
func (c *CueFileEditor) WithCueFileName(filename string) *CueFileEditor {
	c.CueFileName = filename
	return c
}

func (c *CueFileEditor) Update(token *Secret) {
	url := fmt.Sprintf("https://github.com/%s/%s.git", c.GithubOwner, c.GithubRepo)
	dir := dag.Git(url).
		WithAuthToken(token).
		Branch(c.GithubBranch).
		Tree()
	cueFile := dir.File(fmt.Sprintf("%s/%s", c.GithubRepoPath, c.CueFileName))
	_, err := cueFile.Export(context.Background(), c.CueFileName)
	if err != nil {
		panic(err)
	}

	// decoding the cue values as go data stored in memory
	values, err := cue.Decode(c.CueFileName)
	if err != nil {
		panic(err)
	}

	// finally editing the image tag to the desired tag
	newCueVal, err := cue.UpdateValues(values, c.CuePath, c.NewStringValue)
	if err != nil {
		panic(err)
	}

	// writing the updated cue file back to disk
	err = os.WriteFile(c.CueFileName, []byte(fmt.Sprint(newCueVal)), 0644)
	if err != nil {
		panic(err)
	}

	tokenString, err := token.Plaintext(context.Background())
	if err != nil {
		panic(err)
	}

	// now commit this damn thing to github and were good
	cli := github.New(
		tokenString,
		c.GithubOwner,
		c.GithubRepo,
		c.GithubBranch,
		c.GithubRepoPath,
		[]string{c.CueFileName},
		c.GithubCommitterName,
		c.GithubCommitterEmail,
		c.GithubCommitMessage,
	)

	ref, statusCode, err := cli.Commit()
	if err != nil {
		msg := fmt.Sprintf("commit failed with status code: %d. %v", statusCode, err)
		panic(msg)
	}

	if statusCode == http.StatusOK {
		fmt.Printf("commit succeeded with status code: %d. New Commit SHA: %s\n", statusCode, *ref.Object.SHA)
	}
}
