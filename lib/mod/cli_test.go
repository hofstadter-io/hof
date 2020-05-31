package mod_test

import (
	"testing"

	"github.com/rogpeppe/go-internal/testscript"

	"github.com/hofstadter-io/hof/lib/yagu"
)

func TestMod(t *testing.T) {

	yagu.Mkdir(".workdir/tests")

	testscript.Run(t, testscript.Params{
		Dir: "testdata",
		WorkdirRoot: ".workdir/tests",
	})
}

func TestModBugs(t *testing.T) {

	yagu.Mkdir(".workdir/bugs")

	testscript.Run(t, testscript.Params{
		Dir: "testdata/bugs",
		WorkdirRoot: ".workdir/bugs",
	})
}
