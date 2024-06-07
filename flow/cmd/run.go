package cmd

import (
	"fmt"
	"os"
	"strings"

	"cuelang.org/go/cue"
	"github.com/gammazero/workerpool"
	"github.com/olekukonko/tablewriter"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	flowctx "github.com/hofstadter-io/hof/flow/context"
	"github.com/hofstadter-io/hof/flow/flow"
	"github.com/hofstadter-io/hof/flow/middleware"
	"github.com/hofstadter-io/hof/flow/task"  // ensure tasks register
	"github.com/hofstadter-io/hof/flow/tasks" // ensure tasks register
	"github.com/hofstadter-io/hof/lib/hof"
	"github.com/hofstadter-io/hof/lib/yagu"
)

func prepFlow(R *Runtime, val cue.Value) (*flow.Flow, error) {
	node, err := hof.ParseHof[flow.Flow](val)
	if err != nil {
		return nil, err
	}

	c := flowctx.New()
	c.RootValue = val
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Verbosity = R.Flags.Verbosity
	c.ShowErrors = R.Flags.AllErrors

	// how to inject tags into original value
	// fill / return value
	middleware.UseDefaults(c, R.Flags, R.FlowFlags)
	tasks.RegisterDefaults(c)

	f, err := flow.OldFlow(c, val)
	f.Node = node
	return f, err
}

func Run(args []string, rflags flags.RootPflagpole, cflags flags.FlowPflagpole) error {

	wp := workerpool.New(cflags.Parallel)

	// prep our runtime
	R, err := prepRuntime(args, rflags, cflags)
	if err != nil {
		return err
	}

	var src, dst string
	if cflags.Bulk != "" {
		parts := strings.Split(cflags.Bulk, "@")
		if len(parts) != 2 {
			return fmt.Errorf("bad format for -B/--bulk flag, requires <src.path>@<dst.path>")
		}
		src, dst = parts[0], parts[1]
		if src == "" || dst == "" {
			return fmt.Errorf("bad format for -B/--bulk flag, requires <src.path>@<dst.path>")
		}
	}

	errCnt := 0

	for _, WF := range R.Workflows {

		if R.Flags.Verbosity > 0 {
			fmt.Println("running:", WF.Hof.Metadata.Name)
		}

		// runs the workflow in a single value
		fn := func(val cue.Value) error {

			F, err := prepFlow(R, val)
			if err != nil {
				return err
			}

			err = F.Start()
			if err != nil {
				return err
			}

			if R.Flags.Stats {
				err = printFinalContext(F.FlowCtx)
				if err != nil {
					return err
				}
			}

			return nil
		}

		// bulk processing
		if src != "" && dst != "" {
			fmt.Printf("flowing %q in bulk mode using %d workers\n", WF.Hof.Flow.Name, cflags.Parallel)
			// get Src data
			Src := R.Value.LookupPath(cue.ParsePath(src))

			// build up iter from Src
			var iter *cue.Iterator
			switch Src.IncompleteKind() {
			case cue.StructKind:
				iter, err = Src.Fields()
			case cue.ListKind:
				var i cue.Iterator
				i, err = Src.List()
				iter = &i
			default:
				fmt.Println("unknown iterable", Src.Validate())	
			}
			if err != nil {
				return err
			}

			// loop over data
			for iter.Next() {
				data := iter.Value()

				wp.Submit(func(){
					fmt.Println(">>>", data.Path())
		
					v := WF.Root.FillPath(cue.ParsePath(dst), data)

					err := fn(v)
					if err != nil {
						fmt.Println(err)
						errCnt += 1
					}
					fmt.Println()
				})
			}	

			wp.StopWait()

		} else {
			wp.Submit(func(){
				err := fn(WF.Root)
				if err != nil {
					fmt.Println(err)
					errCnt += 1
				}
			})
		} 
	}

	wp.StopWait()
	if errCnt > 0 {
		return fmt.Errorf("%d error(s) were encountered", errCnt)
	}

	return nil
}

func printFinalContext(ctx *flowctx.Context) error {
	// to start, print ids / timings
	// rebuild task dependencies with hof tasks from cue tasks
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


	fmt.Println("\n\n======= final =========")

	return yagu.PrintAsTable(
		[]string{"Name", "Task", "Deps", "Time", "Err"},
		func(table *tablewriter.Table) ([][]string, error) {
			var rows = make([][]string, 0, len(ti))
			// fill with data

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

				deps := ""
				if len(ps) > 0 {
					deps = fmt.Sprint(ps)
				}

				tt := ""
				attrs := t.Orig.Attributes(cue.ValueAttr)
				for _, a := range attrs {
					if a.Name() == "task" {
						tt = fmt.Sprint(a)
						break
					}
				}

				// fmt.Println(t.Final)
				final := t.CueTask.Value()
				// fmt.Println(final)

				err := ""
				if t.Error != nil {
					err = fmt.Sprint(t.Error)
				} else if e := final.LookupPath(cue.ParsePath("error")); e.Exists() {
					err = fmt.Sprint(e)
				}

				row := []string{t.ID, tt, deps, fmt.Sprint(l), err}
				rows = append(rows, row)
			}

			return rows, nil
		},
	)


	return nil
}

//func EnrichFlows(R *Runtime, cflags flags.FlowPflagpole) func (*runtime.Runtime, *flow.Flow) error {
//  return func(r *runtime.Runtime, f *flow.Flow) error {

//    f.FlowCtx = c

//    return nil
//  }
//}
