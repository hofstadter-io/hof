package info

import (
	"fmt"
	"strings"

	"cuelang.org/go/cue"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	hofcontext "github.com/hofstadter-io/hof/flow/context"
)

type Progress struct {
	val  cue.Value
	next hofcontext.Runner
	use  bool
}

func NewProgress(opts *flags.RootPflagpole, popts *flags.FlowFlagpole) *Progress {
	return &Progress{
		use: popts.Progress,
	}
}

func (M *Progress) Run(ctx *hofcontext.Context) (results interface{}, err error) {
	bt := ctx.BaseTask
	fmt.Println("bt:", bt.ID, bt.UUID)
	fmt.Println("task: pre @", strings.Join(ctx.FlowStack, "."), M.val.Path())
	result, err := M.next.Run(ctx)
	fmt.Println("task: post @", strings.Join(ctx.FlowStack, "."), M.val.Path())
	return result, err
}

func (M *Progress) Apply(ctx *hofcontext.Context, runner hofcontext.RunnerFunc) hofcontext.RunnerFunc {
	if !M.use {
		return runner
	}
	return func(val cue.Value) (hofcontext.Runner, error) {
		fmt.Println("task: found @", strings.Join(ctx.FlowStack, "."), val.Path(), val.Attributes(cue.ValueAttr))
		next, err := runner(val)
		if err != nil {
			return nil, err
		}
		return &Progress{
			val:  val,
			next: next,
		}, nil
	}
}
