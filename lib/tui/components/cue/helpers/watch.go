package helpers

import (
	"fmt"
	"time"

	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/watch"
	"github.com/hofstadter-io/hof/lib/yagu"
)


func (sc *SourceConfig) Watch() {
	d := sc.WatchTime

	if d.Nanoseconds() > 0 {
		// startup new watch
		tui.StatusMessage(fmt.Sprintf("start watch for %s", sc.Name))
		err := sc.watch(sc.Name, sc.WatchFunc, d)
		if err != nil {
			tui.Log("error", fmt.Sprintf("watch error in %s: %v", sc.Name, err))
		}

	} else {
		// or stop any watches
		tui.Log("info", fmt.Sprintf("stop watch for %s", sc.Name))
		sc.StopWatch()
	}

}

func (sc *SourceConfig) watch(label string, callback func(), debounce time.Duration) error {
	var (
		files []string
		err error
	)
	if len(sc.WatchGlobs) == 0 {
		switch sc.Source {
		case EvalRuntime:
			if sc._runtime == nil {
				r, err := LoadRuntime(sc.Args)
				if err != nil {
					tui.Log("error", err)
					return err
				}
				sc._runtime = r
			}
			files = sc._runtime.GetLoadedFiles()
		case EvalFile:
			files = sc.Args
		default:
			return fmt.Errorf("auto-file discover not available for %s, you can set globs manually though")
		}
	} else {
		files, err = yagu.FilesFromGlobs(sc.WatchGlobs)
	}
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return fmt.Errorf("did not find any files to watch")
	}

	// always kill old watcher
	sc.StopWatch()

	// make a new runner
	sc.WatchQuit = make(chan bool, 2) // non blocking

	cb := func() error {
		callback()
		return nil
	}

	sc.WatchFunc = callback
	err = watch.Watch(cb, files, label, debounce, sc.WatchQuit, false)

	return err
}

func (sc *SourceConfig) StopWatch() {
	if sc.WatchQuit != nil {
		sc.WatchQuit <- true
		sc.WatchQuit = nil
	}
}
