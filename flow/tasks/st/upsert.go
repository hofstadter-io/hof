package st

import (
	"cuelang.org/go/cue"

	"github.com/hofstadter-io/hof/flow/context"
	"github.com/hofstadter-io/cuetils/structural"
)

func init() {
  context.Register("st.Upsert", NewUpsert)
}

type Upsert struct {}

func NewUpsert(val cue.Value) (context.Runner, error) {
  return &Upsert{}, nil
}

// Tasks must implement a Run func, this is where we execute our task
func (T *Upsert) Run(ctx *context.Context) (interface{}, error) {
  ctx.CUELock.Lock()
  defer ctx.CUELock.Unlock()
	v := ctx.Value

	x := v.LookupPath(cue.ParsePath("val"))
	u := v.LookupPath(cue.ParsePath("up"))

	r, err := structural.UpsertValue(u, x, nil)
	if err != nil {
		return nil, err
	}

	return v.FillPath(cue.ParsePath("out"), r), nil
}
