package gen

import (
	"math/rand"

	"cuelang.org/go/cue"

	hofcontext "github.com/hofstadter-io/hof/flow/context"
)

type Float struct{}

func NewFloat(val cue.Value) (hofcontext.Runner, error) {
	return &Float{}, nil
}

func (T *Float) Run(ctx *hofcontext.Context) (interface{}, error) {
	f := rand.Float64()
	return f, nil
}

type Norm struct{}

func NewNorm(val cue.Value) (hofcontext.Runner, error) {
	return &Norm{}, nil
}

func (T *Norm) Run(ctx *hofcontext.Context) (interface{}, error) {
	f := rand.NormFloat64()
	return f, nil
}
