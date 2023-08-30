package components

import (
	"fmt"
	"time"

	"cuelang.org/go/cue"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/hofstadter-io/hof/lib/cuetils"
	"github.com/hofstadter-io/hof/lib/tui/app"
	"github.com/hofstadter-io/hof/lib/watch"
)

type ValueEvaluator struct {
	*tview.Flex

	App *app.App

	Edit *tview.TextArea
	View *ValueBrowser

	// that's funky!
	debouncer func(func())
}

func NewValueEvaluator(app *app.App) (*ValueEvaluator) {
	root := &ValueEvaluator{
		App: app,
	}

	root.Flex = tview.NewFlex().SetDirection(tview.FlexRow)
	// with two panels

	// TODO, options form

	// editor
	root.Edit = tview.NewTextArea()
	root.Edit.
		SetTitle("expression(s)").
		SetBorder(true)

	// results
	root.View = NewValueBrowser(app, app.Runtime.Value, func(string){})
	root.View.
		SetTitle("results").
		SetBorder(true)

	// layout
	root.Flex.
		AddItem(root.Edit, 0, 1, true).
		AddItem(root.View, 0, 2, false)


	// change debouncer
	root.debouncer = watch.NewDebouncer(time.Millisecond * 500)
	root.Edit.SetChangedFunc(func() {

		// root.App.Logger(".")

		root.debouncer(func(){
			// root.App.Logger("!")
			ctx := root.App.Runtime.CueContext
			val := root.App.Runtime.Value
			src := root.Edit.GetText()

			v := ctx.CompileString(src, cue.Scope(val))

			if v.Err() != nil {
				s := cuetils.CueErrorToString(v.Err())		
				root.App.Logger(s)
			} else {
				root.View.Value = v
				root.View.Rebuild(v.Path().String())
				// root.App.Logger("$")
			}
			root.App.Draw()
		})

	})



	// key handlers
	root.Edit.SetInputCapture(func(ek *tcell.EventKey) *tcell.EventKey {

		root.App.Logger(fmt.Sprintln("edit: ", ek.Key() == tcell.KeyCtrlSpace, (ek.Key() == tcell.KeyCR), ek.Modifiers() & tcell.ModCtrl))

		// we aren't handling this
		return ek
	})

	return root
}
