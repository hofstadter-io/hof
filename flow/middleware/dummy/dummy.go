package dummy

import (
  "fmt"

  "cuelang.org/go/cue"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	hofcontext "github.com/hofstadter-io/hof/flow/context"
)

type Dummy struct {
  val cue.Value
  next hofcontext.Runner
}

func NewDummy(opts *flags.RootPflagpole, popts *flags.FlowFlagpole) (*Dummy) {
  return &Dummy{}
}

func (M *Dummy) Run(ctx *hofcontext.Context) (results interface{}, err error) {
  fmt.Println("dummy: pre @", M.val.Path())
  // should this happen during discovery? (in Apply)
  result, err := M.next.Run(ctx)
  fmt.Println("dummy: post @", M.val.Path())

  return result, err
}

func (M *Dummy) Apply(ctx *hofcontext.Context, runner hofcontext.RunnerFunc) hofcontext.RunnerFunc {
  return func(val cue.Value) (hofcontext.Runner, error) {
    hasAttr := false
    attrs := val.Attributes(cue.ValueAttr)
    for _, attr := range attrs {
      if attr.Name() == "dummy" {
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

    fmt.Println("dummy: found @", val.Path())

    return &Dummy{
      val: val,
      next: next,
    }, nil
  }
}

