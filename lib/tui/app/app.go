package app

import (
	"fmt"
	"os"
	"sync"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/hofstadter-io/hof/lib/tui/events"
	"github.com/hofstadter-io/hof/lib/tui/tview"
)

type App struct {
	*tview.Application
	appLock sync.RWMutex

	rootView tview.Primitive
	lastFocus tview.Primitive

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

func (app *App) GetLastFocus() (last tview.Primitive) {
	return app.lastFocus
}

func (app *App) SetLastFocus(last tview.Primitive) {
	app.lastFocus = last
}

// blocking call
func (app *App) Start(context map[string]any) error {
	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		return err
	}
	defer terminal.Restore(0, oldState)

	// stuff to ensure we don't mess up the user's terminal
	// catch panics, clean up, format error
	defer func() {
		e := recover()
		if e != nil {
			app.stop()
			terminal.Restore(0, oldState)
			fmt.Fprintf(os.Stderr, "Captured a panic(value=%v) lib.Start()... Exiting and cleaning terminal...\nPrint stack trace:\n\n", e)
			panic(e)
		}
	}()

	// start the event engine
	go app.EventBus.Start()

	// set the initial view
	err = app.rootView.Mount(context)
	if err != nil {
		return err
	}

	app.SetRoot(app.rootView, true)

	// blocking
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

func (app *App) DoStop() {
	app.stop()
}

func (app *App) Clear() {
	screen := app.Screen()
	if screen != nil {
		screen.Clear()
		screen.Sync()
	}
}
