package gen

import (
	"fmt"
	"sync"
	"time"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/yagu"
)

func Gen(args []string, rootflags flags.RootPflagpole, cmdflags flags.GenFlagpole) error {
	// return GenLast(args, rootflags, cmdflags)
	verystart := time.Now()

	err := runGen(args, rootflags, cmdflags)
	if err != nil {
		return err
	}

	// final timing
	veryend := time.Now()
	elapsed := veryend.Sub(verystart).Round(time.Millisecond)

	if cmdflags.Stats {
		fmt.Printf("\nGrand Total Elapsed Time: %s\n\n", elapsed)
	}

	return nil
}

func runGen(args []string, rootflags flags.RootPflagpole, cmdflags flags.GenFlagpole) error {
	R := NewRuntime(args, cmdflags)
	err := R.LoadCue()
	if err != nil {
		return err
	}

	LT := len(cmdflags.Template)
	LG := len(cmdflags.Generator)
	globs := cmdflags.WatchGlobs
	xcue := cmdflags.WatchXcue

	// determine watch mode
	//  excplicit: -w
	//  implicit:  -W/-X
	watch := cmdflags.Watch
	if len(globs) > 0 || len(xcue) > 0 {
		watch = true
	}

	/* We will run a generator if either
	   not adhoc or is adhoc with the -G flag
		So let's load them early, there is some helpful info in them
	*/
	if LT == 0 || LG > 0 {
		// load generators just so we can search for watch lists
		err := R.ExtractGenerators()
		if err != nil {
			return err
		}

		errsL := R.LoadGenerators()
		if len(errsL) > 0 {
			for _, e := range errsL {
				fmt.Println(e)
				// cuetils.PrintCueError(e)
			}
			return fmt.Errorf("\nErrors while loading generators\n")
		}
	}

	/* Build up watch list
		We need to buildup the watch list from flags
		and any generator we might run, which might have watch settings
	*/
	for _, G := range R.Generators {
		if G.Disabled {
			continue
		}
		globs = append(globs, G.WatchGlobs...)
		xcue = append(xcue, G.WatchXcue...)
	}

	// this might be empty, we calc anyway for ease and sharing
	wfiles, err := yagu.FilesFromGlobs(globs)
	if err != nil {
		return err
	}
	xfiles, err := yagu.FilesFromGlobs(xcue)
	if err != nil {
		return err
	}

	doGen := func() (chan bool, error) {
		return R.genOnce(watch, xfiles)
	}

	// no watch, gen once or only watch non-cue
	if !watch {
		_, err := doGen()
		return err
	}

	fmt.Printf("found %d glob files from %v\n", len(wfiles), globs)
	fmt.Printf("found %d xcue files from %v\n", len(xfiles), xcue)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = DoWatch(doGen, true, wfiles, "full", make(chan bool, 2))
	}()

	wg.Wait()

	if err != nil {
		return err
	}

	return nil
}

func (R *Runtime) genOnce(watch bool, files []string) (chan bool, error) {
	verystart := time.Now()

	doGen := func() (chan bool, error) {
		fmt.Println("doGen")
		R.ClearGenerators()
		err := R.LoadCue()
		if err != nil {
			return nil, err
		}

		err = R.ExtractGenerators()
		if err != nil {
			return nil, err
		}

		errsL := R.LoadGenerators()
		if len(errsL) > 0 {
			for _, e := range errsL {
				fmt.Println(e)
			}
			return nil, fmt.Errorf("\nErrors while loading generators\n")
		}

		// issue #20 - Don't print and exit on error here, wait until after we have written, so we can still write good files
		errsG := R.RunGenerators()
		errsW := R.WriteOutput()

		// final timing
		veryend := time.Now()
		elapsed := veryend.Sub(verystart).Round(time.Millisecond)

		if R.Flagpole.Stats {
			R.PrintStats()
			fmt.Printf("\nTotal Elapsed Time: %s\n\n", elapsed)
		}

		if len(errsG) > 0 {
			for _, e := range errsG {
				fmt.Println(e)
			}
			return nil, fmt.Errorf("\nErrors while generating output\n")
		}
		if len(errsW) > 0 {
			for _, e := range errsW {
				fmt.Println(e)
			}
			return nil, fmt.Errorf("\nErrors while writing output\n")
		}

		R.PrintMergeConflicts()

		return nil, nil
	} // end doGen

	_, err :=  doGen()
	if err != nil {
		return nil, err
	}

	// return if watching
	if !watch {
		return nil, nil
	}

	quit := make(chan bool, 2)

	go DoWatch(doGen, false, files, "xcue", quit)

	return quit, err
}

