package gen

import (
	"cuelang.org/go/cue"
  "github.com/google/uuid"

  "github.com/hofstadter-io/hof/flow/context"
)

func init() {
  context.Register("gen.UUID", NewUUID)
}

type UUID struct {}

func NewUUID(val cue.Value) (context.Runner, error) {
  return &UUID{}, nil
}

func (T *UUID) Run(ctx *context.Context) (interface{}, error) {
  u := uuid.New()
	return u, nil
}
