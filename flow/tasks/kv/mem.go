package kv

import (
  "fmt"

	"cuelang.org/go/cue"

  hofcontext "github.com/hofstadter-io/hof/flow/context"
)

type Mem struct {}

func NewMem(val cue.Value) (hofcontext.Runner, error) {
  return &Mem{}, nil
}

func (T *Mem) Run(ctx *hofcontext.Context) (interface{}, error) {

  val := ctx.Value

  var k string
  var v interface{}
  var del bool
  var loaded bool
  var err error

  ferr := func () error {
    ctx.CUELock.Lock()
    defer func() {
      ctx.CUELock.Unlock()
    }()

    // lookup key
    key := val.LookupPath(cue.ParsePath("key")) 
    if key.Err() != nil {
      return key.Err() 
    } else if key.Exists() {
      k, err = key.String()
      if err != nil {
        return err
      }
    } else {
      err := fmt.Errorf("unknown key: %s", key)
      return err
    }

    // lookup val
    vv := val.LookupPath(cue.ParsePath("val")) 
    if vv.Exists() {
      v = vv
    }

    // lookup delete
    dv := val.LookupPath(cue.ParsePath("delete"))
    if dv.Exists() {
      del, err = dv.Bool()
      if err != nil {
        return err
      }
    }

    return nil
  }()
  if ferr != nil {
    return nil, ferr
  }

  if v != nil {
    if del {
      ctx.ValStore.Delete(k)
      v = nil
    } else {
      ctx.ValStore.Store(k,v)
    }
  } else {
    if del {
      v, loaded = ctx.ValStore.LoadAndDelete(k)
    } else {
      v, loaded = ctx.ValStore.Load(k)
    }

    // lock when we need to fill in a loaded value
    ctx.CUELock.Lock()
    defer ctx.CUELock.Unlock()

    val = val.FillPath(cue.ParsePath("val"), v)
    val = val.FillPath(cue.ParsePath("loaded"), loaded)

    return val, nil
  }


	return nil, nil
}

