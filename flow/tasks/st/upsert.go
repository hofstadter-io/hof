package st

import (
	"fmt"
	"cuelang.org/go/cue"

	hofcontext "github.com/hofstadter-io/hof/flow/context"
	"github.com/hofstadter-io/hof/lib/cuetils"
	"github.com/hofstadter-io/hof/lib/structural"
)

type Upsert struct{}

func NewUpsert(val cue.Value) (hofcontext.Runner, error) {
	return &Upsert{}, nil
}

// Tasks must implement a Run func, this is where we execute our task
func (T *Upsert) Run(ctx *hofcontext.Context) (interface{}, error) {
	ctx.CUELock.Lock()
	defer ctx.CUELock.Unlock()
	v := ctx.Value

	x := v.LookupPath(cue.ParsePath("val"))
	u := v.LookupPath(cue.ParsePath("upsert"))

	r, err := structural.UpsertValue(u, x, nil)
	if err != nil {
		return nil, err
	}

	s, err := cuetils.PrintCue(r)
	fmt.Println(s)

	r, err = v.FillPath(cue.ParsePath("out"), r), nil
	return r, err
}
