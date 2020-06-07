package cmd_test

import (
	"testing"

	"github.com/hofstadter-io/hof/lib/yagu"
	"github.com/hofstadter-io/hof/script"

	"github.com/hofstadter-io/hof/cmd/hof/cmd"
)

func init() {
	// ensure our root command is setup
	cmd.RootInit()
}

func TestScriptRootCliTests(t *testing.T) {
	// setup some directories
	workdir := ".workdir/cli/root"
	yagu.Mkdir(workdir)

	script.Run(t, script.Params{
		Setup: func(env *script.Env) error {
			// add any environment variables for your tests here

			env.Vars = append(env.Vars, "HOF_TELEMETRY_DISABLED=1")

			return nil
		},
		Funcs: map[string]func(ts *script.Script, args []string) error{
			"__hof": cmd.CallTS,
		},
		Dir:         "hls/cli/root",
		WorkdirRoot: workdir,
	})
}
