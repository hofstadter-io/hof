package render_test

import (
	"os"
	"testing"

	"github.com/hofstadter-io/hof/lib/yagu"
	"github.com/hofstadter-io/hof/script/runtime"
)

func envSetup(env *runtime.Env) error {
	env.Vars = append(env.Vars, "HOF_TELEMETRY_DISABLED=1")
	env.Vars = append(env.Vars, "GITHUB_TOKEN=" + os.Getenv("GITHUB_TOKEN"))
	return nil
}

func TestRender(t *testing.T) {
	yagu.Mkdir(".workdir/render")
	runtime.Run(t, runtime.Params{
		Dir:         "./",
		Glob:        "*.txt",
		WorkdirRoot: ".workdir/render",
		Setup:       envSetup,
	})
}

