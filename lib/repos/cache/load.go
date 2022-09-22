package cache

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"

	"github.com/hofstadter-io/hof/lib/repos/utils"
)

const CustomCacheBaseDirVar = "HOF_CACHE"
var cacheBaseDir  string

func init() {
	e := os.Getenv(CustomCacheBaseDirVar)
	if e != "" {
		cacheBaseDir = e
	} else {
		d, err := os.UserCacheDir()
		if err != nil {
			return
		}

		// save to hof dir for cache across projects
		cacheBaseDir = filepath.Join(d, "hof/mods")
	}
}

func SetBaseDir(basedir string) {
	cacheBaseDir = basedir
}

func Load(url, ver string) (_ billy.Filesystem, err error) {
	FS, err := Read(url, ver)
	if err == nil {
		return FS, nil
	}

	return Cache(url, ver)
}

func Read(url, ver string) (FS billy.Filesystem, err error) {
	remote, owner, repo := utils.ParseModURL(url)
	dir := Outdir(remote, owner, repo, ver)

	// check for existence
	_, err = os.Lstat(dir)
	if err != nil {
		return nil, err
	}

	// load into FS
	FS = osfs.New(dir)

	return FS, nil
}

func Cache(url, ver string) (billy.Filesystem, error) {
	remote, owner, repo := utils.ParseModURL(url)

	// check for in cache
	dir := Outdir(remote, owner, repo, ver)
	if _, err := os.Lstat(dir); err == nil {
		return nil, nil
	}

	FS, err := FetchRepoToMem(url, ver)
	if err != nil {
		return FS, err
	}

	if err := Write(remote, owner, repo, ver, FS); err != nil {
		return FS, fmt.Errorf("While writing to cache\n%w\n", err)
	}

	return FS, nil
}

