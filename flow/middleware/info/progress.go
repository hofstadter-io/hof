package info

import (
  "fmt"

  "cuelang.org/go/cue"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	hofcontext "github.com/hofstadter-io/hof/flow/context"
)

type Progress struct {
  val cue.Value
  next hofcontext.Runner
  use bool
}

func NewProgress(opts *flags.RootPflagpole, popts *flags.FlowFlagpole) (*Progress) {
  return &Progress{
    use: popts.DebugTasks,
  }
}

func (M *Progress) Run(ctx *hofcontext.Context) (results interface{}, err error) {
  fmt.Println("task: pre @", M.val.Path())
  result, err := M.next.Run(ctx)
  fmt.Println("task: post @", M.val.Path())
  return result, err
}

func (M *Progress) Apply(ctx *hofcontext.Context, runner hofcontext.RunnerFunc) hofcontext.RunnerFunc {
  if !M.use {
    return runner
  }
  return func(val cue.Value) (hofcontext.Runner, error) {
    fmt.Println("task: found @", val.Path(), val.Attributes(cue.ValueAttr))
    next, err := runner(val)
    if err != nil {
      return nil, err
    }
    return &Progress{
      val: val,
      next: next,
    }, nil
  }
}

