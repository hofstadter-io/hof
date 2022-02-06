package gen

import (
  "time"

	"cuelang.org/go/cue"

  "github.com/hofstadter-io/hof/flow/context"
)

func init() {
  context.Register("gen.Now", NewNow)
}

type Now struct {}

func NewNow(val cue.Value) (context.Runner, error) {
  return &Now{}, nil
}

func (T *Now) Run(ctx *context.Context) (interface{}, error) {
  t := time.Now()
  f := t.Format(time.RFC3339)
	return f, nil
}
