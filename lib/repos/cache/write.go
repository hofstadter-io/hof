package cache

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"

	"github.com/hofstadter-io/hof/lib/repos/utils"
	"github.com/hofstadter-io/hof/lib/yagu"
)

func ModuleOutdir(remote, owner, repo, tag string) string {
	outdir := filepath.Join(
		modBaseDir,
		remote,
		owner,
		repo+"@"+tag,
	)
	return outdir
}

func SourceOutdir(remote, owner, repo string) string {
	outdir := filepath.Join(
		srcBaseDir,
		remote,
		owner,
		repo,
	)
	return outdir
}

func Write(remote, owner, repo, tag string, FS billy.Filesystem) error {
	if debug {
		fmt.Println("cache.Write:", remote, owner, repo, tag)
	}
	outdir := SourceOutdir(remote, owner, repo)
	if tag != "" {
		outdir = ModuleOutdir(remote, owner, repo, tag)
	}
	err := yagu.Mkdir(outdir)
	if err != nil {
		return err
	}

	err =  yagu.BillyWriteDirToOS(outdir, "/", FS)
	if err != nil {
		return err
	}

	// hacky, but remove .git dir if tag is set
	if tag != "" {
		err = os.RemoveAll(filepath.Join(outdir, ".git"))
		if err != nil {
			return err
		}
	}

	return nil
}

func CopyRepoTag(path, ver string) error {
	if debug {
		fmt.Println("cache.CopyRepoTag:", path, ver)
	}
	remote, owner, repo := utils.ParseModURL(path)
	dir := SourceOutdir(remote, owner, repo)

	// open git source
	R, err := OpenRepoSource(path)
	if err != nil {
		return fmt.Errorf("(crt) open source error: %w for %s@%s", err, path, ver)
	}

	// get workign tree
	wt, err := R.Worktree()
	if err != nil {
		return fmt.Errorf("(crt) worktree error: %w for %s@%s", err, path, ver)
	}

	// checkout tag
	err = wt.Checkout(&gogit.CheckoutOptions{
		Branch: plumbing.ReferenceName("refs/tags/" + ver),
		Force: true,
	})
	if err != nil {
		return fmt.Errorf("(crt) checkout error: %w for %s@%s", err, path, ver)
	}

	// copy
	FS := osfs.New(dir)
	err = Write(remote, owner, repo, ver, FS)
	if err != nil {
		return fmt.Errorf("(crt) writing error: %w for %s@%s", err, path, ver)
	}

	return nil
}

