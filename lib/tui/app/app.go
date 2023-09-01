package app

import (
	"fmt"
	"os"
	"sync"

	"github.com/hofstadter-io/hof/lib/tui/events"
	"github.com/hofstadter-io/hof/lib/tui/tview"
)

type App struct {
	*tview.Application
	appLock sync.RWMutex

	rootView tview.Primitive

	EventBus *events.EventBus

	// TODO core config, redux like data store?
}

func NewApp() (*App, error) {
	app := &App{
		Application: tview.NewApplication(),
		EventBus:    new(events.EventBus),
	}

	err := app.EventBus.Init(app.Application)
	if err != nil {
		return app, err
	}

	app.EnableMouse(true)

	// is this needed? or do we keep it for easy debug of event bus
	// app.EventBus.AddGlobalHandler("/", func(e events.Event){})

	// common, so you can send a message to redraw
	app.EventBus.AddGlobalHandler("/sys/redraw", func(e events.Event) {
		app.Draw()
	})

	return app, nil
}

func (app *App) GetRootView() (root tview.Primitive) {
	return app.rootView
}

func (app *App) SetRootView(root tview.Primitive) {
	app.rootView = root
}

//func (A *App) AddPage(name string, item tview.Primitive, resize, visible bool ) {
//  A.Pages.AddPage(name, item, resize, visible)
//}

// blocking call
func (app *App) Start() error {

	// catch panics, clean up, format error
	defer func() {
		e := recover()
		if e != nil {
			app.stop()
			// Print a formatted panic output
			fmt.Fprintf(os.Stderr, "Captured a panic(value=%v) lib.Start()... Exit vermui and clean terminal...\nPrint stack trace:\n\n", e)
			//debug.PrintStack()
			//gs, err := stack.ParseDump(bytes.NewReader(debug.Stack()), os.Stderr)
			//if err != nil {
			//  debug.PrintStack()
			//  os.Exit(1)
			//}
			//p := &stack.Palette{}
			//buckets := stack.SortBuckets(stack.Bucketize(gs, stack.AnyValue))
			//srcLen, pkgLen := stack.CalcLengths(buckets, false)
			//for _, bucket := range buckets {
			//  io.WriteString(os.Stdout, p.BucketHeader(&bucket, false, len(buckets) > 1))
			//  io.WriteString(os.Stdout, p.StackLines(&bucket.Signature, srcLen, pkgLen, false))
			//}
			panic(e)
		}
	}()

	// start the event engine
	go app.EventBus.Start()

	err := app.rootView.Mount(nil)
	if err != nil {
		panic(err)
	}

	// blocking
	app.SetRoot(app.rootView, true)
	return app.Run()
}

// Close finalizes vermui library,
// should be called after successful initialization when vermui's functionality isn't required anymore.
func (app *App) stop() error {
	app.Stop()
	err := app.EventBus.Stop()
	if err != nil {
		return err
	}
	return nil
}

func (app *App) Clear() {
	screen := app.Screen()
	if screen != nil {
		screen.Clear()
		screen.Sync()
	}
}
