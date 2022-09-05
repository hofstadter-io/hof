package flow_test

import (
	"testing"

	"github.com/hofstadter-io/hof/lib/yagu"
	"github.com/hofstadter-io/hof/script/runtime"
)

func envSetup(env *runtime.Env) error {
	env.Vars = append(env.Vars, "HOF_TELEMETRY_DISABLED=1")
	return nil
}

func TestAPIFlow(t *testing.T) {
	yagu.Mkdir(".workdir/tasks/api")
	runtime.Run(t, runtime.Params{
		Dir:         "testdata/tasks/api",
		Glob:        "*.txt",
		WorkdirRoot: ".workdir/tasks/api",
		Setup:       envSetup,
	})
}

func TestGenFlow(t *testing.T) {
	yagu.Mkdir(".workdir/tasks/gen")
	runtime.Run(t, runtime.Params{
		Dir:         "testdata/tasks/gen",
		Glob:        "*.txt",
		WorkdirRoot: ".workdir/tasks/gen",
		Setup:       envSetup,
	})
}

func TestKVFlow(t *testing.T) {
	yagu.Mkdir(".workdir/tasks/kv")
	runtime.Run(t, runtime.Params{
		Dir:         "testdata/tasks/kv",
		Glob:        "*.txt",
		WorkdirRoot: ".workdir/tasks/kv",
		Setup:       envSetup,
	})
}

func TestOSFlow(t *testing.T) {
	yagu.Mkdir(".workdir/tasks/os")
	runtime.Run(t, runtime.Params{
		Dir:         "testdata/tasks/os",
		Glob:        "*.txt",
		WorkdirRoot: ".workdir/tasks/os",
		Setup:       envSetup,
	})
}

func TestStFlow(t *testing.T) {
	yagu.Mkdir(".workdir/tasks/st")
	runtime.Run(t, runtime.Params{
		Dir:         "testdata/tasks/st",
		Glob:        "*.txt",
		WorkdirRoot: ".workdir/tasks/st",
		Setup:       envSetup,
	})
}
