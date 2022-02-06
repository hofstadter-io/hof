package st

import (
	"cuelang.org/go/cue"

	"github.com/hofstadter-io/hof/flow/context"
	"github.com/hofstadter-io/cuetils/structural"
)

func init() {
  context.Register("st.Patch", NewPatch)
}

type Patch struct {}

func NewPatch(val cue.Value) (context.Runner, error) {
  return &Patch{}, nil
}

// Tasks must implement a Run func, this is where we execute our task
func (T *Patch) Run(ctx *context.Context) (interface{}, error) {
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
