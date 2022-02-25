package tasks

import (
	"cuelang.org/go/cue"

  hofcontext "github.com/hofstadter-io/hof/flow/context"
	"github.com/hofstadter-io/hof/flow/flow"
)

type Nest struct {}

func NewNest(val cue.Value) (hofcontext.Runner, error) {
  return &Nest{}, nil
}

func (T *Nest) Run(ctx *hofcontext.Context) (interface{}, error) {
	val := ctx.Value
  // fmt.Println("nest.Run:", val.Path())

  p, err := flow.NewFlow(ctx, val)
  if err != nil {
    return nil, err
  }

  err = p.Start()
  if err != nil {
    return nil, err
  }

  return p.Final, nil
}
