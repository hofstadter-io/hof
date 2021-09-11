package git

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-billy/v5/osfs"
	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"

	"github.com/hofstadter-io/hof/lib/yagu"
)

// TODO, this file has inconsistency of auth creds adding
// between the functions, this should be cleaned up
// taking note that we want to create more consistency
// across the various repo types

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

// FetchGit clone the repository inside FS.
// If private flag is set, it will look for netrc credentials, fallbacking to SSH
func Fetch(FS billy.Filesystem, remote, owner, repo, tag string, private bool) error {
	srcRepo := path.Join(owner, repo)
	gco := &gogit.CloneOptions{
		URL:   fmt.Sprintf("https://%s/%s", remote, srcRepo),
		Depth: 1,
	}

	if tag != "v0.0.0" {
		gco.ReferenceName = plumbing.NewTagReferenceName(tag)
		gco.SingleBranch = true
	}

	if private {
		fmt.Println("git.Fetch")
		if netrc, err := yagu.NetrcCredentials(remote); err == nil {
			fmt.Println("NetRC")
			gco.Auth = &http.BasicAuth{
				Username: netrc.Login,
				Password: netrc.Password,
			}
		} else if ssh, err := yagu.SSHCredentials(remote); err == nil {
			fmt.Println("SSHCreds")
			gco.Auth = ssh.Keys
			gco.URL = fmt.Sprintf("%s@%s:%s", ssh.User, remote, srcRepo)
			fmt.Println("URL:", gco.URL)
		} else {
			fmt.Println("NoAuth")
			gco.URL = fmt.Sprintf("%s@%s:%s", "git", remote, srcRepo)
		}
	}

	_, err := gogit.Clone(memory.NewStorage(), FS, gco)

	return err
}
