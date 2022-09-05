package fmt_test

import (
	"testing"

	"github.com/hofstadter-io/hof/lib/yagu"
	"github.com/hofstadter-io/hof/script/runtime"
)

func envSetup(env *runtime.Env) error {
	env.Vars = append(env.Vars, "HOF_TELEMETRY_DISABLED=1")
	return nil
}

func TestFormatters(t *testing.T) {
	yagu.Mkdir(".workdir")
	runtime.Run(t, runtime.Params{
		Dir:         "./",
		Glob:        "*.txt",
		WorkdirRoot: ".workdir",
		Setup:       envSetup,
	})
}


