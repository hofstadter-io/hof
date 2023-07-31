package cmd

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/fatih/color"

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
	watch.Watch(genRunner(false), wfiles, "runGen.wGen(full)", wait, quit)
	watch.Watch(genRunner(true),  xfiles, "runGen.xGen(fast)", wait, quit)

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

	R.PrintMergeConflicts()

	if R.GenFlags.Stats {
		R.PrintStats()
		fmt.Printf("\nTotal Elapsed Time: %s\n\n", elapsed)
	}

	if hasErr {
		return fmt.Errorf("ERROR: while running geneators")
	}

	return nil
}

func (R *Runtime) PrintMergeConflicts() {
	for _, G := range R.Generators {
		if G.Disabled {
			continue
		}

		R.printGenMergeConflicts(G)
	}
}

func (R *Runtime) printGenMergeConflicts(G *gen.Generator) {
	for _, F := range G.Files {
		if F.IsConflicted > 0 {
			msg := fmt.Sprintf("MERGE CONFLICT in %s", F.Filepath)
			color.Red(msg)
		}
	}

	for _, SG := range G.Generators {
		R.printGenMergeConflicts(SG)
	}

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

