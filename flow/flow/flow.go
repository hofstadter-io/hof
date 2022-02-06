package flow

import (
	"fmt"

	"cuelang.org/go/cue"
	"cuelang.org/go/tools/flow"

  "github.com/hofstadter-io/hof/flow/context"
  "github.com/hofstadter-io/hof/flow/tasker"
	"github.com/hofstadter-io/cuetils/structural"
)

type Flow struct {
  Orig cue.Value
  Final cue.Value

  Context *context.Context
  Ctrl *flow.Controller
}

func NewFlow(ctx *context.Context, val cue.Value) (*Flow, error) {
  p := &Flow{
    Orig: val,
    Context: ctx,
  }
  return p, nil
}

// This is for the top-level flows
func (P *Flow) Start() error {
  return P.run(P.Orig)
}

func (P *Flow) run(val cue.Value) error {
	// Setup the flow Config
	cfg := &flow.Config{
		//InferTasks:     false,
		//IgnoreConcrete: true,
  }

  // copy orig for good measure
  // This is helpful for when 
  v := P.Orig.Context().CompileString("{...}")
  u := v.Unify(val) 

	// create the workflow which will build the task graph
	P.Ctrl = flow.New(cfg, u, tasker.NewTasker(P.Context))

  final, err := P.Ctrl.Run(P.Context.Context)

  // fmt.Println("flow(end):", P.path, P.rpath)
  P.Final = final
  if err != nil {
    s := structural.FormatCueError(err)
		return fmt.Errorf("Error: %s", s)
  }

  return nil
}

