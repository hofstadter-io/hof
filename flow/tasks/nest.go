package tasks

import (
	"fmt"

	"cuelang.org/go/cue"

	hofcontext "github.com/hofstadter-io/hof/flow/context"
	"github.com/hofstadter-io/hof/flow/flow"
)

// this is buggy, need upstream support
type Nest struct{}

func NewNest(val cue.Value) (hofcontext.Runner, error) {
	return &Nest{}, nil
}

func (T *Nest) Run(ctx *hofcontext.Context) (interface{}, error) {
	val := ctx.Value

	orig := ctx.FlowStack
	ctx.FlowStack = append(orig, fmt.Sprint(val.Path()))

	fmt.Println(ctx.FlowStack)

	p, err := flow.OldFlow(ctx, val)
	if err != nil {
		return nil, err
	}

	fmt.Println(p.Hof)

	err = p.Start()
	if err != nil {
		return nil, err
	}

	ctx.FlowStack = orig

	return p.Final, nil
}
