package gen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/yagu"
)

// determine watch mode
//  explicit: -w
//  implicit:  -W/-X
func shouldWatch(cmdflags flags.GenFlagpole) bool {
	return cmdflags.Watch || len(cmdflags.WatchFull) > 0 || len(cmdflags.WatchFast) > 0
}

func buildWatchLists(R *Runtime, args []string, cmdflags flags.GenFlagpole) (wfiles, xfiles []string, err error) {
	if !shouldWatch(cmdflags) {
		return
	}

	// TODO?, when determined to watch
	// add generator templates / partials

	fullWG := cmdflags.WatchFull
	fastWG := cmdflags.WatchFast

	/* Build up watch list
		We need to buildup the watch list from flags
		and any generator we might run, which might have watch settings
	*/

	if R.Verbosity > 1 {
		fmt.Println("Creating Watch List")
	}
	// todo, infer most entrypoints
	for _, arg := range args {
		// skip stdin arg, or args which are filetype specifiers
		if arg == "-" || strings.HasSuffix(arg, ":") {
			continue
		}
		info, err := os.Stat(arg)
		if err != nil {
			return nil, nil, err
		}
		if info.IsDir() {
			fullWG = append(fullWG, filepath.Join(arg, "/*"))
		} else {
			fullWG = append(fullWG, arg)
		}
	}

	for _, G := range R.Generators {
		// we skip when disabled or package is set
		if G.Disabled {
			continue
		}
		basedir := R.CueModuleRoot
		if G.Name == "AdhocGen" {
			basedir = ""
		}

		for _, wfg := range G.WatchFull {
			fullWG = append(fullWG, filepath.Join(basedir,wfg))
		}
		for _, wfg := range G.WatchFast {
			fastWG = append(fastWG, filepath.Join(basedir,wfg))
		}

		// when package is set or not...
		if G.PackageName == "" {
			// when not set, we are probably in the module
			// thus we are in all-in-one mode or module authoring

			// add templates to full regen globs
			// note, we are not recursing here
			// maybe add a CUE field to disable watch
			// if someone wants to recursively watch
			// some generators but not all?
			for _,T := range G.Templates {	
				for _, glob := range T.Globs {
					fastWG = append(fastWG, filepath.Join(basedir,glob))
				}
			}
			for _,P := range G.Partials {
				for _, glob := range P.Globs {
					fastWG = append(fastWG, filepath.Join(basedir,glob))
				}
			}
			for _,S := range G.Statics {
				for _, glob := range S.Globs {
					fastWG = append(fastWG, filepath.Join(basedir,glob))
				}
			}
			// where's your cover sheet? You got the memo right?

		} else {
			// note, the following probably does not belong in a loop
			// globs = append(globs, "./cue.mod/**/*", "*.cue", "design/**/*")

			// otherwise, this is mostly likely an import
			// let's watch the cue.mod vendor directory
			// will we follow symlinks here?
			// will this break down once `cue mod` is a thing...
			//  and modules live outside of the project, in home dir
			//  really an edge case here...
			// for now this is better
		}
	}
	// add partial templates to xcue globs
	// can do outside loop since all gens have the same value
	fastWG = append(fastWG, R.Flagpole.Partial...)

	// this might be empty, we calc anyway for ease and sharing
	wfiles, err = yagu.FilesFromGlobs(fullWG)
	if err != nil {
		return nil, nil, err
	}
	xfiles, err = yagu.FilesFromGlobs(fastWG)
	if err != nil {
		return nil, nil, err
	}

	// if we are in watch mode, let the user know what is being watched
	fmt.Printf("found %d glob files from %v\n", len(wfiles), fullWG)
	fmt.Printf("found %d fastWG files from %v\n", len(xfiles), fastWG)

	return
}

func DoWatch(F func(fast bool) (chan bool, error), dofast, onfirst bool, files []string, label string, quit chan bool) (error) {
	// now loop
	debounce := NewDebouncer(time.Millisecond * 50)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		var tellDone chan bool

		if onfirst {
			// first call
			fmt.Printf("first (%s)\n", label)
			start := time.Now()
			// first time should always be fast
			tellDone, err = F(true)
			end := time.Now()
			elapsed := end.Sub(start).Round(time.Millisecond)
			fmt.Printf(" done (%s) %v\n", label, elapsed)
			if err != nil {
				return
			}
		}

		// watching loop
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					fmt.Println("event not ok", event)
					break
				}

				if event.Op&fsnotify.Write == fsnotify.Write {
					if tellDone != nil {
						tellDone <- true
					}

					debounce(func() {
						// kill previous sub-spawn (xcue)
						fmt.Printf("regen (%s)\n", label)
						start := time.Now()
						tellDone, err = F(dofast)
						end := time.Now()

						elapsed := end.Sub(start).Round(time.Millisecond)
						fmt.Printf(" done (%s) %v\n", label, elapsed)

						if err != nil {
							fmt.Println("error:", err)
						}
					})
				}

			case err, ok := <-watcher.Errors:
				fmt.Println("error:", err)
				if !ok {
					break
				}

			case <-quit:
				return
			}
		}
	}()

	for _, file := range files {
		err = watcher.Add(file)
		if err != nil {
			return err
		}
	}
	// fmt.Printf("watching (%s) %d files\n", label, len(files))

	wg.Wait()

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
