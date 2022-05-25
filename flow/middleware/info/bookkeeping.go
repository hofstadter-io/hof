package info

import (
	"cuelang.org/go/cue"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	hofcontext "github.com/hofstadter-io/hof/flow/context"
)

type BookkeepingConfig struct {
	Workdir string
}

type Bookkeeping struct {
	cfg  BookkeepingConfig
	val  cue.Value
	next hofcontext.Runner
}

func NewBookkeeping(cfg BookkeepingConfig, opts *flags.RootPflagpole, popts *flags.FlowFlagpole) *Bookkeeping {
	return &Bookkeeping{
		cfg: cfg,
	}
}

func (M *Bookkeeping) Run(ctx *hofcontext.Context) (results interface{}, err error) {
	// bt := ctx.BaseTask
	// fmt.Println("bt:", bt.ID, bt.UUID)
	result, err := M.next.Run(ctx)

	// write out file in background
	return result, err
}

func (M *Bookkeeping) Apply(ctx *hofcontext.Context, runner hofcontext.RunnerFunc) hofcontext.RunnerFunc {
	return func(val cue.Value) (hofcontext.Runner, error) {
		// id := fmt.Sprint(val.Path())
		// fmt.Println("book: found @", val.Path())
		next, err := runner(val)
		if err != nil {
			return nil, err
		}
		return &Bookkeeping{
			val:  val,
			next: next,
		}, nil
	}
}

func (M *Bookkeeping) write(filename string, val cue.Value) error {

	return nil
}
