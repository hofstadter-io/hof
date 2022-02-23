package flow

import (
	"fmt"
	"os"
  "strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	hofcontext "github.com/hofstadter-io/hof/flow/context"
	hofflow "github.com/hofstadter-io/hof/flow/flow"
	"github.com/hofstadter-io/hof/flow/middleware"
	"github.com/hofstadter-io/hof/flow/task"
	"github.com/hofstadter-io/hof/flow/tasks" // ensure tasks register
	"github.com/hofstadter-io/hof/lib/structural"
)

/*
Input is to rigid
- reads from disk (can we workaround upstream?)
-
*/

func Run(entrypoints []string, opts *flags.RootPflagpole, popts *flags.FlowFlagpole) error {
	return run(entrypoints, opts, popts)
}

// refactor out single/multi
func run(entrypoints []string, opts *flags.RootPflagpole, popts *flags.FlowFlagpole) error {
	ctx := cuecontext.New()

	root, err := structural.LoadCueInputs(entrypoints, ctx, nil)
	if err != nil {
    s := structural.FormatCueError(err)
    return fmt.Errorf("root: Error: %s", s)
	}

  if root.Err() != nil {
    s := structural.FormatCueError(root.Err())
    return fmt.Errorf("root.Err(): %s", s)  
  }

  // sharedCtx := buildSharedContext

	// (refactor/flow/many) find  flows
  flows := []*hofflow.Flow{}


  // (temp), give each own context (created in here), or maybe by flag? Need at least the shared mutex
  // (todo) possibly get back new root because middleware injected flags/tags?
  taskCtx, err := buildRootContext(root, opts, popts)
  // taskCtx, err := buildRootContext(sharedContex, root, opts, popts)
  if err != nil {
    return err
  }

  // can middleware add flags to be filled before this? (yes ^^^)
  // how can we run middleware init type things (like creating a waitgroup or execpool, set logging level)

  // this might be buggy?
  root, err = injectTags(root, popts.Tags)
  if err != nil {
    return err
  }

  // lets just print
  if popts.List {
    tags, secrets, errs := getTagsAndSecrets(root)
    if len(errs) > 0 {
      return fmt.Errorf("in getTags: %v", errs)
    }
    if len(tags) > 0 {
      fmt.Println("tags:\n==============")
      for _, v := range tags {
        printHelpValue(v, "tag")
      }
      fmt.Println()
    }
    if len(secrets) > 0 {
      fmt.Println("secrets:\n==============")
      for _, v := range secrets {
        printHelpValue(v, "secret")
      }
      fmt.Println()
    }

    fmt.Println("flows:\n==============")
    err = listFlows(root, opts, popts)
    if err != nil {
      return err
    } 

  }

  flows, err = findFlows(taskCtx, root, opts, popts)
  if err != nil {
    s := structural.FormatCueError(err)
		return fmt.Errorf("Error: %s", s)
  }

  if popts.List {
    return nil
  }

  if len(flows) == 0 {
    return fmt.Errorf("no flows found")
  }

  // start all of the flows
  // TODO, use wait group, accume errors, flag for failure modes
  for _, flow := range flows {
    flow.Root = root
    err := flow.Start()
    if err != nil {
      return err
    }

    err = printFinalContext(flow.HofContext)
    if err != nil {
      return err
    }
  }

  // we are all done running flows now

  // print task dependencies / 

	return nil
}

func printHelpValue(v cue.Value, attr string) {
  path := v.Path()
  docs := v.Doc()
  if len(docs) > 0 {
    for _, d := range docs {
      fmt.Println("//", strings.TrimSpace(d.Text()))
    }
  }
  fmt.Printf("%s: %# v %v\n", path, v, v.Attribute(attr))
}

var walkOptions = []cue.Option{
  cue.Attributes(true),
  cue.Concrete(false),
  cue.Definitions(true),
  cue.Hidden(true),
  cue.Optional(true),
  cue.Docs(true),
}

func buildRootContext(val cue.Value, opts *flags.RootPflagpole, popts *flags.FlowFlagpole) (*hofcontext.Context, error) {
  // lookup the secret label in val
  // and build a filter write for stdout / stderr
  c := hofcontext.New()
  c.RootValue = val
  c.Stdin = os.Stdin
  c.Stdout = os.Stdout
  c.Stderr = os.Stderr
  c.DebugTasks = popts.DebugTasks
  c.Verbosity = opts.Verbose

  // how to inject tags into original value
  // fill / return value
  middleware.UseDefaults(c, opts, popts)
  tasks.RegisterDefaults(c)
  return c, nil
}

// print task dependencies / info
func printFinalContext(ctx *hofcontext.Context) error {
  // to start, print ids / timings
  // rebuild task dependencies with hof tasks from cue tasks

  fmt.Println("\n\n======= final =========")

  tm := map[string]*task.BaseTask{}

  ctx.Tasks.Range(func (key, value interface{}) bool {
    k := key.(string) 
    t := value.(*task.BaseTask)
    tm[k] = t
    return true
  })

  ti := make([]*task.BaseTask,len(tm))
  for _, t := range tm {
    ti[t.CueTask.Index()] = t
  }

  for _, t := range ti {
    b := t.TimeEvents["run.beg"]
    e := t.TimeEvents["run.end"]
    l := e.Sub(b)

    fmt.Println(t.ID, t.CueTask.Index(), t.UUID, l)

    for _, D := range t.CueTask.Dependencies() {
      fmt.Println(" -", D.Index())
    }
  }


  return nil
}
