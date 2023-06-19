package git

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"

	"github.com/go-git/go-billy/v5"
	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"

	"github.com/hofstadter-io/hof/lib/repos/utils"
	"github.com/hofstadter-io/hof/lib/yagu"
)

var debug = false

func IsNetworkReachable(ctx context.Context, mod string) (bool, error) {
	rem := gogit.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{"https://" + mod},
	})

	host, _, _ := utils.ParseModURL(mod)

	auth, err := getAuth(host, "", "")
	if err != nil {
		return false, fmt.Errorf("get auth: %w", err)
	}

	_, err = rem.ListContext(ctx, &gogit.ListOptions{
		Auth: auth,
	})

	// TODO: This isn't ideal. This could be a failure
	// due to bad credentials and it would be better
	// to test for that and prompt the user.
	return err == nil, nil
}

func SyncSource(dir, remote, owner, repo, ver string) error {
	url := path.Join(remote, owner, repo)
	if debug {
		fmt.Println("git.SyncSource:", dir, remote, owner, repo, ver, url)
	}
	_, err := os.Lstat(dir)
	// does not exist
	if err != nil {
		// make plain clone, first fetch
		_, err := PlainClone(dir, remote, owner, repo)
		if err != nil {
			return err
		}
	} else {

		R, err := gogit.PlainOpen(dir)
		if err != nil {
			return err
		}

		opts := &gogit.FetchOptions{
			// Depth: 1,
			Force: true,
			Tags:  gogit.AllTags,
		}
		err = authFetch(opts, remote, owner, repo)
		if err != nil {
			return err
		}

		if debug {
			fmt.Println("git.FetchRepo:", dir, url)
		}

		fmt.Println("sync'n:", url)
		err = R.Fetch(opts)
		if err != nil {
			if strings.Contains(err.Error(), "already up-to-date") {
				return nil
			}

			// fallback on an exec
			cmd := exec.Command("git", "fetch")
			cmd.Dir = dir

			cerr := cmd.Run()
			if cerr != nil {
				fmt.Printf("warn: while sync'n %s: %v\n", url, err, cerr)
			}

		}
	}

	return nil
}

func PlainClone(dir, remote, owner, repo string) (*gogit.Repository, error) {
	if debug {
		fmt.Println("git.PlainClone:", dir, remote, owner, repo)
	}
	srcRepo := path.Join(owner, repo)
	opts := &gogit.CloneOptions{
		URL: fmt.Sprintf("https://%s/%s", remote, srcRepo),
		// Depth: 1,
		Tags: gogit.AllTags,
	}

	err := authClone(opts, remote, owner, repo)
	if err != nil {
		return nil, err
	}

	fmt.Println("fetch'n:", path.Join(remote, owner, repo))
	R, err := gogit.PlainClone(dir, false, opts)
	if err != nil {
		return R, err
	}

	return R, nil
}

func Clone(FS billy.Filesystem, remote, owner, repo string) (*gogit.Repository, error) {
	if debug {
		fmt.Println("git.Clone:", remote, owner, repo)
	}
	srcRepo := path.Join(owner, repo)
	opts := &gogit.CloneOptions{
		URL:   fmt.Sprintf("https://%s/%s", remote, srcRepo),
		Depth: 1,
		Tags:  gogit.AllTags,
	}

	err := authClone(opts, remote, owner, repo)
	if err != nil {
		return nil, err
	}

	R, err := gogit.Clone(memory.NewStorage(), FS, opts)
	if err != nil {
		return R, err
	}

	return R, nil
}

func authClone(opts *gogit.CloneOptions, remote, owner, repo string) error {
	auth, err := getAuth(remote, owner, repo)
	if err != nil {
		return err
	}

	opts.Auth = auth
	return nil
}

func authFetch(opts *gogit.FetchOptions, remote, owner, repo string) error {
	auth, err := getAuth(remote, owner, repo)
	if err != nil {
		return err
	}

	opts.Auth = auth
	return nil
}

var authMap sync.Map

func getAuth(remote, owner, repo string) (auth transport.AuthMethod, err error) {
	fmt.Println("getAuth", remote, owner, repo)
	// cached auth
	val, ok := authMap.Load(remote)
	if ok {
		if debug {
			fmt.Println("found auth:", remote, val)
		}
		return val.(transport.AuthMethod), nil
	}

	// lookup auth
	if netrc, err := yagu.NetrcCredentials(remote); err == nil {
		auth = &http.BasicAuth{
			Username: netrc.Login,
			Password: netrc.Password,
		}
	} else if strings.Contains(remote, "github.com") && os.Getenv("GITHUB_TOKEN") != "" {
		auth = &http.BasicAuth{
			Username: "github-token", // yes, this can be anything except an empty string
			Password: os.Getenv("GITHUB_TOKEN"),
		}
	} else if strings.Contains(remote, "gitlab.com") && os.Getenv("GITLAB_TOKEN") != "" {
		auth = &http.BasicAuth{
			Username: "gitlab-token", // yes, this can be anything except an empty string
			Password: os.Getenv("GITLAB_TOKEN"),
		}
	} else if strings.Contains(remote, "bitbucket.org") {
		if os.Getenv("BITBUCKET_PASSWORD") != "" {
			auth = &http.BasicAuth{
				Username: os.Getenv("BITBUCKET_USERNAME"), // yes, this can be anything except an empty string
				Password: os.Getenv("BITBUCKET_PASSWORD"),
			}
		} else if os.Getenv("BITBUCKET_TOKEN") != "" {
			auth = &http.BasicAuth{
				Username: "bitbucket-token", // yes, this can be anything except an empty string
				Password: os.Getenv("BITBUCKET_TOKEN"),
			}
		}
	//} else if ssh, err := yagu.SSHCredentials(remote); err == nil {
	//  auth = ssh.Keys
	}

	// no auth found, so don't return any
	if auth == nil {
		return nil, nil
	}
	if debug {
		fmt.Println("cache auth:", remote, auth)
	}
	authMap.Store(remote, auth)
	return auth, nil
}
