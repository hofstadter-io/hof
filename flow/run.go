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

  // unsugar the @flow-names into popts
  var entries, flowArgs, tagArgs []string
  for _, e := range entrypoints {
    if strings.HasPrefix(e, "@") {
      flowArgs = append(flowArgs,strings.TrimPrefix(e,"@"))
    } else if strings.HasPrefix(e, "+") {
      tagArgs = append(tagArgs,strings.TrimPrefix(e,"+"))
    } else {
      entries = append(entries,e)
    }
  }

  // update entrypoints and Flow flags
  entrypoints = entries
  popts.Flow = append(popts.Flow, flowArgs...)
  popts.Tags = append(popts.Tags, tagArgs...)

  // load in CUE files
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
    err = printFlows(root, opts, popts)
    if err != nil {
      return err
    } 

  }

  // get flow list and do some checks / updates
  names := flowList(root, opts, popts)
  // check for singular named flow in entire value
  // and run this if no flows set
  if len(popts.Flow) == 0 && len(names) == 1 && names[0] != "<unnamed>" {
    popts.Flow = []string{names[0]}
  }

	// (refactor/flow/many) find  flows
  flows := []*hofflow.Flow{}

  flows, err = findFlows(taskCtx, root, opts, popts)
  if err != nil {
    s := structural.FormatCueError(err)
		return fmt.Errorf("Error: %s", s)
  }

  if popts.List {
    return nil
  }

  if len(flows) == 0 {
    fmt.Println("available:\n==============")
    err = printFlows(root, opts, popts)
    fmt.Println()
    return fmt.Errorf("no flows matched")
  }

  // start all of the flows
  // TODO, use wait group, accume errors, flag for failure modes
  for _, flow := range flows {
    flow.Root = root
    err := flow.Start()
    if err != nil {
      return err
    }

    if popts.Stats {
      err = printFinalContext(flow.HofContext)
      if err != nil {
        return err
      }
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
    if t == nil {
      // panic("nil t")
      fmt.Println("nil t")
      continue
    }
    b := t.TimeEvents["run.beg"]
    e := t.TimeEvents["run.end"]
    l := e.Sub(b)

    // is := []int{}
    ps := []cue.Path{}
    for _, D := range t.CueTask.Dependencies() {
      // is = append(is, D.Index())
      ps = append(ps, D.Path())
    }
    if len(ps) > 0 {
      fmt.Println(t.ID, l, ps)
    } else {
      fmt.Println(t.ID, l)
    }
  }


  return nil
}
