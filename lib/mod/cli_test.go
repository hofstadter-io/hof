package mod_test

import (
	"testing"

	"github.com/rogpeppe/go-internal/testscript"

	"github.com/hofstadter-io/hof/lib/yagu"
)

func envSetup(env *testscript.Env) error {

	env.Vars = append(env.Vars, "HOF_TELEMETRY_DISABLED=1")

	return nil
}

func TestMod(t *testing.T) {

	yagu.Mkdir(".workdir/tests")

	testscript.Run(t, testscript.Params{
		Setup: envSetup,
		Dir: "testdata",
		WorkdirRoot: ".workdir/tests",
	})
}

func TestModBugs(t *testing.T) {

	yagu.Mkdir(".workdir/bugs")

	testscript.Run(t, testscript.Params{
		Setup: envSetup,
		Dir: "testdata/bugs",
		WorkdirRoot: ".workdir/bugs",
	})
}
