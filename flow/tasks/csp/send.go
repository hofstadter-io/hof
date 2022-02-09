package csp

import (
  "fmt"

	"cuelang.org/go/cue"

  "github.com/hofstadter-io/hof/flow/context"
)

func init() {
  context.Register("csp.Send", NewSend)
}

type Send struct {}

func NewSend(val cue.Value) (context.Runner, error) {
  return &Send{}, nil
}

func (T *Send) Run(ctx *context.Context) (interface{}, error) {
  fmt.Println("csp.Send", ctx.Value)

	v := ctx.Value
  var (
    err error
    name string
    key string
    val cue.Value
  )

  ferr := func () error {
    ctx.CUELock.Lock()
    defer func() {
      ctx.CUELock.Unlock()
    }()

    val := v.LookupPath(cue.ParsePath("val")) 
    if !val.Exists() {
      return fmt.Errorf("in csp.Send task %s: missing field 'val'", v.Path())
    }
    if val.Err() != nil {
      return val.Err()
    }

    kv := v.LookupPath(cue.ParsePath("name")) 
    if kv.Exists() {
      if kv.Err() != nil {
        return kv.Err()
      }
      name, err = kv.String()
      if err != nil {
        return err 
      }
    }

    nv := v.LookupPath(cue.ParsePath("name")) 
    if !nv.Exists() {
      return fmt.Errorf("in csp.Send task %s: missing field 'name'", v.Path())
    }
    if nv.Err() != nil {
      return nv.Err()
    }
    name, err = nv.String()
    if err != nil {
      return err 
    }

    return nil
  }()
  if ferr != nil {
    return nil, ferr
  }

  // load mailbox
  fmt.Println("mailbox?:", name)
  ci, loaded := ctx.Mailbox.Load(name)
  if !loaded {
    return nil, fmt.Errorf("channel %q not found", name)
  }

  msg := Msg {
    Key: key,
    Val: val,
  }
  fmt.Println("sending:", msg)
  // send a Msg
  c := ci.(chan Msg)
  c <- msg 

	return nil, nil 
}
