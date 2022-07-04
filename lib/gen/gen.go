package gen

import (
	"fmt"
	"time"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/cuetils"
)

func Gen(args []string, rootflags flags.RootPflagpole, cmdflags flags.GenFlagpole) error {

	if len(cmdflags.Generator) == 0 && len(cmdflags.Template) > 0 {
		return Render(args, rootflags, cmdflags)
	}

	verystart := time.Now()

	var errs []error

	R := NewRuntime(args, cmdflags)

	errs = R.LoadCue()
	if len(errs) > 0 {
		for _, e := range errs {
			cuetils.PrintCueError(e)
		}
		return fmt.Errorf("\nErrors while loading cue files\n")
	}

	errsL := R.LoadGenerators()
	if len(errsL) > 0 {
		for _, e := range errsL {
			fmt.Println(e)
			// cuetils.PrintCueError(e)
		}
		return fmt.Errorf("\nErrors while loading generators\n")
	}

	// issue #20 - Don't print and exit on error here, wait until after we have written, so we can still write good files
	errsG := R.RunGenerators()
	errsW := R.WriteOutput()

	// final timing
	veryend := time.Now()
	elapsed := veryend.Sub(verystart).Round(time.Millisecond)

	if cmdflags.Stats {
		R.PrintStats()
		fmt.Printf("\nTotal Elapsed Time: %s\n\n", elapsed)
	}

	if len(errsG) > 0 {
		for _, e := range errsG {
			fmt.Println(e)
		}
		return fmt.Errorf("\nErrors while generating output\n")
	}
	if len(errsW) > 0 {
		for _, e := range errsW {
			fmt.Println(e)
		}
		return fmt.Errorf("\nErrors while writing output\n")
	}

	R.PrintMergeConflicts()

	return nil
}
