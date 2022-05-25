package info

import (
	"fmt"

	"cuelang.org/go/cue"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	hofcontext "github.com/hofstadter-io/hof/flow/context"
)

type Print struct {
	val  cue.Value
	next hofcontext.Runner
}

func NewPrint(opts *flags.RootPflagpole, popts *flags.FlowFlagpole) *Print {
	return &Print{}
}

func (M *Print) Run(ctx *hofcontext.Context) (results interface{}, err error) {
	result, err := M.next.Run(ctx)

	attrs := M.val.Attributes(cue.ValueAttr)
	for _, attr := range attrs {
		if attr.Name() == "print" {
			for i := 0; i < attr.NumArgs(); i++ {
				p, err := attr.String(i)
				if err != nil {
					return result, err
				}
				if p == "" {
					fmt.Printf("%s: %v\n", M.val.Path(), result)
				} else {
					r, ok := result.(cue.Value)
					var v cue.Value
					if ok {
						v = r.LookupPath(cue.ParsePath(p))
					} else {
						v = M.val.LookupPath(cue.ParsePath(p))
					}
					fmt.Printf("%s: %v\n", v.Path(), v)
				}
			}
			break
		}
	}

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
			val:  val,
			next: next,
		}, nil
	}
}
