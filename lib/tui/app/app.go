package app

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/hofstadter-io/hof/lib/runtime"
)

type App struct {
	*tview.Application

	loglevel string

	Logger func(string)

	Runtime *runtime.Runtime

	Pages *tview.Pages

	helpShown bool
	HelpModal tview.Primitive
}

func NewApp() *App {
	app := &App{
		Application: tview.NewApplication(),
		loglevel: "info",
	}

	app.EnableMouse(true)


	text := tview.NewTextView()
	text.
		SetBorder(true).
		SetTitle("   Help   ")
	fmt.Fprint(text, helpMsg)	

	width, height := 100, 30

	app.HelpModal = tview.
		NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(text, height, 1, true).
			AddItem(nil, 0, 1, false), width, 1, true).
		AddItem(nil, 0, 1, false)

	app.Pages = tview.NewPages().
		AddPage("app-help", app.HelpModal, true, false)

	// setup keys
	app.SetInputCapture(func(ek *tcell.EventKey) *tcell.EventKey {
		app.Logger(fmt.Sprintf("app: %q %v %v\n", ek.Rune(), ek.Key(), ek.Modifiers()))

		// enum keys
		switch ek.Key() {
		case tcell.KeyEscape:
			if app.helpShown {
				app.Pages.HidePage("app-help")
				app.helpShown = false
				return nil
			}
			// we aren't handling this
			return ek

		}

		// rune keys
		switch ek.Rune() {
		case '?':
			app.helpShown = true
			app.Pages.ShowPage("app-help")
			app.Pages.SendToFront("app-help")
			return nil
		}
		// we aren't handling this
		return ek
	})

	return app
}

const helpMsg = `
Hof TUI help system

<space>, <enter>:         open / close a value

C - dive into a value    (make it the new root)
u - back out of a value  (unwind the vals to root)

j, down arrow, right arrow: Move (the selection) down by one node.
k, ⬆⇦: Move (the selection) up by one node.
g, home: Move (the selection) to the top.
G, end: Move (the selection) to the bottom.

J:											Move (the selection) up one level
K: 											Move (the selection) to the last node one level down
Ctrl-F, page down:			Move (the selection) down by one page.
Ctrl-B, page up:				Move (the selection) up by one page.

? - this help

`
