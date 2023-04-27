package cache

import (
	"context"
	"fmt"
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

func FetchRepoSource(mod, ver string) (billy.Filesystem, error) {
	if debug {
		fmt.Println("cache.FetchRepoSource:", mod)
	}

	rmt, err := remote.Parse(mod)
	if err != nil {
		return nil, fmt.Errorf("remote parse: %w", err)
	}

	dir := SourceOutdirParts(rmt.Host, rmt.Owner, rmt.Name)

	// TODO:
	//   * Use a passed-in context.
	//   * Choose a better timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// only fetch if we haven't already this run
	if _, ok := syncedRepos.Load(mod); ok {
		if err := rmt.Pull(ctx, dir, ver); err != nil {
			return nil, fmt.Errorf("remote pull: %w", err)
		}

		syncedRepos.Store(mod, true)
	}

	return osfs.New(dir), nil
}
