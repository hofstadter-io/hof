package pages

import (
	"strings"
	"github.com/rivo/tview"

	"github.com/hofstadter-io/hof/lib/tui/app"
	"github.com/hofstadter-io/hof/lib/tui/components"
	"github.com/hofstadter-io/hof/lib/tui/layouts"
)

type EvalPage struct {
	*Page

	App *app.App

	View *components.ValueBrowser
	Eval *components.ValueEvaluator
	Repl *components.Shell
}

func NewEvalPage(app *app.App) *EvalPage {
	page := &EvalPage{
		Page: new(Page),
		App:  app,
	}

	// setup shell
	page.Repl = components.NewShell(app)
	app.Logger = page.Repl.Append

	// setup file browser
	onNodeSelect := func(path string) {
		// app.Logger("onNodeSelect: " + path)
	}
	page.View = components.NewValueBrowser(app, app.Runtime.Value, onNodeSelect)
	// file browser
	n := app.Runtime.Value.Path().String()
	page.View.Rebuild(n)



	page.View.SetTitle(strings.Join(page.App.Runtime.Entrypoints, " "))

	page.Eval = components.NewValueEvaluator(app)


	page.Name = "Value Browser"
	layout := layouts.NewEvaluatorLayout(page.View, page.Eval, page.Repl)
	page.Flex = layout.Flex

	return page
}

func (P *EvalPage) Focus(delegate func(p tview.Primitive)) {
	delegate(P.View)
}
