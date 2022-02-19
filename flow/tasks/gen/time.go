package gen

import (
  "time"

	"cuelang.org/go/cue"

  hofcontext "github.com/hofstadter-io/hof/flow/context"
)

type Now struct {}

func NewNow(val cue.Value) (hofcontext.Runner, error) {
  return &Now{}, nil
}

func (T *Now) Run(ctx *hofcontext.Context) (interface{}, error) {
  t := time.Now()
  f := t.Format(time.RFC3339)
	return f, nil
}
