package cmd

import (
	"fmt"
	"strings"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
)

func Get(args []string, rflags flags.RootPflagpole, cflags flags.EnvPflagpole) error {
	args, cueargs := splitArgs(args)

	R, err := prepRuntime(cueargs, rflags)
	if err != nil {
		return err
	}

	for _, e := range R.Envs {
		if len(args) > 0 {
			for _, a := range args {
				if strings.HasPrefix(e.Hof.Env.Name, a) {
					fmt.Printf("%s: %#+v\n", e.Hof.Env.Name, e.Value)
				}
			}
		}
	}

	return nil
}
