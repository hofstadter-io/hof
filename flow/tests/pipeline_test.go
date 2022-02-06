package flow_test

import (
	"testing"

	"github.com/hofstadter-io/hof/lib/yagu"
	"github.com/hofstadter-io/hof/script/runtime"
)

func TestFlow(t *testing.T) {
	yagu.Mkdir(".workdir/tests")
	runtime.Run(t, runtime.Params{
		Dir:         "testdata",
		Glob:        "nested*.txt",
		WorkdirRoot: ".workdir/tests",
	})
}
