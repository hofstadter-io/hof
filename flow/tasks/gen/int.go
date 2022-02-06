package gen

import (
  "math/rand"

	"cuelang.org/go/cue"

  "github.com/hofstadter-io/hof/flow/context"
)

func init() {
  context.Register("gen.Int", NewInt)
}

type Int struct {}

func NewInt(val cue.Value) (context.Runner, error) {
  return &Int{}, nil
}

func (T *Int) Run(ctx *context.Context) (interface{}, error) {

  val := ctx.Value

  var n int

  ferr := func () error {
    ctx.CUELock.Lock()
    defer func() {
      ctx.CUELock.Unlock()
    }()

    // lookup key
    nv := val.LookupPath(cue.ParsePath("n")) 
    if nv.Exists() {
      if nv.Err() != nil {
        return nv.Err() 
      }
      ni, err := nv.Int64()
      if err != nil {
        return err
      }
      n = int(ni)
    }

    return nil
  }()
  if ferr != nil {
    return nil, ferr
  }

  var i int
  if n == 0 {
    i = rand.Int()
  } else {
    i = rand.Intn(n)
  }

	return i, nil
}
