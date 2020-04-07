package lib

import (
	"fmt"
	"time"

	"github.com/hofstadter-io/hof/lib/util"
)

func Gen(entrypoints, expressions []string, mode string) (error) {
	verystart := time.Now()

	var errs []error

	R := NewRuntime(entrypoints, expressions)

	errs = R.LoadCue()
	if len(errs) > 0 {
		for _, e := range errs {
			util.PrintCueError(e)
		}
		return fmt.Errorf("\nErrors while loading cue files\n")
	}

	errsL := R.LoadGenerators()
	if len(errsL) > 0 {
		for _, e := range errsL {
			util.PrintCueError(e)
		}
		return fmt.Errorf("\nErrors while loading generators\n")
	}

	// issue #20 - Don't print and exit on error here, wait until after we have written, so we can still write good files
	errsG := R.RunGenerators()
	errsW := R.WriteOutput()

	// final timing
	veryend := time.Now()
	elapsed := veryend.Sub(verystart).Round(time.Millisecond)


	R.PrintStats()
	fmt.Printf("\nTotal Elapsed Time: %s\n\n", elapsed)

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


