package mod_test

import (
	"os"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"

	"github.com/hofstadter-io/hof/lib/yagu"
)

func envSetup(env *testscript.Env) error {
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		env.Vars = append(env.Vars, "GITHUB_TOKEN="+token)
	}
	env.Vars = append(env.Vars, "HOF_TELEMETRY_DISABLED=1")
	return nil
}

func TestModTests(t *testing.T) {
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
