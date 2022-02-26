package st

import (
	"cuelang.org/go/cue"

	hofcontext "github.com/hofstadter-io/hof/flow/context"
	"github.com/hofstadter-io/hof/lib/structural"
)

type Diff struct {}

func NewDiff(val cue.Value) (hofcontext.Runner, error) {
  return &Diff{}, nil
}

// Tasks must implement a Run func, this is where we execute our task
func (T *Diff) Run(ctx *hofcontext.Context) (interface{}, error) {
  ctx.CUELock.Lock()
  defer ctx.CUELock.Unlock()

	v := ctx.Value

	o := v.LookupPath(cue.ParsePath("orig"))
	n := v.LookupPath(cue.ParsePath("next"))

	r, err := structural.DiffValue(o, n, nil)
	if err != nil {
		return v, err
	}

	return v.FillPath(cue.ParsePath("diff"), r), nil
}
