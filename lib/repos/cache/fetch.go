package cache

import (
	"fmt"
	"os"

	"github.com/hofstadter-io/hof/lib/repos/fetch"
	"github.com/hofstadter-io/hof/lib/repos/utils"
)

func Fetch(url, ver string) error {
	remote, owner, repo := utils.ParseModURL(url)

	// check for in cache
	dir := Outdir(remote, owner, repo, ver)
	if _, err := os.Lstat(dir); err == nil {
		return nil
	}

	FS, err := fetch.FetchRepoToMem(url, ver)
	if err != nil {
		return err
	}

	if err := Write(remote, owner, repo, ver, FS); err != nil {
		return fmt.Errorf("While writing to cache\n%w\n", err)
	}

	return nil
}
