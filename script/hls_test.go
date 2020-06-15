package script_test

import (
	"testing"

	"github.com/hofstadter-io/hof/script/runtime"
)

func TestHLS(t *testing.T) {
	runtime.Run(t, runtime.Params{
		Dir: "tests/browser",
		Glob: "*.hls",
	})

	runtime.Run(t, runtime.Params{
		Dir: "tests/http",
		Glob: "*.hls",
	})
}

