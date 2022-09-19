package datamodel_test

import (
	"testing"

	"github.com/hofstadter-io/hof/lib/yagu"
	"github.com/hofstadter-io/hof/script/runtime"
)

func envSetup(env *runtime.Env) error {
	env.Vars = append(env.Vars, "HOF_TELEMETRY_DISABLED=1")
	return nil
}

func TestDatamodel(t *testing.T) {
	yagu.Mkdir(".workdir")
	runtime.Run(t, runtime.Params{
		Dir:         "./testdata",
		Glob:        "*.txt",
		WorkdirRoot: ".workdir/testdata",
		Setup:       envSetup,
	})
}

func TestBugcases(t *testing.T) {
	yagu.Mkdir(".workdir")
	runtime.Run(t, runtime.Params{
		Dir:         "./bugcases",
		Glob:        "*.txt",
		WorkdirRoot: ".workdir/bugcases",
		Setup:       envSetup,
	})
}

