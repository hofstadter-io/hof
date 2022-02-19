package dummy

import (
  "fmt"

  "cuelang.org/go/cue"

	hofcontext "github.com/hofstadter-io/hof/flow/context"
)

type Dummy struct {
  val cue.Value
  next hofcontext.RunnerFunc
}

func (M *Dummy) Run(ctx *hofcontext.Context) (results any, err error) {
  fmt.Println("dummy")
  r, err := M.next(M.val)
  if err != nil {
    return nil, err
  }
  return r.Run(ctx)
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

    if !hasAttr {
      return runner(val)
    }

    return &Dummy{
      val: val,
      next: runner,
    }, nil
  }
}

