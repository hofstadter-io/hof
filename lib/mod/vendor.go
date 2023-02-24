package mod

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/repos/cache"
	"github.com/hofstadter-io/hof/lib/repos/utils"
	"github.com/hofstadter-io/hof/lib/yagu"
)

func Vendor(rflags flags.RootPflagpole) (error) {
	upgradeHofMods()

	cm, err := loadRootMod()
	if err != nil {
		return err
	}

	if rflags.Verbosity > 0 {
		fmt.Println("vendoring deps for:", cm.Module)
	}

	err = cm.ensureCached()
	if err != nil {
		return err
	}

	return cm.Vendor("copy", rflags.Verbosity)
}

func (cm *CueMod) Vendor(method string, verbosity int) (err error) {
	pkgDir := filepath.Join(cm.Basedir, "cue.mod", "pkg")

	// dynamically determine
	if method == "" {

		for path, _ := range cm.Require {
			info, err := os.Lstat(filepath.Join(pkgDir, path))
			if err == nil {
				if info.Mode() & fs.ModeSymlink == fs.ModeSymlink {
					method = "link"
					break
				} else if info.IsDir() {
					method = "copy"
					break
				}
			}
		}

		if method == "" {
			for path, _ := range cm.Indirect {
				if _, ok := cm.Replace[path]; ok {
					info, err := os.Lstat(filepath.Join(pkgDir, path))
					if err == nil {
						if info.Mode() & fs.ModeSymlink == fs.ModeSymlink {
							method = "link"
							break
						} else if info.IsDir() {
							method = "copy"
							break
						}
					}
				}
			}
		}

		if method == "" {
			method = "link"
		}
	}

	// cleanup before remaking
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
			src = cache.ModuleOutdir(remote, owner, repo, ver)
		}

		err := linkOrCopy(method, src, dst)
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

func linkOrCopy(method string, src, dst string) (err error) {
	if method == "link" {
		err = os.MkdirAll(filepath.Dir(dst), 0755)
		if err != nil {
			return err
		}
		err = os.Symlink(src, dst)

	} else if method == "copy" {
		err = yagu.CopyDir(src, dst)
	} else {
		panic("uknown mod install method: " + method)
	}

	return err
}
