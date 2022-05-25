package runtime

import (
	"io/ioutil"

	"github.com/hofstadter-io/hof/lib/gotils/txtar"
)

// unquote unquotes files.
func (ts *Script) CmdUnquote(neg int, args []string) {
	if neg != 0 {
		ts.Fatalf("unsupported: !? unquote")
	}
	for _, arg := range args {
		file := ts.MkAbs(arg)
		data, err := ioutil.ReadFile(file)
		ts.Check(err)
		data, err = txtar.Unquote(data)
		ts.Check(err)
		err = ioutil.WriteFile(file, data, 0666)
		ts.Check(err)
	}
}
