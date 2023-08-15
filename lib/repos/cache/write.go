package cache

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"

	"github.com/hofstadter-io/hof/lib/repos/utils"
	"github.com/hofstadter-io/hof/lib/yagu"
)

func CleanCache() error {
	return os.RemoveAll(cacheBaseDir)
}

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
	return SourceOutdirParts(remote, owner, repo)
}

func SourceOutdirParts(parts ...string) string {
	return filepath.Join(srcBaseDir, filepath.Join(parts...))
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

	err = yagu.BillyWriteDirToOS(outdir, "/", FS)
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

func CopyRepoTag(path, ver string) (string, error) {
	if debug {
		fmt.Println("cache.CopyRepoTag:", path, ver)
	}
	remote, owner, repo := utils.ParseModURL(path)
	dir := SourceOutdir(remote, owner, repo)

	// open git source
	R, err := OpenRepoSource(path)
	if err != nil {
		return ver, fmt.Errorf("(crt) open source error: %w for %s@%s", err, path, ver)
	}

	// get working tree
	wt, err := R.Worktree()
	if err != nil {
		return ver, fmt.Errorf("(crt) worktree error: %w for %s@%s", err, path, ver)
	}

	lver := ver
	parts := strings.Split(lver, "-")
	if strings.HasPrefix(lver, "v0.0.0-") {
		lver = strings.Join(parts[2:], "-")
	}

	// fmt.Println("PVL", path, ver, lver)

	// checkout tag
	err = wt.Checkout(&gogit.CheckoutOptions{
		Branch: plumbing.NewTagReferenceName(lver),
		Force:  true,
	})
	if err != nil {
		// fmt.Printf("(crt) -- checkout error: %v for %s@%s\n", err, path, ver)

		// err = fmt.Errorf("(crt) checkout error: %w for %s@%s", err, path, ver)
		// try branch
		err = wt.Checkout(&gogit.CheckoutOptions{
			Branch: plumbing.NewRemoteReferenceName("origin", lver),
			Force:  true,
		})

		if err != nil {
			err = wt.Checkout(&gogit.CheckoutOptions{
				Hash:  plumbing.NewHash(lver),
				Force: true,
			})
			if err != nil {
				return ver, fmt.Errorf("(crt) checkout error: unable to find version %q for module %q: %w", ver, path, err)
			}

		} else {
			h, err := R.Head()
			lver = strings.Join(append(parts[:2], h.Hash().String()), "-")
			fmt.Println("Checking out branch:", path, ver, lver, err)
			ver = lver
		}
	}

	// copy
	FS := osfs.New(dir)
	err = Write(remote, owner, repo, ver, FS)
	if err != nil {
		return ver, fmt.Errorf("(crt) writing error: %w for %s@%s", err, path, ver)
	}

	return ver, nil
}
