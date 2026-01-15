package dag

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"dagger.io/dagger"
	"github.com/fsnotify/fsnotify"
)

type Runner func(fsnotify.Event) error

func (d *Dag) StartWatchers() (err error) {
	fmt.Println("watching?", len(d.watched))
	for _, w := range d.watched {
		w.quit = make(chan bool)

		fmt.Println("watching:", w.watch.Path, w.watch.Include, w.watch.Exclude)
		paths, err := d.gatherPaths(w)
		if err != nil {
			d.StopWatchers()
			return err
		}

		c := d.dag.Container()
		w.ctr = c.WithMountedCache("/", w.vol)

		handle := func(evt fsnotify.Event) error {

			fmt.Println("fsnotify:", evt)

			return nil
		}

		go startWatcher(handle, paths, time.Duration(0), w.quit, true)

	}

	return nil
}

func (d *Dag) gatherPaths(idx *hashCacheIndex) (paths []string, err error) {
	paths, err = idx.dir.Glob(d.ctx, "**/*.*")
	if err != nil {
		return nil, err
	}

	if idx.watch.TrimPrefix != "" {
		for i, p := range paths {
			paths[i] = filepath.Join(idx.watch.TrimPrefix, p)
		}
	}

	return paths, nil
}

func (d *Dag) walkDaggerDir(dir *dagger.Directory) ([]string, error) {
	paths := make([]string, 0)

	entries, err := dir.Entries(d.ctx)
	if err != nil {
		return nil, err
	}

	paths = append(paths, entries...)

	for _, entry := range entries {

		if strings.HasSuffix(entry, "/") {
			subdir := dir.Directory(entry)
			subpaths, err := d.walkDaggerDir(subdir)
			if err != nil {
				return nil, err
			}
			paths = append(paths, subpaths...)
		}
	}

	return paths, nil
}

func (d *Dag) StopWatchers() (err error) {
	for _, w := range d.watched {
		if w.quit != nil {
			w.quit <- true
		}
	}
	return nil
}

func startWatcher(Run Runner, files []string, deboundDelay time.Duration, quit chan bool, verbose bool) error {
	// now loop
	// debounce := NewDebouncer(deboundDelay)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		// this is to prevent more than one running at a time
		//var tellDone chan bool
		//tellDone = make(chan bool)

		// watching loop
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					if verbose {
						fmt.Println("event not ok", event)
					}
					break
				}

				// debounce(func() {
				err = Run(event)
				if err != nil {
					fmt.Println("watch.run.error:", err)
				}
				// })

			case err, ok := <-watcher.Errors:
				if err != nil /* && verbose */ {
					fmt.Println("error:", ok, err)
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
	fmt.Printf("watching %d files\n", len(files))

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
