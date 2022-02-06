package os

import (
  "fmt"
  g_os "os"

	"cuelang.org/go/cue"

  "github.com/hofstadter-io/hof/flow/context"
)

func init() {
  context.Register("os.ReadFile", NewReadFile)
}

type ReadFile struct {}

func NewReadFile(val cue.Value) (context.Runner, error) {
  return &ReadFile{}, nil
}

func (T *ReadFile) Run(ctx *context.Context) (interface{}, error) {

	v := ctx.Value
  var fn string
  var err error

  ferr := func () error {
    ctx.CUELock.Lock()
    defer func() {
      ctx.CUELock.Unlock()
    }()
    f := v.LookupPath(cue.ParsePath("filename"))

    fn, err = f.String()
    return err
  }()
  if ferr != nil {
    return nil, ferr
  }

  bs, err := g_os.ReadFile(fn)
  if err != nil {
    return nil, err
  }

  var res cue.Value
  ferr = func () error {
    ctx.CUELock.Lock()
    defer func() {
      ctx.CUELock.Unlock()
    }()

    c := v.LookupPath(cue.ParsePath("contents"))

    // switch on c's type to fill appropriately
    switch k := c.IncompleteKind(); k {
    case cue.StringKind:
      res = v.FillPath(cue.ParsePath("contents"), string(bs))
    case cue.BytesKind:
      res = v.FillPath(cue.ParsePath("contents"), bs)

    case cue.StructKind:
      ctx := v.Context()
      c := ctx.CompileBytes(bs)
      if c.Err() != nil {
        return c.Err() 
      }
      res = v.FillPath(cue.ParsePath("contents"), c)

    case cue.BottomKind:
      res = v.FillPath(cue.ParsePath("contents"), string(bs))

    default:
      return fmt.Errorf("Unsupported Content type in ReadFile task: %q", k)
    }
    return nil
  }()
  if ferr != nil {
    return nil, ferr
  }

	return res, nil
}
