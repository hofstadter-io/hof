package gen

import (
  "math/rand"

	"cuelang.org/go/cue"

  "github.com/hofstadter-io/hof/flow/context"
)

func init() {
  context.Register("gen.Float", NewFloat)
  context.Register("gen.Norm", NewNorm)
}

type Float struct {}

func NewFloat(val cue.Value) (context.Runner, error) {
  return &Float{}, nil
}

func (T *Float) Run(ctx *context.Context) (interface{}, error) {
  f := rand.Float64()
	return f, nil
}

type Norm struct {}

func NewNorm(val cue.Value) (context.Runner, error) {
  return &Norm{}, nil
}

func (T *Norm) Run(ctx *context.Context) (interface{}, error) {
  f := rand.NormFloat64()
	return f, nil
}
