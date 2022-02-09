package csp

import (
  "fmt"

	"cuelang.org/go/cue"

  "github.com/hofstadter-io/hof/flow/context"
)

func init() {
  context.Register("csp.Chan", NewChan)
}

type Chan struct {}

func NewChan(val cue.Value) (context.Runner, error) {
  return &Chan{}, nil
}

func (T *Chan) Run(ctx *context.Context) (interface{}, error) {
	v := ctx.Value

  fmt.Println("csp.Chan", v)

  var err error
  var mailbox string
  var buf int

  ferr := func () error {
    ctx.CUELock.Lock()
    defer func() {
      ctx.CUELock.Unlock()
    }()

    nv := v.LookupPath(cue.ParsePath("mailbox")) 
    if !nv.Exists() {
      return fmt.Errorf("in csp.Chan task %s: missing field 'mailbox'", v.Path())
    }
    if nv.Err() != nil {
      return nv.Err()
    }
    mailbox, err = nv.String()
    if err != nil {
      return err 
    }

    iv := v.LookupPath(cue.ParsePath("buf")) 
    if iv.Exists() {
      if iv.Err() != nil {
        return iv.Err()
      }
      i64, err := iv.Int64()
      if err != nil {
        return err 
      }
      buf = int(i64)
    }

    return nil
  }()
  if ferr != nil {
    return nil, ferr
  }

  // make mailbox in it doesn't exist
  // todo, lookup prior art in CSP
  _, loaded := ctx.Mailbox.Load(mailbox)
  if !loaded {
    fmt.Println("new mailbox!  ", mailbox)
    c := make(chan Msg,buf)
    ctx.Mailbox.Store(mailbox, c)
  }
  fmt.Println("mailbox saved")

	return nil, nil 
}


