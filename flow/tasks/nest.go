package tasks

import (
  "fmt"

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
  fmt.Println("nest.Run:", val.Path())

  p, err := flow.NewFlow(ctx, val)
  if err != nil {
    fmt.Println("Error(nest/new):", err)
    return nil, nil
  }

  err = p.Start()
  if err != nil {
    fmt.Println("Error(nest/run):", err)
    return nil, err
  }

  return p.Final, nil
}
