package mod

import (
	"github.com/go-git/go-billy/v5/osfs"

	"github.com/hofstadter-io/hof/lib/cuetils"
	"github.com/hofstadter-io/hof/lib/repos/cache"
)

func loadRootMod() (*CueMod, error) {
	basedir, err := cuetils.FindModuleAbsPath("")
	if err != nil {
		return nil, err
	}

	FS := osfs.New(basedir)

	return ReadModule(basedir, FS)
}

func (cm *CueMod) ensureCached() error {
	for path, ver := range cm.Require {
		_, err := cache.Load(path, ver)
		if err != nil {
			return err
		}
	}
	for path, ver := range cm.Indirect {
		_, err := cache.Load(path, ver)
		if err != nil {
			return err
		}
	}
	return nil
}
