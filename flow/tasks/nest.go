package tasks

import (
	"fmt"

	"cuelang.org/go/cue"

	flowctx "github.com/hofstadter-io/hof/flow/context"
	"github.com/hofstadter-io/hof/flow/flow"
	"github.com/hofstadter-io/hof/lib/hof"
)

// this is buggy, need upstream support
type Nest struct{}

func NewNest(val cue.Value) (flowctx.Runner, error) {
	return &Nest{}, nil
}

func (T *Nest) Run(ctx *flowctx.Context) (interface{}, error) {
	val := ctx.Value

	orig := ctx.FlowStack
	ctx.FlowStack = append(orig, fmt.Sprint(val.Path()))

	n, err := hof.ParseHof[flow.Flow](val)
	if err != nil {
		return nil, err
	}

	p, err := flow.OldFlow(ctx, val)
	if err != nil {
		return nil, err
	}

	p.Node = n

	err = p.Start()
	if err != nil {
		return nil, fmt.Errorf("in nested task: %w", err)
	}

	ctx.FlowStack = orig

	return p.Final, nil
}
