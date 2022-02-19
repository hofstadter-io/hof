package tasker

import (
	"fmt"

	"cuelang.org/go/cue"
	cueflow "cuelang.org/go/tools/flow"

  hofcontext "github.com/hofstadter-io/hof/flow/context"
  "github.com/hofstadter-io/hof/lib/cuetils"
)

func NewTasker(ctx *hofcontext.Context) cueflow.TaskFunc {
  // This function implements the Runner interface.
  // It parses Cue values, you will see all of them recursively
  return func(val cue.Value) (cueflow.Runner, error) {
    if len(val.Path().Selectors()) == 0 {
      return nil, nil
    }

    // Check that we have something that looks like a task
    // (look for attributes that match hof ones)
    attrs := val.Attributes(cue.ValueAttr)
    // skip if no attributes
    if len(attrs) == 0 {
      return nil, nil
    }

    // look for @task([string,...])
    for _, attr := range attrs {
      // TODO, iterate over all attrs and build them up?
      if attr.Name() == "task" {
        if attr.NumArgs() == 0 {
          return nil, fmt.Errorf("No type provided to task: %s", attr)
        }
        t, err := maybeTask(ctx, val, attr)
        if err != nil {
          fmt.Println("maybeTask err:", err)
        }
        return t, err 
      }
    }
    return nil, nil
  }
}

func maybeTask(ctx *hofcontext.Context, val cue.Value, attr cue.Attribute) (cueflow.Runner, error) {
  if ctx.DebugTasks {
    fmt.Println("task?:", attr)
  }

  taskId, err := attr.String(0)
  if err != nil {
    return nil, err
  }

  // lookup context.RunnerFunc 
  runnerFunc := ctx.Lookup(taskId)
  if runnerFunc == nil {
    return nil, fmt.Errorf("unknown task: %q at %q", attr, val.Path())
  }

  // some way to validate task against it's schema
  // (1) schemas self register
  // (2) here, we lookup schemas by taskId 
  // (3) use custom Require (or other validator)

  // create hof task from val
  // these live under /flow/tasks
  // and are of type context.RunnerFunc
  task, err := runnerFunc(val)
  if err != nil {
    return nil, err
  }

  // do per-task setup / common base / initial value / bookkeeping

  // wrap our RunnerFunc with cue/flow RunnerFunc
  return cueflow.RunnerFunc(func(t *cueflow.Task) error {
    c := hofcontext.Copy(ctx)

    c.Value = t.Value()

    // run the hof task 
    value, err := task.Run(c)
    if err != nil {
      err = fmt.Errorf("in %q, %v", val.Path(), err)
      c.Error = err
      return err
    }
    
    switch val := value.(type) {
    case cue.Value:
      attr := val.Attribute("print")
      err = cuetils.PrintAttr(attr, val)
    }

    if value != nil {
      err = t.Fill(value)
      if err != nil {
        c.Error = err
        return err
      }
    }
    return nil
  }), nil
}

