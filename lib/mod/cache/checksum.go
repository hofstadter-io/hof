package cache

import (
	"os"

	"golang.org/x/mod/sumdb/dirhash"
)

func Checksum(lang, mod, ver string) (string, error) {
	remote, owner, repo := parseModURL(mod)
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
