package gen

import (
	"fmt"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/cuetils"
	"github.com/hofstadter-io/hof/lib/yagu"
)

func Gen(args []string, rootflags flags.RootPflagpole, cmdflags flags.GenFlagpole) error {
	// always generate at startup
	err := GenOnce(args, rootflags, cmdflags)
	if err != nil {
		return err
	}

	// if no watches, then only wanted to gen once
	if len(cmdflags.Watch) == 0 {
		return nil
	}

	// otherwise in watch mode
	files, err := yagu.FilesFromGlobs(cmdflags.Watch)
	if err != nil {
		return err
	}

	debounce := NewDebouncer(time.Millisecond * 50)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Op&fsnotify.Write == fsnotify.Write {

					debounce(func() {
						fmt.Printf("watch: starting regen... ")
						start := time.Now()
						derr := GenOnce(args, rootflags, cmdflags)
						end := time.Now()

						elapsed := end.Sub(start).Round(time.Millisecond)
						fmt.Printf("done  %v\n", elapsed)

						if derr != nil {
							fmt.Println("error:", err)
						}
					})
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Println("error:", err)
			}
		}
	}()


	for _, file := range files {
		err = watcher.Add(file)
		if err != nil {
			return err
		}
	}
	fmt.Printf("watching %d files\n", len(files))

	<-done

	return nil
}

func GenOnce(args []string, rootflags flags.RootPflagpole, cmdflags flags.GenFlagpole) error {
	verystart := time.Now()

	var errs []error

	if len(cmdflags.Template) > 0 {
		// Todo, just construct generator here? 
		// and add to CRT below?
		// there might be more going on in there than we want to deal with right now to merge
		err := Render(args, rootflags, cmdflags)
		if err != nil {
			errs = append(errs, err)
		}
		if len(cmdflags.Generator) == 0 {
			var err error
			if len(errs) > 0 {
				for _, e := range errs {
					cuetils.PrintCueError(e)
				}
				err = fmt.Errorf("\nErrors during adhoc gen\n")
			}
			if cmdflags.Stats {
				veryend := time.Now()
				elapsed := veryend.Sub(verystart).Round(time.Millisecond)
				fmt.Printf("\nTotal Elapsed Time: %s\n\n", elapsed)
			}
			return err
		}
	}


	R := NewRuntime(args, cmdflags)

	errs = R.LoadCue()
	if len(errs) > 0 {
		for _, e := range errs {
			cuetils.PrintCueError(e)
		}
		return fmt.Errorf("\nErrors while loading cue files\n")
	}

	// unless len(-T) > 0 && len(-G) == 0
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

func NewDebouncer(after time.Duration) func(f func()) {
	d := &debouncer{after: after}

	return func(f func()) {
		d.add(f)
	}
}

type debouncer struct {
	mu    sync.Mutex
	after time.Duration
	timer *time.Timer
}

func (d *debouncer) add(f func()) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		d.timer.Stop()
	}
	d.timer = time.AfterFunc(d.after, f)
}
