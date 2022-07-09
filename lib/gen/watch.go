package gen

import (
	"fmt"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

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
