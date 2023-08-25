package cuecmd_test

import (
	"os"
	"testing"

	"github.com/hofstadter-io/hof/lib/yagu"
	"github.com/hofstadter-io/hof/script/runtime"
)

func envSetup(env *runtime.Env) error {
	env.Vars = append(env.Vars, "HOF_TELEMETRY_DISABLED=1")
	return nil
}

func setupWorkdir(dir string) {
	os.RemoveAll(dir)
	yagu.Mkdir(dir)
}

func TestDef(t *testing.T) {
	d := ".workdir/def"
	setupWorkdir(d)
	runtime.Run(t, runtime.Params{
		Setup:       envSetup,
		Dir:         "testdata",
		Glob:        "def_*.txt",
		WorkdirRoot: d,
	})
}

func TestEval(t *testing.T) {
	d := ".workdir/eval"
	setupWorkdir(d)
	runtime.Run(t, runtime.Params{
		Setup:       envSetup,
		Dir:         "testdata",
		Glob:        "eval_*.txt",
		WorkdirRoot: d,
	})
}

func TestExport(t *testing.T) {
	d := ".workdir/export"
	setupWorkdir(d)
	runtime.Run(t, runtime.Params{
		Setup:       envSetup,
		Dir:         "testdata",
		Glob:        "export_*.txt",
		WorkdirRoot: d,
	})
}

func TestVet(t *testing.T) {
	d := ".workdir/vet"
	setupWorkdir(d)
	runtime.Run(t, runtime.Params{
		Setup:       envSetup,
		Dir:         "testdata",
		Glob:        "vet_*.txt",
		WorkdirRoot: d,
	})
}

