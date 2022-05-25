package script

import (
	"io/ioutil"
	"strings"

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

func envSetup(env *runtime.Env) error {

	// .env can contain lines of ENV=VAR
	content, err := ioutil.ReadFile(".env")
	if err != nil {
		// ignore errors, as the file likely doesn't exist
		return nil
	}

	for _, line := range strings.Split(string(content), "\n") {
		if strings.Contains(line, "=") {
			if line[0:1] == "#" {
				continue
			}
			env.Vars = append(env.Vars, line)
		}
	}

	return nil
}
