package os

import (
	"os"

	"cuelang.org/go/cue"

	"github.com/hofstadter-io/hof/flow/context"
)

func init() {
  context.Register("os.Mkdir", NewMkdir)
}

type Mkdir struct {}

func NewMkdir(val cue.Value) (context.Runner, error) {
  return &Mkdir{}, nil
}

func (T *Mkdir) Run(ctx *context.Context) (interface{}, error) {

	v := ctx.Value

  var dir string
  var err error
  
  ferr := func () error {
    ctx.CUELock.Lock()
    defer func() {
      ctx.CUELock.Unlock()
    }()

    d := v.LookupPath(cue.ParsePath("dir")) 
    if d.Err() != nil {
      return d.Err()
    } else if d.Exists() {
      dir, err = d.String()
      if err != nil {
        return err
      }
    }
    return nil
  }()
  if ferr != nil {
    return nil, ferr
  }

  // TODO, make mode configurable
  err = os.MkdirAll(dir, 0755)

	return nil, err
}

