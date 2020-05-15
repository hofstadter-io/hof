package cache

import (
	"os"
	"strings"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
)

func Load(lang, mod, ver string) (FS billy.Filesystem, err error) {
	flds := strings.Split(mod, "/")
	remote := flds[0]
	owner := flds[1]
	repo := flds[2]
	tag := ver

	dir := Outdir(lang, remote, owner, repo, tag)

	// fmt.Println("Cache Load:", dir)

	_, err = os.Lstat(dir)
	if err != nil {
		if _, ok := err.(*os.PathError); !ok && err.Error() != "file does not exist" && err.Error() != "no such file or directory" {
			return nil, err
		}
	}

	FS = osfs.New(dir)

	return FS, nil
}
