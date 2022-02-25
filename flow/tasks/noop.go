package tasks

import (
	"cuelang.org/go/cue"

  hofcontext "github.com/hofstadter-io/hof/flow/context"
)

type Noop struct {}

func NewNoop(val cue.Value) (hofcontext.Runner, error) {
  return &Noop{}, nil
}

func (T *Noop) Run(ctx *hofcontext.Context) (interface{}, error) {
	return nil, nil
}
