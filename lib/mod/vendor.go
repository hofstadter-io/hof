package mod

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/repos/cache"
	"github.com/hofstadter-io/hof/lib/repos/utils"
	"github.com/hofstadter-io/hof/lib/yagu"
)

func Vendor(rflags flags.RootPflagpole, mflags flags.ModPflagpole) (error) {
	upgradeHofMods()

	cm, err := loadRootMod()
	if err != nil {
		return err
	}

	return cm.Vendor(false, rflags.Verbosity)
}

func (cm *CueMod) Vendor(link bool, verbosity int) (err error) {

	// cleanup before remaking
	pkgDir := filepath.Join(cm.Basedir, "cue.mod", "pkg")
	err = os.RemoveAll(pkgDir)
	if err != nil {
		return err
	}

	err = os.Mkdir(pkgDir, 0755)
	if err != nil {
		return err
	}

	written := make(map[string]string)
	// we symlink deps in the following preference order
	// 1. replace
	// 2. require
	// 2. indirect

	process := func (orig, path, ver string) error {
		// check
		_, ok := written[orig]
		if ok {
			if verbosity > 0 {
				fmt.Println("warning, replaced twice:", orig)
			}
			return nil
		}

		src, dst := "", filepath.Join(cm.Basedir, "cue.mod", "pkg", orig)

		if ver == "" {
			src = filepath.Join(cm.Basedir, path)
		} else {
			remote, owner, repo := utils.ParseModURL(path)
			src = cache.Outdir(remote, owner, repo, ver)
		}

		err := linkOrCopy(link, src, dst)
		if err != nil {
			return err
		}

		written[orig] = src
		return nil
	}

	// for each replace
	for orig, repl := range cm.Replace {
		err := process(orig, repl.Path, repl.Version)
		if err != nil {
			return err
		}
	}

	for path, ver := range cm.Require {
		if _, ok := cm.Replace[path]; ok {
			continue
		}
		err := process(path, path, ver)
		if err != nil {
			return err
		}
	}

	for path, ver := range cm.Indirect {
		if _, ok := cm.Replace[path]; ok {
			continue
		}
		err := process(path, path, ver)
		if err != nil {
			return err
		}
	}

	return nil
}

func linkOrCopy(link bool, src, dst string) (err error) {
	if link {
		err = os.MkdirAll(filepath.Dir(dst), 0755)
		if err != nil {
			return err
		}

		err = os.Symlink(src, dst)
		if err != nil {
			return err
		}
	} else {
		err = yagu.CopyDir(src, dst)
	}

	return nil
}
