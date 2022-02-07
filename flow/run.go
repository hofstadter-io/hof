package flow

import (
	go_ctx "context"
	"fmt"
	"os"
  "strings"
	"sync"

	// "time"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/flow/context"
	"github.com/hofstadter-io/hof/flow/flow"
	_ "github.com/hofstadter-io/hof/flow/tasks" // ensure tasks register
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

	val, err := structural.LoadCueInputs(entrypoints, ctx, nil)
	if err != nil {
    s := structural.FormatCueError(err)
		return fmt.Errorf("Error: %s", s)
	}

  if val.Err() != nil {
    s := structural.FormatCueError(val.Err())
    return fmt.Errorf("val.Err(): %s", s)  
  }

  // sharedCtx := buildSharedContext

	// (refactor/flow/many) find  flows
  flows := []*flow.Flow{}



  // (temp), give each own context (created in here), or maybe by flag? Need at least the shared mutex
  taskCtx, err := buildRootContext(val, opts, popts)
  // taskCtx, err := buildRootContext(sharedContex, val, opts, popts)
  if err != nil {
    return err
  }

  // this might be buggy?
  val, err = injectTags(val, popts.Tags)
  if err != nil {
    return err
  }

  // lets just print
  if popts.List {
    tags, secrets, errs := getTagsAndSecrets(val)
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
    err = listFlows(val, opts, popts)
    if err != nil {
      return err
    } 

  }

  flows, err = findFlows(taskCtx, val, opts, popts)
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
    err := flow.Start()
    if err != nil {
      return err
    }
  }

  //time.Sleep(time.Second)
  //fmt.Println("done")
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

func buildRootContext(val cue.Value, opts *flags.RootPflagpole, popts *flags.FlowFlagpole) (*context.Context, error) {
  // lookup the secret label in val
  // and build a filter write for stdout / stderr
  c := &context.Context{
    Stdin: os.Stdin,
    Stdout: os.Stdout,
    Stderr: os.Stderr,
    Context: go_ctx.Background(),
    CUELock: new(sync.Mutex),
    ValStore: new(sync.Map),
    Mailbox: new(sync.Map),
    DebugTasks: popts.DebugTasks,
    Verbosity: opts.Verbose,
  }
  return c, nil
}
