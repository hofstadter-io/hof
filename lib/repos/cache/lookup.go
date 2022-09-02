package cache

import (
	"os"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"

	"github.com/hofstadter-io/hof/lib/repos/utils"
)

func Load(url, ver string) (FS billy.Filesystem, err error) {
	remote, owner, repo := utils.ParseModURL(url)
	dir := Outdir(remote, owner, repo, ver)

	// check for directory
	_, err = os.Lstat(dir)
	if err != nil {
		if _, ok := err.(*os.PathError); !ok && err.Error() != "file does not exist" && err.Error() != "no such file or directory" {
			return nil, err
		}
	}

	// load into FS
	FS = osfs.New(dir)

	return FS, nil
}
