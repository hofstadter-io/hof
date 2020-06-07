package cmd_test

import (
	"testing"

	"github.com/hofstadter-io/hof/lib/gotils/testscript"
	"github.com/hofstadter-io/hof/lib/yagu"

	"github.com/hofstadter-io/hof/cmd/hof/cmd"
)

func TestScriptBisectCliTests(t *testing.T) {
	// setup some directories

	dir := "bisect"

	workdir := ".workdir/cli/" + dir
	yagu.Mkdir(workdir)

	testscript.Run(t, testscript.Params{
		Setup: func(env *testscript.Env) error {
			// add any environment variables for your tests here

			env.Vars = append(env.Vars, "HOF_TELEMETRY_DISABLED=1")

			return nil
		},
		Funcs: map[string]func(ts *testscript.TestScript, args []string) error{
			"__hof": cmd.CallTS,
		},
		Dir:         "testscripts/cli/bisect",
		WorkdirRoot: workdir,
	})
}
