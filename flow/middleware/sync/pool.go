package sync

import (
	"fmt"

	"cuelang.org/go/cue"
	"github.com/gammazero/workerpool"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	flowctx "github.com/hofstadter-io/hof/flow/context"
	"github.com/hofstadter-io/hof/lib/hof"
)

type Pool struct {
	val  cue.Value
	next flowctx.Runner
}

func NewPool(opts flags.RootPflagpole, popts flags.FlowPflagpole) *Pool {
	// fmt.Println("Pool: new")
	return &Pool{}
}

func (M *Pool) Run(ctx *flowctx.Context) (results interface{}, err error) {
	val := ctx.Value

	node, err := hof.ParseHof[any](val)
	if err != nil {
		return nil, err
	}
	if node == nil  {
		panic("we should have found a node to even get here")
	}

	hofp := node.Hof.Flow.Pool

	// Make (pool) is setup before hand? (as middleware?)
	if hofp.Take {

		pool, ok := ctx.Pools.Load(hofp.Name)
		if !ok {
			return nil, fmt.Errorf("unknown exec pool %q @ %s\n", hofp.Name, val.Path())
		}

		P, ok := pool.(*workerpool.WorkerPool)
		P.SubmitWait(func() {
			results, err = M.next.Run(ctx)
		})
	} else {
		// fmt.Println("Pool: skip @", M.val.Path())
		results, err = M.next.Run(ctx)
	}

	return results, err
}

func (M *Pool) Apply(ctx *flowctx.Context, runner flowctx.RunnerFunc) flowctx.RunnerFunc {
	// fmt.Println("pool.Apply call")
	return func(val cue.Value) (flowctx.Runner, error) {

		// fmt.Println("pool.Apply func")

		// parse out the local #hof node data
		node, err := hof.ParseHof[any](val)
		if err != nil {
			return nil, err
		}
		if node == nil  {
			panic("we should have found a node to even get here")
		}

		// convenience
		hofp := node.Hof.Flow.Pool

		// hmmm, not sure if this actually runs it?
		// probably does, which is why it is paired with a no-op
		// so that a task is found, maybe this is where we are doubling up?
		next, err := runner(val)
		if err != nil {
			return nil, err
		}

		// if not a pool making task, return
		if hofp.Name == "" {
			return next, nil
		}

		// fmt.Printf("Pool: found @ %s %#v\n", val.Path(), hofp)

		// setup pool by name here

		if hofp.Make {
			pool := workerpool.New(hofp.Number)
			// fmt.Println("Pool: store @", val.Path(), pool)
			ctx.Pools.Store(hofp.Name, pool)
		}

		return &Pool{
			// required
			val:  val,
			next: next,
			// extra
		}, nil
	}
}
