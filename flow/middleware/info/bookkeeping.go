package info

import (
  // "fmt"

  "cuelang.org/go/cue"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	hofcontext "github.com/hofstadter-io/hof/flow/context"
)

type Bookkeeping struct {
  val cue.Value
  next hofcontext.Runner
}

func NewBookkeeping(opts *flags.RootPflagpole, popts *flags.FlowFlagpole) (*Bookkeeping) {
  return &Bookkeeping{}
}

func (M *Bookkeeping) Run(ctx *hofcontext.Context) (results interface{}, err error) {
  result, err := M.next.Run(ctx)
  return result, err
}

func (M *Bookkeeping) Apply(ctx *hofcontext.Context, runner hofcontext.RunnerFunc) hofcontext.RunnerFunc {
  return func(val cue.Value) (hofcontext.Runner, error) {
    // id := fmt.Sprint(val.Path())
    // fmt.Println("book: found@", val.Path())
    next, err := runner(val)
    if err != nil {
      return nil, err
    }
    return &Bookkeeping{
      val: val,
      next: next,
    }, nil
  }
}


