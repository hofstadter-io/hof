package playground

import (
	"fmt"
	"cuelang.org/go/cue"

	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/components/cue/helpers"
)

func (C *Playground) GetConnValue() cue.Value {
	tui.Log("trace", fmt.Sprintf("Play.GetConnValue from: %s/%s", C.Id(), C.Name()))

	return C.final.viewer.GetConnValue()
}

func (C *Playground) GetConnValueExpr(expr string) func () cue.Value {
	tui.Log("trace", fmt.Sprintf("Play.GetConnValueExpr from: %s/%s %s", C.Id(), C.Name(), expr))
	p := cue.ParsePath(expr)

	return func() cue.Value {
		return C.GetConnValue().LookupPath(p)
	}

}

func (C *Playground) SetConnection(args []string, valueGetter func() cue.Value) {
	c := C.scope.config

	c.Source = helpers.EvalConn
	c.ConnGetter = valueGetter
	c.Args = args

	C.scope.config = c
	C.scope.viewer.SetSourceConfig(c)

	C.scope.viewer.SetTitle(C.scope.viewer.BuildStatusString())


	C.useScope = true
	C.seeScope = true

	C.Rebuild(true)
}

