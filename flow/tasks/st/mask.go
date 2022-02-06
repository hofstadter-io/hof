package st

import (
	"cuelang.org/go/cue"

	"github.com/hofstadter-io/hof/flow/context"
	"github.com/hofstadter-io/cuetils/structural"
)

func init() {
  context.Register("st.Mask", NewMask)
}

type Mask struct {}

func NewMask(val cue.Value) (context.Runner, error) {
  return &Mask{}, nil
}

// Tasks must implement a Run func, this is where we execute our task
func (T *Mask) Run(ctx *context.Context) (interface{}, error) {
  ctx.CUELock.Lock()
  defer ctx.CUELock.Unlock()

	v := ctx.Value

	x := v.LookupPath(cue.ParsePath("val"))
	m := v.LookupPath(cue.ParsePath("mask"))

	r, err := structural.MaskValue(m, x, nil)
	if err != nil {
		return nil, err
	}

	return v.FillPath(cue.ParsePath("out"), r), nil
}
