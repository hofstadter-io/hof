package yagu

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func FindGitRepoAbsPath(dir string) (string, error) {
	var err error
	if dir != "" {
		dir, err = os.Getwd()
		if err != nil {
			return "", err
		}
	}

	found := false

	for !found && dir != "/" {
		try := filepath.Join(dir, ".git")
		info, err := os.Stat(try)
		if err == nil && info.IsDir() {
			found = true
			break
		}

		next := filepath.Clean(filepath.Join(dir, ".."))
		dir = next
	}

	if !found {
		return "", nil
	}

	return dir, nil
}

// FindRemoteRepoRootAndClone walks a path back until
// it finds a valid git repo, cloning it into a tmp dir
func FindRemoteRepoRootAndClone(repo, ver string) (root, dir string, err error) {
	parts := strings.Split(repo, "/")
	for i := len(parts); i > 1; i-- {
		r := strings.Join(parts[:i], "/")
		fmt.Println("trying...", r)
		dir, err = CloneRepoIntoTmp("https://" + r, ver)
		if err == nil {
			return r, dir, nil
		}
		err = os.RemoveAll(dir)
		if err != nil {
			return r, dir, err
		}
	}

	return "", "", fmt.Errorf("failed to find remote git repo from %s@%s", repo, ver)
}

func CloneRepoIntoTmp(srcUrl, srcVer string) (string, error) {

	co, err := SetupGitOptions(srcUrl, srcVer)
	if err != nil {
		return "", err
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

func CloneRepoIntoDir(srcUrl, srcVer, dest string) error {

	co, err := SetupGitOptions(srcUrl, srcVer)
	if err != nil {
		return err
	}

	// Clones the repository into the worktree (fs) and storer all the .git
	// content into the storer
	_, err = git.PlainClone(dest, false, co)
	if err != nil {
		return err
	}

	return nil
}

func SetupGitOptions(srcUrl, srcVer string) (*git.CloneOptions, error) {
	co := &git.CloneOptions{
		URL: srcUrl,
		// Progress: os.Stdout,
		// SingleBranch: true,
	}

	if srcVer != "" {
		co.ReferenceName = plumbing.ReferenceName(srcVer)
	}

	err := SetupGitAuth(srcUrl, srcVer, co)
	if err != nil {
		return co, err
	}

	return co, nil
}

func SetupGitAuth(srcUrl, srcVer string, co *git.CloneOptions) error {

	// GitHub variations
	if strings.Contains(srcUrl, "github.com") {
		if os.Getenv("GITHUB_TOKEN") != "" {
			co.Auth = &http.BasicAuth{
				Username: "github-token", // yes, this can be anything except an empty string
				Password: os.Getenv("GITHUB_TOKEN"),
			}
			return nil
		}
	}

	// GitLab variations
	if strings.Contains(srcUrl, "gitlab.com") {
		if os.Getenv("GITLAB_TOKEN") != "" {
			co.Auth = &http.TokenAuth{
				Token: os.Getenv("GITLAB_TOKEN"),
			}
			return nil
		}
		if os.Getenv("GITLAB_USERNAME") != "" && os.Getenv("GITLAB_PASSWORD") != "" {
			co.Auth = &http.BasicAuth{
				Username: os.Getenv("GITLAB_USERNAME"),
				Password: os.Getenv("GITLAB_PASSWORD"),
			}
			return nil
		}
	}

	// BitBucket variations
	if strings.Contains(srcUrl, "bitbucket.org") {
		if os.Getenv("BITBUCKET_TOKEN") != "" {
			co.Auth = &http.TokenAuth{
				Token: os.Getenv("BITBUCKET_TOKEN"),
			}
			return nil
		}
		if os.Getenv("BITBUCKET_USERNAME") != "" && os.Getenv("BITBUCKET_PASSWORD") != "" {
			co.Auth = &http.BasicAuth{
				Username: os.Getenv("BITBUCKET_USERNAME"),
				Password: os.Getenv("BITBUCKET_PASSWORD"),
			}
			return nil
		}
	}

	return nil
}
