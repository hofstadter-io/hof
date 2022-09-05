package script

import (
	"github.com/hofstadter-io/hof/script/runtime"
)

func Run(glob string) error {
	r := runtime.Runner{
		// LogLevel: flags.RootVerbosePflag,
		LogLevel: "yes please",
	}

	p := runtime.Params{
		Mode:        "run",
		Setup:       envSetup,
		Dir:         ".",
		Glob:        glob,
		WorkdirRoot: ".",
		TestWork:    true,
	}

	runtime.RunT(r, p)

	// TODO check output / status?

	return nil
}
