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
	err := P.run()
	// fmt.Println("Start().Err", P.Orig.Path(), err)	
	return err
}

func (P *Flow) run() error {
	// fmt.Println("FLOW.run:", P.FlowCtx.RootValue.Path(), P.Root.Path())
	// root := P.FlowCtx.RootValue
	root := P.Root
	// Setup the flow Config
	cfg := &cueflow.Config{
		// InferTasks:      true,
		IgnoreConcrete:  true,
		FindHiddenTasks: true,
		UpdateFunc: func(c *cueflow.Controller, t *cueflow.Task) error {
			//if t != nil {
			//  fmt.Println("Flow.Update()", t.Index(), t.Path())
			//} else {
			//  fmt.Println("Flow.Update()", "nil task")
			//}
			if t != nil {
				v := t.Value()

				node, err := hof.ParseHof[any](v)
				if err != nil {
					return err
				}
				if node == nil  {
					panic("we should have found a node to even get here")
				}

				if node.Hof.Flow.Print.Level > 0 && !node.Hof.Flow.Print.Before {
					pv := v.LookupPath(cue.ParsePath(node.Hof.Flow.Print.Path))
					if node.Hof.Path == "" {
						fmt.Printf("%s", node.Hof.Flow.Print.Path)
					} else if node.Hof.Flow.Print.Path == "" {
						fmt.Printf("%s", node.Hof.Path)
					} else {
						fmt.Printf("%s.%s", node.Hof.Path, node.Hof.Flow.Print.Path)
					}
					fmt.Printf(": %v\n", pv)
				}
			}
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

	// fmt.Println("Flow.run() start")
	err := P.Ctrl.Run(P.FlowCtx.GoContext)
	// fmt.Println("Flow.run() end", err)

	// fmt.Println("flow(end):", P.path, P.rpath)
	P.Final = P.Ctrl.Value()
	if err != nil {
		s := cuetils.CueErrorToString(err)
		// fmt.Println("Flow ERR in?", P.Orig.Path(), s)
		
		//fmt.Println(P)
		return fmt.Errorf("Error in %s | %s: %s", P.Hof.Metadata.Name, P.Orig.Path(), s)
	}
	// fmt.Println("NOT HERE", P.Orig.Path())

	return nil
}
