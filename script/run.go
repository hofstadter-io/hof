package script

import (
	"github.com/hofstadter-io/script/runtime"
	"github.com/hofstadter-io/script/util"
)

func Run(glob string) error {
	r := util.Runner{
		// LogLevel: flags.RootVerbosePflag,
		LogLevel: "yes please",
	}

	p := runtime.Params{
		Mode: "run",
		Setup: envSetup,
		Dir: ".",
		Glob: glob,
		WorkdirRoot: ".",
		TestWork: true,
	}

	runtime.RunT(r, p)

	// TODO check output / status?

	return nil
}

