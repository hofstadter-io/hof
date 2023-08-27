package flow_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hofstadter-io/hof/lib/yagu"
	"github.com/hofstadter-io/hof/script/runtime"
)

func envSetup(env *runtime.Env) error {
	env.Vars = append(env.Vars, "HOF_TELEMETRY_DISABLED=1")

	vars := []string{
		"GITHUB_TOKEN",
		"HOF_FMT_VERSION",
		"DOCKER_HOST",
		"CONTAINERD_ADDRESS",
		"CONTAINERD_NAMESPACE",
	}

	for _,v := range vars {
		val := os.Getenv(v)
		jnd := fmt.Sprintf("%s=%s", v, val)
		env.Vars = append(env.Vars, jnd)
	}

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

func TestBulk(t *testing.T) {
	yagu.Mkdir(".workdir/bulk")
	runtime.Run(t, runtime.Params{
		Dir:         "testdata/bulk/",
		Glob:        "*.txt",
		WorkdirRoot: ".workdir/bulk/",
		Setup:       envSetup,
	})
}

