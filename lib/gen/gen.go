package gen

import (
	"fmt"
	"sync"
	"time"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/cuetils"
	"github.com/hofstadter-io/hof/lib/yagu"
)

func Gen(args []string, rootflags flags.RootPflagpole, cmdflags flags.GenFlagpole) error {
	LT := len(cmdflags.Template)
	LG := len(cmdflags.Generator)
	LW := len(cmdflags.Watch)

	// this might be empty, we calc anyway for ease and sharing
	files, err := yagu.FilesFromGlobs(cmdflags.Watch)
	if err != nil {
		return err
	}

	var errT, errG error

	var wg sync.WaitGroup

	// Run adhoc
	if LT > 0 {

		doRender := func() (chan bool, error) {
			return Render(args, rootflags, cmdflags)
		}

		// no watch, gen once or only watch non-cue
		if LW == 0 {
			_, err := doRender()
			return err
		}

		// otherwise, we are watching for full reload
		wg.Add(1)
		go func() {
			defer wg.Done()
			errT = DoWatch(doRender, files, "adhoc-full", nil)
		}()
	}

	if LT == 0 || LG > 0 {
		// generator modules handled from here on out
		doGen := func() (chan bool, error) {
			return GenOnce(args, rootflags, cmdflags)
		}

		// no watch, gen once or only watch non-cue
		if LW == 0 {
			_, err := doGen()
			return err
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			errG = DoWatch(doGen, files, "hgmod-full", make(chan bool, 2))
		}()
	}

	wg.Wait()

	if errT != nil {
		return errT
	}
	if errG != nil {
		return errG
	}

	return nil
}

func GenOnce(args []string, rootflags flags.RootPflagpole, cmdflags flags.GenFlagpole) (chan bool, error) {
	verystart := time.Now()

	var errs []error

	R := NewRuntime(args, cmdflags)

	errs = R.LoadCue()
	if len(errs) > 0 {
		for _, e := range errs {
			cuetils.PrintCueError(e)
		}
		return nil, fmt.Errorf("\nErrors while loading cue files\n")
	}

	doGen := func() (chan bool, error) {
		R.ClearGenerators()
		R.ExtractGenerators()

		errsL := R.LoadGenerators()
		if len(errsL) > 0 {
			for _, e := range errsL {
				fmt.Println(e)
				// cuetils.PrintCueError(e)
			}
			return nil, fmt.Errorf("\nErrors while loading generators\n")
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
	}

	_, err :=  doGen()
	if err != nil {
		return nil, err
	}

	// return if watch-xcue not set
	if len(cmdflags.WatchXcue) == 0 {
		return nil, nil
	}

	// we need to watch and do our faster regen for template author DX
	files, err := yagu.FilesFromGlobs(cmdflags.WatchXcue)
	if err != nil {
		return nil, err
	}

	quit := make(chan bool, 2)

	go DoWatch(doGen, files, "hgmod-xcue", quit)

	return quit, err
}

