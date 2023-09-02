package components

import (
	"time"

	"cuelang.org/go/cue"

	"github.com/hofstadter-io/hof/lib/runtime"
	"github.com/hofstadter-io/hof/lib/tui"
	// "github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/tview"
	"github.com/hofstadter-io/hof/lib/watch"
)

type ValueEvaluator struct {
	*tview.Flex

	Runtime *runtime.Runtime

	Edit *tview.TextArea
	View *ValueBrowser

	// that's funky!
	debouncer func(func())
}

func NewValueEvaluator(R *runtime.Runtime) (*ValueEvaluator) {

	VE := &ValueEvaluator{
		Flex: tview.NewFlex(),
		Runtime: R,

	}

	VE.Flex = tview.NewFlex().SetDirection(tview.FlexRow)
	// with two panels

	// TODO, options form

	// editor
	VE.Edit = tview.NewTextArea()
	VE.Edit.
		SetTitle("expression(s)").
		SetBorder(true)

	// results
	VE.View = NewValueBrowser(VE.Runtime.Value, "cue", func(string){})
	VE.View.
		SetTitle("results").
		SetBorder(true)

	// layout
	VE.Flex.
		AddItem(VE.Edit, 0, 1, true).
		AddItem(VE.View, 0, 2, false)


	return VE
}

func (VE *ValueEvaluator) Rebuild() {
	val := VE.Runtime.Value
	ctx := val.Context()

	src := VE.Edit.GetText()

	v := ctx.CompileString(src, cue.Scope(val), cue.InferBuiltins(true))

	VE.View.Value = v
	VE.View.Rebuild(v.Path().String())

	tui.Draw()
}

func (VE *ValueEvaluator) Mount(context map[string]any) error {

	// change debouncer
	VE.debouncer = watch.NewDebouncer(time.Millisecond * 500)
	VE.Edit.SetChangedFunc(func() {
		tui.Log("info", "VE debouncer.setup")

		VE.debouncer(func(){
			tui.Log("warn", "VE debouncer.run")
			VE.Rebuild()
		})

	})

	return nil
}

func (C *ValueEvaluator) Focus(delegate func(p tview.Primitive)) {
	if C.View.HasFocus() {
		delegate(C.View)
		return
	}
	// otherwise, assume we want to keep the view focus
	delegate(C.Edit)
	return
}
