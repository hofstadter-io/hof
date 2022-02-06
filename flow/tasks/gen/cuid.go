package gen

import (
	"cuelang.org/go/cue"
  "github.com/lucsky/cuid"

  "github.com/hofstadter-io/hof/flow/context"
)

func init() {
  context.Register("gen.CUID", NewCUID)
  context.Register("gen.Slug", NewSlug)
}

type CUID struct {}

func NewCUID(val cue.Value) (context.Runner, error) {
  return &CUID{}, nil
}

func (T *CUID) Run(ctx *context.Context) (interface{}, error) {
  u := cuid.New()
	return u, nil
}

type Slug struct {}

func NewSlug(val cue.Value) (context.Runner, error) {
  return &Slug{}, nil
}

func (T *Slug) Run(ctx *context.Context) (interface{}, error) {
  u := cuid.Slug()
	return u, nil
}
