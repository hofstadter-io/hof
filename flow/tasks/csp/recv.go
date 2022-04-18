package csp

import (
	"fmt"

	"cuelang.org/go/cue"

	hofcontext "github.com/hofstadter-io/hof/flow/context"
	"github.com/hofstadter-io/hof/flow/flow"
	"github.com/hofstadter-io/hof/lib/cuetils"
)

type Recv struct {}

func NewRecv(val cue.Value) (hofcontext.Runner, error) {
  return &Recv{}, nil
}

func (T *Recv) Run(ctx *hofcontext.Context) (interface{}, error) {
  // fmt.Println("csp.Recv", ctx.Value)

	v := ctx.Value
  var (
    err error
    mailbox string
    quit string 
  )

  ferr := func () error {
    ctx.CUELock.Lock()
    defer func() {
      ctx.CUELock.Unlock()
    }()

    q := v.LookupPath(cue.ParsePath("quitMailbox"))
    if q.Exists() {
      if q.Err() != nil {
        return q.Err()
      }
      quit, err = q.String()
      if err != nil {
        return err
      }
    }

    nv := v.LookupPath(cue.ParsePath("mailbox")) 
    if !nv.Exists() {
      return fmt.Errorf("in csp.Recv task %s: missing field 'mailbox'", v.Path())
    }
    if nv.Err() != nil {
      return nv.Err()
    }
    mailbox, err = nv.String()
    if err != nil {
      return err 
    }

    return nil
  }()
  if ferr != nil {
    return nil, ferr
  }

  // load mailbox
  // fmt.Println("mailbox?:", mailbox)
  ci, loaded := ctx.Mailbox.Load(mailbox)
  if !loaded {
    return nil, fmt.Errorf("channel %q not found", mailbox)
  }

  c := ci.(chan Msg)

  var quitChan chan Msg
  if quit != "" {
    // fmt.Println("quitMailbox?:", quit)
    qi, loaded := ctx.Mailbox.Load(quit)
    if !loaded {
      return nil, fmt.Errorf("channel %q not found", quit)
    }
    quitChan = qi.(chan Msg)
  }

  handler := v.LookupPath(cue.ParsePath("handler"))
  if !handler.Exists() {
    // fmt.Println("got here")
    return nil, handler.Err()
  }

  // fmt.Println("handler:", handler)

  for {
    select {
    case <-quitChan:
      break 
  
    case msg := <-c:
      fmt.Println("msg:", msg)
      var H cue.Value

      ferr := func () error {
        ctx.CUELock.Lock()
        defer func() {
          ctx.CUELock.Unlock()
        }()

        H = ctx.Value.Context().CompileString("{...}")
        H = H.Unify(handler) 
        H = H.FillPath(cue.ParsePath("msg"), msg)

        return nil
      }()
      if ferr != nil {
        return nil, ferr
      }

      s, err := cuetils.PrintCue(H)
      if err != nil {
        fmt.Println("Error(csp/recv/print):", err)
        return nil, nil
      }
      fmt.Println("H:", s)

      p, err := flow.NewFlow(ctx, H)
      if err != nil {
        fmt.Println("Error(csp/recv/new):", err)
        return nil, nil
      }

      err = p.Start()
      if err != nil {
        fmt.Println("Error(csp/recv/run):", err)
        return nil, nil
      }
    }
  }
}
