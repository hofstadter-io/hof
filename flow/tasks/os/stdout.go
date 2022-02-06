package os

import (
  "bufio"
  "fmt"

	"cuelang.org/go/cue"

  "github.com/hofstadter-io/hof/flow/context"
)

func init() {
  context.Register("os.Stdout", NewStdout)
}

type Stdout struct {}

func NewStdout(val cue.Value) (context.Runner, error) {
  return &Stdout{}, nil
}

func (T *Stdout) Run(ctx *context.Context) (interface{}, error) {
  bufStdout := bufio.NewWriter(ctx.Stdout)
  defer bufStdout.Flush()

  v := ctx.Value
  var m string
  var err error

  ferr := func () error {
    ctx.CUELock.Lock()
    defer func() {
      ctx.CUELock.Unlock()
    }()

    msg := v.LookupPath(cue.ParsePath("text")) 
    if msg.Err() != nil {
      return msg.Err() 
    } else if msg.Exists() {
      m, err = msg.String()
      if err != nil {
        return err
      }
    } else {
      err := fmt.Errorf("unknown msg:", msg)
      return err
    }
    return nil
  }()
  if ferr != nil {
    return nil, ferr
  }

  fmt.Fprint(bufStdout, m)
	return nil, nil
}
