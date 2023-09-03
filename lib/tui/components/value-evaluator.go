package components

import (
	"time"

	"cuelang.org/go/cue"
	"github.com/gdamore/tcell/v2"

	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/tview"
	"github.com/hofstadter-io/hof/lib/watch"
)

type ValueEvaluator struct {
	*tview.Flex

	Value cue.Value

	Edit *tview.TextArea
	View *ValueBrowser

	flexDir  int
	useScope bool

	// that's funky!
	debouncer func(func())
}

func NewValueEvaluator(val cue.Value) (*ValueEvaluator) {

	C := &ValueEvaluator{
		Flex: tview.NewFlex(),
		Value: val,
		flexDir: tview.FlexRow,
		useScope: true,
	}

	C.Flex = tview.NewFlex().SetDirection(C.flexDir)
	// with two panels

	// TODO, options form

	// editor
	C.Edit = tview.NewTextArea()
	C.Edit.
		SetTitle("expression(s)").
		SetBorder(true)

	// results
	C.View = NewValueBrowser(C.Value, "cue", func(string){})
	C.View.
		SetTitle("results").
		SetBorder(true)
	C.View.UsingScope = C.useScope

	// layout
	C.Flex.
		AddItem(C.Edit, 0, 1, true).
		AddItem(C.View, 0, 1, false)

	// events (hotkeys)
	C.SetInputCapture(func(evt *tcell.EventKey) *tcell.EventKey {
		switch evt.Key() {
		case tcell.KeyRune:
			if (evt.Modifiers() & tcell.ModAlt) == tcell.ModAlt {
				switch evt.Rune() {
				case 'f':
					if C.flexDir == tview.FlexRow {
						C.flexDir = tview.FlexColumn
					} else {
						C.flexDir = tview.FlexRow
					}
					C.Flex.SetDirection(C.flexDir)

				case 's':
					C.useScope = !C.useScope
					C.View.UsingScope = C.useScope
					C.Rebuild()

				default: 
					tui.Log("trace", "val-eval bypassing " + string(evt.Rune()))
					return evt
				}

				tui.Draw()
				return nil
			}

			return evt

		default:
			return evt
		}

		tui.Draw()
		return nil
	})	

	return C
}

func (C *ValueEvaluator) SetScope(visible bool) {
	C.useScope = visible
}

func (C *ValueEvaluator) Rebuild() {
	val := C.Value
	ctx := val.Context()

	src := C.Edit.GetText()

	var v cue.Value
	if C.useScope {
		v = ctx.CompileString(src, cue.InferBuiltins(true), cue.Scope(val))
	} else {
		v = ctx.CompileString(src, cue.InferBuiltins(true))
	}

	C.View.Value = v
	C.View.Rebuild(v.Path().String())

	tui.Draw()
}

func (C *ValueEvaluator) Mount(context map[string]any) error {
	// setup debouncer
	C.debouncer = watch.NewDebouncer(time.Millisecond * 500)

	// trigger rebuild on editor changes
	C.Edit.SetChangedFunc(func() {
		C.debouncer(func(){
			C.Rebuild()
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
