package st

import (
	"cuelang.org/go/cue"

	hofcontext "github.com/hofstadter-io/hof/flow/context"
	"github.com/hofstadter-io/hof/lib/structural"
)

type Patch struct {}

func NewPatch(val cue.Value) (hofcontext.Runner, error) {
  return &Patch{}, nil
}

// Tasks must implement a Run func, this is where we execute our task
func (T *Patch) Run(ctx *hofcontext.Context) (interface{}, error) {
  ctx.CUELock.Lock()
  defer ctx.CUELock.Unlock()

	v := ctx.Value

	o := v.LookupPath(cue.ParsePath("orig"))
	n := v.LookupPath(cue.ParsePath("patch"))

	r, err := structural.PatchValue(o, n, nil)
	if err != nil {
		return nil, err
	}

	return v.FillPath(cue.ParsePath("next"), r), nil
}
