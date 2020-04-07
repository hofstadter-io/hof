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

	errs = R.LoadGenerators()
	if len(errs) > 0 {
		for _, e := range errs {
			util.PrintCueError(e)
		}
		return fmt.Errorf("\nErrors while loading generators\n")
	}

	errs = R.RunGenerators()
	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Println(e)
		}
		return fmt.Errorf("\nErrors while generating output\n")
	}

	// wait to print error as this is the last thing
	errs = R.WriteOutput()

	// final timing
	veryend := time.Now()
	elapsed := veryend.Sub(verystart).Round(time.Millisecond)


	R.PrintStats()
	fmt.Printf("\nTotal Elapsed Time: %s\n\n", elapsed)

	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Println(e)
		}
		return fmt.Errorf("\nErrors while writing output\n")
	}
	R.PrintMergeConflicts()

	return nil
}


