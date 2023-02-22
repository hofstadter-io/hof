package mod

import (
	"github.com/go-git/go-billy/v5/osfs"

	"github.com/hofstadter-io/hof/lib/cuetils"
)

func loadRootMod() (*CueMod, error) {
	basedir, err := cuetils.FindModuleAbsPath("")
	if err != nil {
		return nil, err
	}

	FS := osfs.New(basedir)

	return ReadModule(basedir, FS)
}
