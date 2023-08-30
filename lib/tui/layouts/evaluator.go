package layouts

import (
	"github.com/rivo/tview"
)

type EvaluatorLayout struct {
	*tview.Flex

	Main *tview.Flex
	View tview.Primitive
	Eval tview.Primitive
	Repl tview.Primitive
}

func NewEvaluatorLayout(view, eval, repl tview.Primitive) *EvaluatorLayout {

	layout := &EvaluatorLayout{
		View: view,
		Eval: eval,
		Repl: repl,
	}

	// main layout
	layout.Main = tview.NewFlex()
	layout.Main.
		AddItem(layout.View, 0, 1, false).
		AddItem(layout.Eval, 0, 1, false)

	// page-layout
	layout.Flex = tview.NewFlex().SetDirection(tview.FlexRow)
	layout.Flex.
		AddItem(layout.Main, 0, 2, false).
		AddItem(layout.Repl, 16, 1, false)

	return layout
}
