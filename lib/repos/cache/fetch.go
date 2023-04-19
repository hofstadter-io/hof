package cache

import (
	"context"
	"fmt"
	"path"
	"time"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	gogit "github.com/go-git/go-git/v5"

	"github.com/hofstadter-io/hof/lib/repos/remote"
	"github.com/hofstadter-io/hof/lib/repos/utils"
)

func OpenRepoSource(path string) (*gogit.Repository, error) {
	if debug {
		fmt.Println("cache.OpenRepoSource:", path)
	}

	remote, owner, repo := utils.ParseModURL(path)
	dir := SourceOutdir(remote, owner, repo)
	return gogit.PlainOpen(dir)
}

func FetchRepoSource(rpath, ver string) (billy.Filesystem, error) {
	if debug {
		fmt.Println("cache.FetchRepoSource:", rpath)
	}

	rmt, err := remote.Parse(rpath)
	if err != nil {
		return nil, fmt.Errorf("remote parse: %w", err)
	}

	var (
		p   = rmt.Parts()
		dir = SourceOutdirParts(p...)
		url = path.Join(p...)
	)

	// TODO:
	//   * Use a passed-in context.
	//   * Choose a better timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// only fetch if we haven't already this run
	if _, ok := syncedRepos.Load(url); ok {
		if err := rmt.Pull(ctx, remote.LocalDir(dir), remote.Version(ver)); err != nil {
			return nil, fmt.Errorf("remote pull: %w", err)
		}

		syncedRepos.Store(url, true)
	}

	return osfs.New(dir), nil
}
