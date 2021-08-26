package mod_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hofstadter-io/hof/lib/yagu"
	"github.com/hofstadter-io/hof/script/runtime"
)

func envSetup(env *runtime.Env) error {
	vars := []string{
		"GITHUB_TOKEN",
		"GITLAB_TOKEN",
	}
	for _, v := range vars {
		if val := os.Getenv(v); val != "" {
			env.Vars = append(env.Vars, fmt.Sprintf("%s=%s", v, val))
		}
	}
	env.Vars = append(env.Vars, "HOF_TELEMETRY_DISABLED=1")
	return nil
}

func TestModTests(t *testing.T) {
	yagu.Mkdir(".workdir/tests")
	runtime.Run(t, runtime.Params{
		Setup: envSetup,
		Dir: "testdata",
		Glob: "*.txt",
		WorkdirRoot: ".workdir/tests",
	})
}

func TestModBugs(t *testing.T) {
	yagu.Mkdir(".workdir/bugs")
	runtime.Run(t, runtime.Params{
		Setup: envSetup,
		Dir: "testdata/bugs",
		Glob: "*.txt",
		WorkdirRoot: ".workdir/bugs",
	})
}
