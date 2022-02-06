package fs

import (
	"fmt"
  "sync"
  "time"

	"cuelang.org/go/cue"
  "github.com/fsnotify/fsnotify"
  "github.com/mattn/go-zglob"

  "github.com/hofstadter-io/hof/flow/context"
  "github.com/hofstadter-io/hof/flow/flow"
)

func init() {
  context.Register("fs.Watch", NewWatch)
}

type Watch struct {}

func NewWatch(val cue.Value) (context.Runner, error) {
  return &Watch{}, nil
}

func (T *Watch) Run(ctx *context.Context) (interface{}, error) {

  // todo, check failure modes, fill, not return error?
  // (in all tasks)
  // do failed message handlings fail the client connection and IRC flow?

	val := ctx.Value

  var globs []string
  var handler cue.Value
  
  ferr := func () error {
    ctx.CUELock.Lock()
    defer func() {
      ctx.CUELock.Unlock()
    }()


    handler = val.LookupPath(cue.ParsePath("handler"))
    if !handler.Exists() {
      return fmt.Errorf("fs.Watch task missing 'handler' field at %s", val.Path())
    }
    if handler.Err() != nil {
      return handler.Err()
    }

    globListVal := val.LookupPath(cue.ParsePath("globs")) 
    if !globListVal.Exists() {
      return fmt.Errorf("fs.Watch task missing 'globs' field at %s", val.Path())
    }
    if globListVal.Err() != nil {
      return globListVal.Err()
    }

    iter, err := globListVal.List()
    if err != nil {
      return err
    }

    for iter.Next() {
      gv := iter.Value()
      if gv.Err() != nil {
        return gv.Err()
      }
      gs, err := gv.String()
      if err != nil {
        return err
      }

      globs = append(globs, gs)
    }

    return nil
  }()
  if ferr != nil {
    return nil, ferr
  }

  files := []string{}
  for _, glob := range globs {
    matches, err := zglob.Glob(glob)
    if err != nil {
      return nil, err
    }
    files = append(files, matches...)
  }

  fmt.Printf("watching %d files\n", len(files))

  debounce := New(time.Millisecond*250)

  watcher, err := fsnotify.NewWatcher()
	if err != nil {
	  return nil, err	
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

        if ctx.DebugTasks {
          fmt.Println("event:", event)
        }

        // todo, fill event into handler
        // v = v.FillPath(cue.ParsePath("event"), event)

				if event.Op&fsnotify.Write == fsnotify.Write {

          debounce(func() {
            if ctx.DebugTasks {
              fmt.Println("start:", event)
            }
            // todo
            // TODO, compile and run pipeline
            v := val.Context().CompileString("{...}")
            v = v.Unify(handler) 

            // fmt.Println(v)

            p, err := flow.NewFlow(ctx, v)
            if err != nil {
              fmt.Println("Error(flow/new):", err)
              return
            }

            err = p.Start()
            if err != nil {
              fmt.Println("Error(flow/run):", err)
              return
            }
            if ctx.DebugTasks {
              fmt.Println("end:", event)
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
      return nil, err	
    }
  }

	<-done

  return nil, nil
}

func New(after time.Duration) func(f func()) {
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
