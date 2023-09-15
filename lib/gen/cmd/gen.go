package cmd

import (
	"fmt"
	"os"
	"sync"
	"time"

	"cuelang.org/go/cue"
	"github.com/fatih/color"

	flowcontext "github.com/hofstadter-io/hof/flow/context"
	"github.com/hofstadter-io/hof/flow/middleware"
	"github.com/hofstadter-io/hof/flow/tasks"
	"github.com/hofstadter-io/hof/flow/flow"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/gen"
	"github.com/hofstadter-io/hof/lib/watch"
)

// this function needs to be broken down and split around ./runtime.go and lib/gen/*.go
func (R *Runtime) runGen() (err error) {
	// fix default Diff3 flag value when running hof gen
	// needs to be interwoven here, probably?
	// it's usage pattern is specific to our use cases right now, and want diff3 true for generators, but overridable if set to false
	hasDiff3Flag := false
	for _, arg := range os.Args {
		if arg == "--diff3" || arg == "-D" {
			hasDiff3Flag = true
			break
		}
	}

	// We need to set Diff3 default to true
	// when the user supplies generators and does not set flag
	if len(R.GenFlags.Template) == 0 {
		if !hasDiff3Flag {
			R.GenFlags.Diff3 = true
		}
	}

	if len(R.GenFlags.Template) == 0 {
		R.Diff3FlagSet = hasDiff3Flag
	}

	// everything above is about the --diff3 flag


	// First time code gen & output
	err = R.genOnce()
	if err != nil {
		return err
	}

	// return if we are not going into watch mode
	doWatch := shouldWatch(R.GenFlags)
	if !doWatch {
		return nil
	}

	// find our full & fast files
	wfiles, xfiles, err := R.buildWatchLists()
	if err != nil {
		return err
	}

	// reduce code below with this Runner builder
	genRunner := func(fast bool) watch.Runner {
		return func() error {
			var err error
			R.Lock()
			defer R.Unlock()

			err = R.Reload(fast)
			if err != nil {
				return err
			}
			return R.genOnce()
		}
	}

	// start up our watchers
	quit := make(chan bool, 3)
	wait := time.Millisecond * 50
	watch.Watch(genRunner(false), wfiles, "runGen.wGen(full)", wait, quit, true)
	watch.Watch(genRunner(true),  xfiles, "runGen.xGen(fast)", wait, quit, true)

	fmt.Println("watching for changes...")

	// main process waits here for ctrl-c
	// this is hacky, do this right
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()

	// send twice, once for each runner
	quit <- true
	quit <- true

	return nil
}


func (R *Runtime) genOnce() error {
	verystart := time.Now()

	// issue #20 - Don't print and exit on error here, wait until after we have written, so we can still write good files
	errsG := R.RunGenerators()
	errsW := R.WriteOutput()

	// final timing
	veryend := time.Now()
	elapsed := veryend.Sub(verystart).Round(time.Millisecond)

	// TODO (correctness)
	// ordering for the remainder of this function is unclear
	hasErr := false

	if len(errsG) > 0 {
		hasErr = true
		for _, e := range errsG {
			fmt.Println(e)
		}
	}
	if len(errsW) > 0 {
		hasErr = true
		for _, e := range errsW {
			fmt.Println(e)
		}
	}

	// TODO (shadow) not sure if we want to clean up gens without error?
	// right now, if any error, then no clean
	if !hasErr {
		errsS := R.CleanupRemainingShadow()
		if len(errsS) > 0 {
			hasErr = true
			for _, e := range errsS {
				fmt.Println(e)
			}
		}
	}

	nc := R.PrintMergeConflicts()
	if nc > 0 {
		fmt.Printf("\n%d merge conflict(s) found\n", nc)
	} else {
		// run post-exec here
		for _, G := range R.Generators {

			// maybe run post-flow per generator
			postFlow := G.CueValue.LookupPath(cue.ParsePath("PostFlow"))
			if postFlow.Exists() {
				if R.Flags.Verbosity > 0 {
					fmt.Println("running post-flow:", postFlow)
				}
				if !R.GenFlags.Exec {
					fmt.Println("skipping post-flow, use --exec to run")
				} else {
					ctx := flowcontext.New()
					ctx.RootValue = postFlow
					ctx.Stdin = os.Stdin
					ctx.Stdout = os.Stdout
					ctx.Stderr = os.Stderr
					ctx.Verbosity = R.Flags.Verbosity

					// how to inject tags into original value
					// fill / return value
					middleware.UseDefaults(ctx, R.Flags, flags.FlowPflags)
					tasks.RegisterDefaults(ctx)

					p, err := flow.OldFlow(ctx, postFlow)
					if err != nil {
						return err
					}

					err = p.Start()
					if err != nil {
						return err
					}

					// do we really want to fill anything in the value afterwards?
					G.CueValue = G.CueValue.FillPath(cue.ParsePath("PostFlow"), postFlow)
					if G.CueValue.Err() != nil {
						return err
					}
				}

			} else if !postFlow.Exists() {
				if G.Verbosity > 0 {
					fmt.Println("post-exec not found")
				}
			} else if postFlow.Err() != nil {
				return postFlow.Err()
			}
		}

	}


	if R.Flags.Stats {
		R.PrintStats()
		fmt.Printf("\nTotal Elapsed Time: %s\n\n", elapsed)
	}


	if hasErr {
		return fmt.Errorf("ERROR: while running generators")
	}

	return nil
}

func (R *Runtime) PrintMergeConflicts() (numConflict int) {
	for _, G := range R.Generators {
		if G.Disabled {
			continue
		}

		numConflict += R.printGenMergeConflicts(G)
	}
	return numConflict
}

func (R *Runtime) printGenMergeConflicts(G *gen.Generator) (numConflict int) {
	for _, F := range G.Files {
		if F.IsConflicted > 0 {
			numConflict += 1
			msg := fmt.Sprintf("MERGE CONFLICT in %s", F.Filepath)
			color.Red(msg)
		}
	}

	for _, SG := range G.Generators {
		numConflict += R.printGenMergeConflicts(SG)
	}

	return numConflict
}

func (R *Runtime) CleanupRemainingShadow() (errs []error) {
	if R.Flags.Verbosity > 0 {
		fmt.Println("Cleaning shadow")
	}

	for _, G := range R.Generators {
		gerrs := G.CleanupShadow(R.OutputDir(R.GenFlags.Outdir), R.ShadowDir(), R.Flags.Verbosity, R.GenFlags.KeepDeleted)
		errs = append(errs, gerrs...)
	}

	return errs
}

