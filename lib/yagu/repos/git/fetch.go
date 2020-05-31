package git

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-billy/v5/osfs"
	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
)

func NewRemote(srcUrl string) (*GitRepo, error) {

	rc := &config.RemoteConfig{
		Name: "origin",
		URLs: []string{
			"https://" + srcUrl,
		},
	}

	lo := &gogit.ListOptions{}

	if strings.Contains(srcUrl, "github.com") && os.Getenv("GITHUB_TOKEN") != "" {
		lo.Auth = &http.BasicAuth{
			Username: "github-token", // yes, this can be anything except an empty string
			Password: os.Getenv("GITHUB_TOKEN"),
		}
		// co.URL = "git@" + strings.Replace(srcUrl, "/", ":", 1)
	}

	// fmt.Println("URL:", rc.URLs[0])

	// Clones the repository into the worktree (fs) and storer all the .git
	// content into the storer
	st := memory.NewStorage()
	remote := gogit.NewRemote(st, rc)

	return &GitRepo{
		Store:       st,
		Remote:      remote,
		ListOptions: lo,
	}, nil
}

func CloneLocalRepo(location string) (*GitRepo, error) {
	fs := osfs.New(location)

	// Only returning the Billy FS in this case
	return &GitRepo{
		FS: fs,
	}, nil
}

func CloneRepoRef(srcUrl string, ref *plumbing.Reference) (*GitRepo, error) {

	co := &gogit.CloneOptions{
		URL:           "https://" + srcUrl,
		SingleBranch:  true,
		ReferenceName: ref.Name(),
	}

	fmt.Println("cloning:", co.URL, ref)

	if strings.Contains(srcUrl, "github.com") && os.Getenv("GITHUB_TOKEN") != "" {
		co.Auth = &http.BasicAuth{
			Username: "github-token", // yes, this can be anything except an empty string
			Password: os.Getenv("GITHUB_TOKEN"),
		}
	}

	// Clones the repository into the worktree (fs) and storer all the .git
	// content into the storer
	st := memory.NewStorage()
	fs := memfs.New()
	r, err := gogit.Clone(st, fs, co)
	if err != nil {
		return nil, err
	}

	return &GitRepo{
		Store: st,
		FS:    fs,
		Repo:  r,
	}, nil
}
