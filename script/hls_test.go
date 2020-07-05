package script_test

import (
	"testing"

	"github.com/hofstadter-io/hof/script/runtime"
)

func TestScriptBrowser(t *testing.T) {
	runtime.RunTester(t, runtime.Params{
		Dir: "tests/browser",
		Glob: "*.hls",
	})
}

func TestScriptCmds(t *testing.T) {
	runtime.RunTester(t, runtime.Params{
		Dir: "tests/cmds",
		Glob: "*.hls",
	})
}

func TestScriptHTTP(t *testing.T) {
	runtime.RunTester(t, runtime.Params{
		Dir: "tests/http",
		Glob: "*.hls",
	})
}

