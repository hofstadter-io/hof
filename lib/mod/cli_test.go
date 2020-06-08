package mod_test

import (
	"os"
	"testing"

	"github.com/hofstadter-io/hof/lib/yagu"
	"github.com/hofstadter-io/hof/script"
)

func envSetup(env *script.Env) error {
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		env.Vars = append(env.Vars, "GITHUB_TOKEN="+token)
	}
	env.Vars = append(env.Vars, "HOF_TELEMETRY_DISABLED=1")
	return nil
}

func TestModTests(t *testing.T) {
	yagu.Mkdir(".workdir/tests")
	script.Run(t, script.Params{
		Setup: envSetup,
		Dir: "testdata",
		Glob: "*.txt",
		WorkdirRoot: ".workdir/tests",
	})
}

func TestModBugs(t *testing.T) {
	yagu.Mkdir(".workdir/bugs")
	script.Run(t, script.Params{
		Setup: envSetup,
		Dir: "testdata/bugs",
		Glob: "*.txt",
		WorkdirRoot: ".workdir/bugs",
	})
}
