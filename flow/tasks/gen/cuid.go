package gen

import (
	"cuelang.org/go/cue"
  "github.com/lucsky/cuid"

  hofcontext "github.com/hofstadter-io/hof/flow/context"
)

type CUID struct {}

func NewCUID(val cue.Value) (hofcontext.Runner, error) {
  return &CUID{}, nil
}

func (T *CUID) Run(ctx *hofcontext.Context) (interface{}, error) {
  u := cuid.New()
	return u, nil
}

type Slug struct {}

func NewSlug(val cue.Value) (hofcontext.Runner, error) {
  return &Slug{}, nil
}

func (T *Slug) Run(ctx *hofcontext.Context) (interface{}, error) {
  u := cuid.Slug()
	return u, nil
}
