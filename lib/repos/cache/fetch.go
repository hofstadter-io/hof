package cache

import (
	"fmt"
	"path"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	gogit "github.com/go-git/go-git/v5"

	"github.com/hofstadter-io/hof/lib/repos/git"
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

	remote, owner, repo := utils.ParseModURL(rpath)
	dir := SourceOutdir(remote, owner, repo)
	url := path.Join(remote, owner, repo)

	// only fetch if we haven't already this run
	_, ok := syncedRepos.Load(url)
	if !ok {

		err := git.SyncSource(dir, remote, owner, repo, ver)
		if err != nil {
			return nil, err
		}

		syncedRepos.Store(url, true)
	}

	FS := osfs.New(dir)

	return FS, nil
}
