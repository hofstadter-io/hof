package sync

import (
	"fmt"

	"cuelang.org/go/cue"
	"github.com/gammazero/workerpool"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	hofcontext "github.com/hofstadter-io/hof/flow/context"
)

type Pool struct {
	val  cue.Value
	next hofcontext.Runner
}

func NewPool(opts *flags.RootPflagpole, popts *flags.FlowFlagpole) *Pool {
	// fmt.Println("Pool: new")
	return &Pool{}
}

func (M *Pool) Run(ctx *hofcontext.Context) (results interface{}, err error) {
	val := ctx.Value
	attrs := val.Attributes(cue.ValueAttr)
	hasAttr := false
	var a cue.Attribute
	for _, attr := range attrs {
		n := attr.Name()
		if n == "pool" {
			a = attr
			hasAttr = true
			break
		}
	}

	if hasAttr && a.NumArgs() == 1 {
		pn, err := a.String(0)
		if err != nil {
			return nil, err
		}

		pool, ok := ctx.Pools.Load(pn)
		if !ok {
			return nil, fmt.Errorf("unknown exec pool %q @ %s\n", pn, val.Path())
		}

		P, ok := pool.(*workerpool.WorkerPool)
		P.SubmitWait(func() {
			// fmt.Println("Pool: run @", M.val.Path())
			results, err = M.next.Run(ctx)
		})
	} else {
		// fmt.Println("Pool: skip @", M.val.Path())
		results, err = M.next.Run(ctx)
	}

	// fmt.Println("Pool: post @", M.val.Path())

	return results, err
}

func (M *Pool) Apply(ctx *hofcontext.Context, runner hofcontext.RunnerFunc) hofcontext.RunnerFunc {
	return func(val cue.Value) (hofcontext.Runner, error) {
		hasAttr := false
		attrs := val.Attributes(cue.ValueAttr)
		var a cue.Attribute
		for _, attr := range attrs {
			n := attr.Name()
			if n == "pool" {
				a = attr
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

		// fmt.Printf("Pool: found @ %s %v %p\n", val.Path(), a, M)

		// setup pool by name here

		pn, err := a.String(0)
		if err != nil {
			return nil, err
		}

		if a.NumArgs() > 1 {
			// fmt.Println("Pool: make @", val.Path(), a)
			max, err := a.Int(1)
			if err != nil {
				return nil, err
			}

			pool := workerpool.New(int(max))
			// fmt.Println("Pool: store @", val.Path(), pool)
			ctx.Pools.Store(pn, pool)
		}

		return &Pool{
			// required
			val:  val,
			next: next,
			// extra
		}, nil
	}
}
