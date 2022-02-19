package gen

import (
	"cuelang.org/go/cue"
  "github.com/google/uuid"

  hofcontext "github.com/hofstadter-io/hof/flow/context"
)

type UUID struct {}

func NewUUID(val cue.Value) (hofcontext.Runner, error) {
  return &UUID{}, nil
}

func (T *UUID) Run(ctx *hofcontext.Context) (interface{}, error) {
  u := uuid.New()
	return u, nil
}
