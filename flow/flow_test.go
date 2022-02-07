package flow_test

import (
	"testing"

	"github.com/hofstadter-io/hof/lib/yagu"
	"github.com/hofstadter-io/hof/script/runtime"
)

func TestMainFlow(t *testing.T) {
	yagu.Mkdir(".workdir/main")
	runtime.Run(t, runtime.Params{
		Dir:         "testdata",
		Glob:        "*.txt",
		WorkdirRoot: ".workdir/main",
	})
}

func TestAPIFlow(t *testing.T) {
	yagu.Mkdir(".workdir/api")
	runtime.Run(t, runtime.Params{
		Dir:         "testdata/api",
		Glob:        "*.txt",
		WorkdirRoot: ".workdir/api",
	})
}

func TestOSFlow(t *testing.T) {
	yagu.Mkdir(".workdir/os")
	runtime.Run(t, runtime.Params{
		Dir:         "testdata/os",
		Glob:        "*.txt",
		WorkdirRoot: ".workdir/os",
	})
}

func TestStFlow(t *testing.T) {
	yagu.Mkdir(".workdir/st")
	runtime.Run(t, runtime.Params{
		Dir:         "testdata/st",
		Glob:        "*.txt",
		WorkdirRoot: ".workdir/st",
	})
}
