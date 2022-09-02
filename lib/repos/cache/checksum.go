package cache

import (
	"os"

	"golang.org/x/mod/sumdb/dirhash"

	"github.com/hofstadter-io/hof/lib/repos/utils"
)

func Checksum(mod, ver string) (string, error) {
	remote, owner, repo := utils.ParseModURL(mod)
	tag := ver

	dir := Outdir(remote, owner, repo, tag)
	// fmt.Println("Cache Checksum:", dir)

	_, err := os.Lstat(dir)
	if err != nil {
		return "", err
	}

	h, err := dirhash.HashDir(dir, mod, dirhash.Hash1)

	return h, err
}
