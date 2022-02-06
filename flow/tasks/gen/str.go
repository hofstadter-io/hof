package gen

import (
  "math/rand"

	"cuelang.org/go/cue"

  "github.com/hofstadter-io/hof/flow/context"
)

// default runes if none provided
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init() {
  context.Register("gen.Str", NewStr)
}

type Str struct {}

func NewStr(val cue.Value) (context.Runner, error) {
  return &Str{}, nil
}

func (T *Str) Run(ctx *context.Context) (interface{}, error) {

  val := ctx.Value

  n := 12
  runes := letters

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

    // lookup key
    rv := val.LookupPath(cue.ParsePath("runes")) 
    if rv.Exists() {
      if rv.Err() != nil {
        return rv.Err() 
      }
      s, err := rv.String()
      if err != nil {
        return err
      }
      runes = []rune(s)
    }

    return nil
  }()
  if ferr != nil {
    return nil, ferr
  }

  b := make([]rune, n)
  for i := range b {
      b[i] = runes[rand.Intn(len(runes))]
  }

  r := string(b)

	return r, nil
}
