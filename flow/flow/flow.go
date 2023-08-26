package flow

import (
	"fmt"
	// "sync"

	"cuelang.org/go/cue"
	cueflow "cuelang.org/go/tools/flow"

	flowctx "github.com/hofstadter-io/hof/flow/context"
	"github.com/hofstadter-io/hof/flow/tasker"
	"github.com/hofstadter-io/hof/lib/cuetils"
	"github.com/hofstadter-io/hof/lib/hof"
)

type Flow struct {
	*hof.Node[Flow]

	Root  cue.Value
	Orig  cue.Value
	Final cue.Value

	FlowCtx *flowctx.Context
	Ctrl    *cueflow.Controller
}

func NewFlow(node *hof.Node[Flow]) *Flow {
	return &Flow{
		Node: node,
		Root: node.Value,
		Orig: node.Value,
	}
}

func OldFlow(ctx *flowctx.Context, val cue.Value) (*Flow, error) {
	p := &Flow{
		Root:    val,
		Orig:    val,
		FlowCtx: ctx,
	}
	return p, nil
}

// This is for the top-level flows
func (P *Flow) Start() error {
	return P.run()
}

func (P *Flow) run() error {
	// root := P.FlowCtx.RootValue
	root := P.Root
	// Setup the flow Config
	cfg := &cueflow.Config{
		InferTasks:     true,
		IgnoreConcrete: true,
		UpdateFunc: func(c *cueflow.Controller, t *cueflow.Task) error {
			return nil
		},
	}

	// This is for flows down from the root val
	// This is needed because nested flows (like IRC / API handler)
	// ... break if this check is not performed
	// ... and we blindly set the RootPath the value Path
	if P.Orig != P.Root {
		cfg.Root = P.Orig.Path()
	}

	// copy orig for good measure
	// This is helpful for when
	v := P.Orig.Context().CompileString("{...}")
	u := v.Unify(root)

	// create the workflow which will build the task graph
	P.Ctrl = cueflow.New(cfg, u, tasker.NewTasker(P.FlowCtx))

	err := P.Ctrl.Run(P.FlowCtx.GoContext)

	// fmt.Println("flow(end):", P.path, P.rpath)
	P.Final = P.Ctrl.Value()
	if err != nil {
		s := cuetils.CueErrorToString(err)
		fmt.Println(s)
		fmt.Println(P)
		return fmt.Errorf("Error in %s: %s", P.Hof.Metadata.Name, s)
	}

	return nil
}
