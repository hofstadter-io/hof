package cache

import (
	"fmt"
	"path/filepath"

	"github.com/go-git/go-billy/v5"

	"github.com/hofstadter-io/hof/lib/yagu"
)

func Outdir(lang, remote, owner, repo, tag string) string {
	outdir := filepath.Join(
		LocalCacheBaseDir,
		"mod",
		lang,
		remote,
		owner,
		repo + "@" + tag,
	)
	return outdir
}

func Write(lang, remote, owner, repo, tag string, FS billy.Filesystem) error {
	fmt.Printf("Saving %s mod %s/%s/%s@%s\n", lang, remote, owner, repo, tag)
	outdir := Outdir(lang, remote, owner, repo, tag)
	err := yagu.Mkdir(outdir)
	if err != nil {
		return err
	}
	return yagu.BillyWriteDirToOS(outdir, "/", FS)
}
