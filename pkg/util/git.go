package util

import (
	"io/ioutil"
	"os"
	"strings"

  "gopkg.in/src-d/go-git.v4"
  "gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)


func CloneRepo(srcUrl, srcVer string) (string, error) {

	co := &git.CloneOptions{
			URL: srcUrl,
			// Progress: os.Stdout,
			// SingleBranch: true,
	}

	if strings.Contains(srcUrl, "github.com") && os.Getenv("GITHUB_TOKEN") != "" {
		co.Auth = &http.BasicAuth{
			Username: "github-token", // yes, this can be anything except an empty string
			Password: os.Getenv("GITHUB_TOKEN"),
		}
	}

	if srcVer != "" {
		co.ReferenceName = plumbing.ReferenceName(srcVer)
	}

	// temp dir to clone to
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		return "", err
	}

	// Clones the repository into the worktree (fs) and storer all the .git
	// content into the storer
	_, err = git.PlainClone(dir, false, co)
	if err != nil {
		return dir, err
	}

	return dir, nil
}
