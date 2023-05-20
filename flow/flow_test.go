package flow_test

import (
	"os"
	"testing"

	"github.com/hofstadter-io/hof/lib/yagu"
	"github.com/hofstadter-io/hof/script/runtime"
)

func envSetup(env *runtime.Env) error {
	env.Vars = append(env.Vars, "HOF_TELEMETRY_DISABLED=1")
	env.Vars = append(env.Vars, "GITHUB_TOKEN=" + os.Getenv("GITHUB_TOKEN"))
	env.Vars = append(env.Vars, "HOF_FMT_VERSION=" + os.Getenv("HOF_FMT_VERSION"))
	return nil
}

func doTaskTest(dir string, t *testing.T) {
	yagu.Mkdir(".workdir/tasks/" + dir)
	runtime.Run(t, runtime.Params{
		Dir:         "testdata/tasks/" + dir,
		Glob:        "*.txt",
		WorkdirRoot: ".workdir/tasks/" + dir,
		Setup:       envSetup,
	})
}

func TestAPIFlow(t *testing.T) {
	doTaskTest("api", t)
}

func TestGenFlow(t *testing.T) {
	doTaskTest("gen", t)
}

func TestHofFlow(t *testing.T) {
	doTaskTest("hof", t)
}

func TestKVFlow(t *testing.T) {
	doTaskTest("kv", t)
}

func TestOSFlow(t *testing.T) {
	doTaskTest("os", t)
}

func TestStFlow(t *testing.T) {
	doTaskTest("st", t)
}
