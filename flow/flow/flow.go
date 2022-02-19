package flow

import (
	"fmt"
  // "sync"

	"cuelang.org/go/cue"
	cueflow "cuelang.org/go/tools/flow"

  hofcontext "github.com/hofstadter-io/hof/flow/context"
  "github.com/hofstadter-io/hof/flow/tasker"
	"github.com/hofstadter-io/hof/lib/structural"
)

type Flow struct {
  Root cue.Value
  Orig cue.Value
  Final cue.Value

  HofContext *hofcontext.Context
  Ctrl *cueflow.Controller
}

func NewFlow(ctx *hofcontext.Context, val cue.Value) (*Flow, error) {
  p := &Flow{
    Orig: val,
    HofContext: ctx,
  }
  return p, nil
}

// This is for the top-level flows
func (P *Flow) Start() error {
  return P.run()
}

func (P *Flow) run() error {
  root := P.Root
  // root := P.Orig
	// Setup the flow Config
	cfg := &cueflow.Config{
    Root: P.Orig.Path(),
		InferTasks:     true,
		//IgnoreConcrete: true,
    UpdateFunc: func(c *cueflow.Controller, t *cueflow.Task) error {

      //if t != nil {
        //fmt.Printf("task(%d): %s %s => %v\n", t.Index(), t.Path(), t.State(), t.Value())
      //}

      return nil
    },
  }

  // copy orig for good measure
  // This is helpful for when 
  v := P.Orig.Context().CompileString("{...}")
  u := v.Unify(root) 

	// create the workflow which will build the task graph
	P.Ctrl = cueflow.New(cfg, u, tasker.NewTasker(P.HofContext))

  final, err := P.Ctrl.Run(P.HofContext.GoContext)

  // fmt.Println("flow(end):", P.path, P.rpath)
  P.Final = final
  if err != nil {
    s := structural.FormatCueError(err)
		return fmt.Errorf("Error: %s", s)
  }

  return nil
}

