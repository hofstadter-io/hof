package os

import (
  "fmt"
  "os"

	"cuelang.org/go/cue"

  hofcontext "github.com/hofstadter-io/hof/flow/context"
)

type WriteFile struct {}

func NewWriteFile(val cue.Value) (hofcontext.Runner, error) {
  return &WriteFile{}, nil
}

func (T *WriteFile) Run(ctx *hofcontext.Context) (interface{}, error) {

	v := ctx.Value
  var fn string // filename
  var bs []byte // contents
  var m int64   // mode
  var a bool    // append

  ferr := func () error {
    ctx.CUELock.Lock()
    defer func() {
      ctx.CUELock.Unlock()
    }()
    var err error

    f := v.LookupPath(cue.ParsePath("filename"))
    if !f.Exists() {
      return fmt.Errorf("missing required field 'filename'")
    }

    fn, err = f.String()
    if err != nil {
      return err
    }

    // switch on c's type to fill appropriately
    c := v.LookupPath(cue.ParsePath("contents"))
    if !c.Exists() {
      return fmt.Errorf("missing required field 'contents'")
    }

    switch k := c.IncompleteKind(); k {
    case cue.StringKind:
      s, err := c.Bytes()
      if err != nil {
        return err
      }
      bs = []byte(s)
      
    case cue.BytesKind:
      bs, err = c.Bytes()
      if err != nil {
        return err
      }

    default:
      return fmt.Errorf("Unsupported content type in WriteFile task: %q", k)
    }

    mode := v.LookupPath(cue.ParsePath("mode"))
    if !mode.Exists() {
      return fmt.Errorf("missing required field 'mode'")
    }
    m, err = mode.Int64()
    if err != nil {
      return err
    }

    av := v.LookupPath(cue.ParsePath("append"))
    if av.Exists() {
      a, err = av.Bool()
      if err != nil {
        return err
      }
    }

    return nil
  }()
  if ferr != nil {
    return nil, ferr
  }

  // todo, check failure modes, fill, not return error?
  // (in all tasks)

  om := os.O_CREATE|os.O_WRONLY
  if a {
    om = om | os.O_APPEND
  } else {
    om = om | os.O_TRUNC
  }

  f, err := os.OpenFile(fn, om, os.FileMode(m))
  if err != nil {
    return nil, err 
  }

  defer f.Close()
  if _, err := f.Write(bs); err != nil {
    return nil, err
  }
	return nil, err
}
