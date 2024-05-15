package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"dagger/cue-file-editor/internal/cue"
	// "dagger/cue-file-editor/internal/dagger"
	"dagger/cue-file-editor/internal/github"
)

type CueFileEditor struct {
	GithubOwner         string
	GithubRepo          string
	GithubBranch        string
	GithubCommiterName  string
	GithubCommiterEmail string
	GithubCommitMessage string
	GithubRepoPath      string
	CueFileName         string
	CuePath             string
	NewStringValue      string
}

func New(
	// The owner of the github org that stores the cue file you want to edit
	//
	//+optional
	githubOwner string,

	// the name of the repo that stores the cue file you want to edit
	//
	//+optional
	githubRepo string,

	// the branch of the repo that stores the cue file you want to edit
	//
	//+optional
	githubBranch string,

	// the name of the github user that will commit the changes to your cue file
	//
	//+optional
	githubCommiterName string,

	// the email of the github user that will commit the changes to your cue file
	//
	//+optional
	githubCommiterEmail string,

	// The commit message that will be used when committing the cue file changes
	//
	//+optional
	githubCommitMsg string,

	//	the subdirectory within the github repo that contains the cue file
	//
	//+optional
	githubRepoPath string,

	// The name of the cue config file you want to update
	//
	//+default="values.cue"
	//+optional
	cueFile string,

	// A cue path (dot seperated) to a specific key within a cue configuration file
	//
	// //+optional
	cuePath string,

	// A concrete string value you want to use to your cue file
	//
	// //+optional
	newStringVal string,

) *CueFileEditor {
	return &CueFileEditor{
		GithubOwner:         githubOwner,
		GithubRepo:          githubRepo,
		GithubBranch:        githubBranch,
		GithubCommiterName:  githubCommiterName,
		GithubCommiterEmail: githubCommiterEmail,
		GithubCommitMessage: githubCommitMsg,
		GithubRepoPath:      githubRepoPath,
		CueFileName:         cueFile,
		CuePath:             cuePath,
		NewStringValue:      newStringVal,
	}
}

func (c *CueFileEditor) WithGithubOwner(owner string) *CueFileEditor {
	c.GithubOwner = owner
	return c
}
func (c *CueFileEditor) WithGithubRepo(repo string) *CueFileEditor {
	c.GithubRepo = repo
	return c
}
func (c *CueFileEditor) WithGithubBranch(branch string) *CueFileEditor {
	c.GithubBranch = branch
	return c
}
func (c *CueFileEditor) WithGithubCommiterName(name string) *CueFileEditor {
	c.GithubCommiterName = name
	return c
}
func (c *CueFileEditor) WithGithubCommiterEmail(email string) *CueFileEditor {
	c.GithubCommiterEmail = email
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
func (c *CueFileEditor) WithCuePath(path string) *CueFileEditor {
	c.CuePath = path
	return c
}
func (c *CueFileEditor) WithNewStringValue(val string) *CueFileEditor {
	c.NewStringValue = val
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
		fmt.Println(err)
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
		c.GithubCommiterName,
		c.GithubCommiterEmail,
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
