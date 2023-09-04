package components

import (
	"time"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"github.com/gdamore/tcell/v2"

	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/tview"
	"github.com/hofstadter-io/hof/lib/watch"
)

type ValueEvaluator struct {
	*tview.Flex

	Value  cue.Value
	Scope  cue.Value

	// TODO
	ValueGetter func() cue.Value
	ScopeGetter func() cue.Value

	Edit *tview.TextArea
	View *ValueBrowser
	// TODO, make Views...

	flexDir  int
	useScope bool

	// that's funky!
	debouncer func(func())
}

func NewValueEvaluator(src string, val, scope cue.Value) (*ValueEvaluator) {

	C := &ValueEvaluator{
		Flex: tview.NewFlex(),
		Value: val,
		Scope: scope,
		flexDir: tview.FlexRow,
	}

	C.Flex = tview.NewFlex().SetDirection(C.flexDir)
	// with two panels

	// TODO, options form

	// editor
	C.Edit = tview.NewTextArea()
	C.Edit.
		SetTitle("expression(s)").
		SetBorder(true)

	if src != "" {
		C.Edit.SetText(src, false)
	}

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

	C.setupKeybinds()
	C.setupMousebinds()
	return C
}

func (C *ValueEvaluator) SetScope(visible bool) {
	C.useScope = visible
}

func (C *ValueEvaluator) Rebuild(context map[string]any) {
	val := C.Value
	var ctx *cue.Context
	if val.Exists() {
		ctx = val.Context()
	} else {
		ctx = cuecontext.New()
	}

	src := C.Edit.GetText()
	if src == "" && val.Exists() {
		C.View.Value = C.Value
		// we really need everything here (like from Runtime.Value)
		// C.Edit.SetText(fmt.Sprintf("%#v", C.Value), false)
	} else {
		var v cue.Value
		if C.useScope {
			if !C.Scope.Exists() {
				tui.Log("error", "no scope set to activate")
				tui.Tell("error", "no scope set to activate")
				return
			} else {
				s := C.Scope
				if C.ScopeGetter != nil {
					s = C.ScopeGetter()
				}
				// tui.Log("warn", fmt.Sprintf("%#v", s))
				ctx := s.Context()
				v = ctx.CompileString(src, cue.InferBuiltins(true), cue.Scope(s))
			}
		} else {
			v = ctx.CompileString(src, cue.InferBuiltins(true))
		}

		// only update view value, that way, if we erase everything, we still see the value
		C.View.Value = v
		// C.Value = v
	}

	C.View.Rebuild("")

	tui.Draw()
}

func (C *ValueEvaluator) Mount(context map[string]any) error {
	// setup debouncer
	C.debouncer = watch.NewDebouncer(time.Millisecond * 500)

	// trigger rebuild on editor changes
	C.Edit.SetChangedFunc(func() {
		C.debouncer(func(){
			C.Rebuild(nil)
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


func (C *ValueEvaluator) setupKeybinds() {
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

				case 'S':
					if !C.Scope.Exists() {
						tui.Log("error", "no scope set to activate")
						tui.Tell("error", "no scope set to activate")
						return nil
					} else {
						C.useScope = !C.useScope
						C.View.UsingScope = C.useScope
						C.Rebuild(nil)
					}

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
}

func (C *ValueEvaluator) setupMousebinds() {
	// events (hotkeys)
	//handle := func (P tview.Primitive) (func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse)) {
	//  return func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
	//    mx, my := event.Position()	
	//    x,y,w,h := P.GetRect()
	//    // if action != tview.MouseMove {
	//      // tui.Log("trace", fmt.Sprintf("%v %d,%d %d,%d %d,%d", action, mx,my, x,y, x+h,y+w))
	//    // }
	//    if action == tview.MouseLeftDoubleClick {
	//      if my == y && (x <= mx && mx <= x+w) {
	//        tui.Log("warn", fmt.Sprintf("bang! %d", h))
	//        if h == 2 {
	//          C.Flex.ResizeItem(P, 0, 1)
	//        } else {
	//          C.Flex.ResizeItem(P, 2, 0)
	//        }
	//      }
	//    }
	//    return action, event
	//  }
	//}
	//C.Edit.SetMouseCapture(handle(C.Edit))	
	//C.View.SetMouseCapture(handle(C.View))	
}
