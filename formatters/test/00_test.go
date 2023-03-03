package fmt_test

import (
	"os"
	"testing"

	"github.com/hofstadter-io/hof/lib/yagu"
	"github.com/hofstadter-io/hof/script/runtime"
)

func envSetup(env *runtime.Env) error {
	env.Vars = append(env.Vars, "HOF_TELEMETRY_DISABLED=1")
	env.Vars = append(env.Vars, "DOCKER_HOST=" + os.Getenv("DOCKER_HOST"))
	env.Vars = append(env.Vars, "GITHUB_TOKEN=" + os.Getenv("GITHUB_TOKEN"))
	env.Vars = append(env.Vars, "HOF_FMT_VERSION=" + os.Getenv("HOF_FMT_VERSION"))
	env.Vars = append(env.Vars, "HOF_FMT_DEBUG=" + os.Getenv("HOF_FMT_DEBUG"))
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
