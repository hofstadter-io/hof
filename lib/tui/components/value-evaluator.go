package components

import (
	"time"

	"cuelang.org/go/cue"
	"github.com/gdamore/tcell/v2"

	// "github.com/hofstadter-io/hof/lib/cuetils"
	"github.com/hofstadter-io/hof/lib/runtime"
	"github.com/hofstadter-io/hof/lib/tui"
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
	VE.View = NewValueBrowser(VE.Runtime.Value, func(string){})
	VE.View.
		SetTitle("results").
		SetBorder(true)
	VE.View.SetMode("code")

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
	tui.SendCustomEvent("/console/warn", "eval updated")

	tui.Draw()
}

func (VE *ValueEvaluator) Mount(context map[string]any) error {

	// change debouncer
	VE.debouncer = watch.NewDebouncer(time.Millisecond * 300)
	VE.Edit.SetChangedFunc(func() {
		tui.SendCustomEvent("/console/warn", "edit changed")

		VE.debouncer(func(){
			VE.Rebuild()
		})

	})



	// key handlers
	VE.Edit.SetInputCapture(func(evt *tcell.EventKey) *tcell.EventKey {
		tui.SendCustomEvent("/console/warn", "edit capture")

		switch evt.Key() {
		case tcell.KeyRune:
			switch evt.Rune() {
			default:
				return evt
			}
		default:
			return evt
		}

		// VB.Rebuild("")

		return nil
	})

	return nil
}
