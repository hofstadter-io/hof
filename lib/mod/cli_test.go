package mod_test

import (
	"os"
	"testing"

	"github.com/hofstadter-io/hof/lib/yagu"
	"github.com/hofstadter-io/hof/script/runtime"
)

func envSetup(RT *runtime.Runtime) error {
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		RT.AddEnvVar("GITHUB_TOKEN", token)
	}
	RT.AddEnvVar("HOF_TELEMETRY_DISABLED", "1")
	return nil
}

func TestModTests(t *testing.T) {
	yagu.Mkdir(".workdir/tests")
	runtime.RunTester(t, runtime.Params{
		Setup: envSetup,
		Dir: "testdata",
		Glob: "*.hls",
		WorkdirRoot: ".workdir/tests",
	})
}

func TestModBugs(t *testing.T) {
	yagu.Mkdir(".workdir/bugs")
	runtime.RunTester(t, runtime.Params{
		Setup: envSetup,
		Dir: "testdata/bugs",
		Glob: "*.hls",
		WorkdirRoot: ".workdir/bugs",
	})
}
