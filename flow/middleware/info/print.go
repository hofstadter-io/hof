package info

import (
  "fmt"

  "cuelang.org/go/cue"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	hofcontext "github.com/hofstadter-io/hof/flow/context"
)

type Print struct {
  val cue.Value
  next hofcontext.Runner
}

func NewPrint(opts *flags.RootPflagpole, popts *flags.FlowFlagpole) (*Print) {
  return &Print{}
}

func (M *Print) Run(ctx *hofcontext.Context) (results interface{}, err error) {
  result, err := M.next.Run(ctx)
  fmt.Printf("PRINT: %q: %v\n", M.val.Path(), M.val)
  return result, err
}

func (M *Print) Apply(ctx *hofcontext.Context, runner hofcontext.RunnerFunc) hofcontext.RunnerFunc {
  return func(val cue.Value) (hofcontext.Runner, error) {
    hasAttr := false
    attrs := val.Attributes(cue.ValueAttr)
    for _, attr := range attrs {
      if attr.Name() == "print" {
        hasAttr = true
        break
      }
    }

    next, err := runner(val)
    if err != nil {
      return nil, err
    }

    if !hasAttr {
      return next, nil
    }

    return &Print{
      val: val,
      next: next,
    }, nil
  }
}


