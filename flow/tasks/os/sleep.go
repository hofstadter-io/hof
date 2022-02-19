package os

import (
	"time"

	"cuelang.org/go/cue"

  hofcontext "github.com/hofstadter-io/hof/flow/context"
)

type Sleep struct {}

func NewSleep(val cue.Value) (hofcontext.Runner, error) {
  return &Sleep{}, nil
}

func (T *Sleep) Run(ctx *hofcontext.Context) (interface{}, error) {

	v := ctx.Value

  var D time.Duration

  
  ferr := func () error {
    ctx.CUELock.Lock()
    defer func() {
      ctx.CUELock.Unlock()
    }()

    d := v.LookupPath(cue.ParsePath("duration")) 
    if d.Err() != nil {
      return d.Err()
    } else if d.Exists() {
      ds, err := d.String()
      if err != nil {
        return err
      }
      D, err = time.ParseDuration(ds)
      if err != nil {
        return  err
      }
    }
    return nil
  }()
  if ferr != nil {
    return nil, ferr
  }

  time.Sleep(D)

  var res interface{}
  func () {
    ctx.CUELock.Lock()
    defer func() {
      ctx.CUELock.Unlock()
    }()
    res = v.FillPath(cue.ParsePath("done"), true)
  }()

	return res, nil
}

