package cache

import (
	"os"
	"strings"

	"golang.org/x/mod/sumdb/dirhash"
)

func Checksum(lang, mod, ver string) (string, error) {

	flds := strings.SplitN(mod, "/", 3)
	remote := flds[0]
	owner := flds[1]
	repo := flds[2]
	tag := ver

	dir := Outdir(lang, remote, owner, repo, tag)
	// fmt.Println("Cache Checksum:", dir)

	_, err := os.Lstat(dir)
	if err != nil {
		return "", err
	}

	h, err := dirhash.HashDir(dir, mod, dirhash.Hash1)

	return h, err
}
