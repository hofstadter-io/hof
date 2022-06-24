package flow_test

import (
	"testing"

	"github.com/hofstadter-io/hof/lib/yagu"
	"github.com/hofstadter-io/hof/script/runtime"
)

func TestRender(t *testing.T) {
	yagu.Mkdir(".workdir/render")
	runtime.Run(t, runtime.Params{
		Dir:         "./",
		Glob:        "*.txt",
		WorkdirRoot: ".workdir/render",
	})
}

