package tasker

import (
	"fmt"

	"cuelang.org/go/cue"
	cueflow "cuelang.org/go/tools/flow"

  hofcontext "github.com/hofstadter-io/hof/flow/context"
  "github.com/hofstadter-io/hof/flow/task"
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
  //if ctx.DebugTasks {
    //fmt.Println("task?:", attr)
  //}

  taskId, err := attr.String(0)
  if err != nil {
    return nil, err
  }

  // lookup context.RunnerFunc 
  runnerFunc := ctx.Lookup(taskId)
  if runnerFunc == nil {
    return nil, fmt.Errorf("unknown task: %q at %q", attr, val.Path())
  }

  // Note, we apply this in the reverse order so that the Use order is like a stack
  // (i.e. the first is the most outer, which is typical for how these work for servers
  // apply plugin / middleware
  for i := len(ctx.Middlewares)-1; i>=0; i-- {
    ware := ctx.Middlewares[i]
    runnerFunc = ware.Apply(ctx, runnerFunc)
  }



  // some way to validate task against it's schema
  // (1) schemas self register
  // (2) here, we lookup schemas by taskId 
  // (3) use custom Require (or other validator)

  // create hof task from val
  // these live under /flow/tasks
  // and are of type context.RunnerFunc
  T, err := runnerFunc(val)
  if err != nil {
    return nil, err
  }

  // do per-task setup / common base / initial value / bookkeeping
  bt := task.NewBaseTask(val)
  ctx.Tasks.Store(bt.ID, bt)

  // wrap our RunnerFunc with cue/flow RunnerFunc
  return cueflow.RunnerFunc(func(t *cueflow.Task) error {
    // why do we need a copy?
    // maybe for local Value / CurrTask
    c := hofcontext.Copy(ctx)

    c.Value = t.Value()
    c.BaseTask = bt

    bt.CueTask = t
    bt.Start = c.Value
    bt.Final = c.Value

    // run the hof task 
    bt.AddTimeEvent("run.beg")
    // (update)
    value, err := T.Run(c)
    bt.AddTimeEvent("run.end")

    if err != nil {
      err = fmt.Errorf("in %q, %v", val.Path(), err)
      c.Error = err
      bt.Error = err
      return err
    }

    if value != nil {
      bt.AddTimeEvent("fill.beg")
      err = t.Fill(value)
      bt.AddTimeEvent("fill.end")
      if err != nil {
        c.Error = err
        bt.Error = err
        return err
      }

      bt.Final = t.Value() 

    }
    return nil
  }), nil
}

