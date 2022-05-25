package st

import (
	"cuelang.org/go/cue"

	hofcontext "github.com/hofstadter-io/hof/flow/context"
	"github.com/hofstadter-io/hof/lib/structural"
)

type Replace struct{}

func NewReplace(val cue.Value) (hofcontext.Runner, error) {
	return &Replace{}, nil
}

// Tasks must implement a Run func, this is where we execute our task
func (T *Replace) Run(ctx *hofcontext.Context) (interface{}, error) {
	ctx.CUELock.Lock()
	defer ctx.CUELock.Unlock()

	v := ctx.Value

	x := v.LookupPath(cue.ParsePath("val"))
	repl := v.LookupPath(cue.ParsePath("replace"))

	r, err := structural.ReplaceValue(repl, x, nil)
	if err != nil {
		return nil, err
	}

	return v.FillPath(cue.ParsePath("out"), r), nil
}
