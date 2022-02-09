package ext

import (
  "fmt"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/format"

  "github.com/hofstadter-io/hof/flow/context"
)

func init() {
  context.Register("ext.CueFormat", NewCueFormat)
}

type CueFormat struct {}

func NewCueFormat(val cue.Value) (context.Runner, error) {
  return &CueFormat{}, nil
}

func (T *CueFormat) Run(ctx *context.Context) (interface{}, error) {

	v := ctx.Value
  var val cue.Value

  ferr := func () error {
    ctx.CUELock.Lock()
    defer func() {
      ctx.CUELock.Unlock()
    }()

    val := v.LookupPath(cue.ParsePath("val")) 
    if !val.Exists() {
      return fmt.Errorf("in task %s: missing field 'value'", v.Path())
    }
    if val.Err() != nil {
      return val.Err()
    }

    return nil
  }()
  if ferr != nil {
    return nil, ferr
  }

	syn := val.Syntax(
		cue.Final(),
		cue.ResolveReferences(true),
		cue.Definitions(true),
		cue.Hidden(true),
		cue.Optional(true),
		cue.Attributes(true),
		cue.Docs(true),
	)

	bs, err := format.Node(syn)
	if err != nil {
    return nil, err
	}

  ctx.CUELock.Lock()
  defer ctx.CUELock.Unlock()
  res := v.FillPath(cue.ParsePath("out"), string(bs))

	return res, nil 
}

