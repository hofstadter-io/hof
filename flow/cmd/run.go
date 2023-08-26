package cmd

import (
	"fmt"
	"os"
	"strings"

	"cuelang.org/go/cue"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/flow/flow"
	flowctx "github.com/hofstadter-io/hof/flow/context"
	"github.com/hofstadter-io/hof/flow/middleware"
	"github.com/hofstadter-io/hof/flow/task" // ensure tasks register
	"github.com/hofstadter-io/hof/flow/tasks" // ensure tasks register
)

func Run(args []string, rflags flags.RootPflagpole, cflags flags.FlowPflagpole) error {
	// unsugar the @flow-names into popts
	var entries, flowArgs, tagArgs []string
	for _, e := range args {
		if strings.HasPrefix(e, "@") {
			flowArgs = append(flowArgs, strings.TrimPrefix(e, "@"))
		} else if strings.HasPrefix(e, "+") {
			tagArgs = append(tagArgs, strings.TrimPrefix(e, "+"))
		} else {
			entries = append(entries, e)
		}
	}
	// update entrypoints and Flow flags
	args = entries
	rflags.Tags = append(rflags.Tags, tagArgs...)
	cflags.Flow = append(cflags.Flow, flowArgs...)


	// prep our runtime
	R, err := prepRuntime(args, rflags, cflags)
	if err != nil {
		return err
	}

	// this sets up the flows to run
	//err = R.EnrichFlows(cflags.Flow, EnrichFlows)
	//if err != nil {
	//  return err
	//}

	for _, flow := range R.Workflows {
		prepFlow(R, flow)

		if R.Flags.Verbosity > 0 {
			fmt.Println("running:", flow.Hof.Metadata.Name)
		}

		err := flow.Start()
		if err != nil {
			return err
		}

		if R.Flags.Stats {
			err = printFinalContext(flow.FlowCtx)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func prepFlow(R *Runtime, f *flow.Flow) {
	c := flowctx.New()
	c.RootValue = R.Value
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Verbosity = R.Flags.Verbosity

	// how to inject tags into original value
	// fill / return value
	middleware.UseDefaults(c, R.Flags, R.FlowFlags)
	tasks.RegisterDefaults(c)

	f.FlowCtx = c
}

func printFinalContext(ctx *flowctx.Context) error {
	// to start, print ids / timings
	// rebuild task dependencies with hof tasks from cue tasks

	fmt.Println("\n\n======= final =========")

	tm := map[string]*task.BaseTask{}

	ctx.Tasks.Range(func(key, value interface{}) bool {
		k := key.(string)
		t := value.(*task.BaseTask)
		tm[k] = t
		return true
	})

	ti := make([]*task.BaseTask, len(tm))
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

//func EnrichFlows(R *Runtime, cflags flags.FlowPflagpole) func (*runtime.Runtime, *flow.Flow) error {
//  return func(r *runtime.Runtime, f *flow.Flow) error {

//    f.FlowCtx = c

//    return nil
//  }
//}
