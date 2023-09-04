package eval

import (
	"fmt"

	"github.com/gdamore/tcell/v2"

	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/tview"
)

// This is probably now what Wrap*Handlers is in tview, and Panel might now be a different kind of component, since others embed and extend it
func setupEventHandlers(
	P *Panel,
	customKey func(event *tcell.EventKey) *tcell.EventKey,
	customMouse func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse),
) {

	P.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse){

		if customMouse != nil {
			action, event = customMouse(action, event)
			if event == nil {
				return action, nil
			}
		}
		if i := P.ChildFocus(); i >= 0 {
			return action, event
		}
		if action > 0 {
			tui.Log("trace", fmt.Sprintf("Panel.mouseHandler %v %v %v %d", P.Id(), action, event.Buttons(), P.Flex.ChildFocus()))
			// return action, nil
		}

		return action, event
	})

	P.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		if customKey != nil {
			event = customKey(event)
			if event == nil {
				return nil
			}
		}
		
		mostFocusedPanel := P.GetMostFocusedPanel()
		// unfortunately, these events travel top-down, rather than bottom up...
		// we should probably rework the the call stack, so that...
		//   1. this same travel to the child until no more child focus happens
		//   2. we then unwind, running any user defined handlers, so most child element can override parent handlers (by returning nil)

		// if a child has focus, let's pass it on to them
		focus := P.Flex.HasFocus()
		child := P.Flex.ChildFocus()
		count := P.Flex.GetItemCount()
		tui.Log("trace", fmt.Sprintf("Panel.inputHandler.1 at %d of %d | %v %v", child, count, focus, P.Id()))

		alt := event.Modifiers() & tcell.ModAlt == tcell.ModAlt
		ctrl := event.Modifiers() & tcell.ModCtrl == tcell.ModCtrl
		meta := event.Modifiers() & tcell.ModMeta == tcell.ModMeta
		shift := event.Modifiers() & tcell.ModShift == tcell.ModShift

		ctx := make(map[string]any)

		// we only care about ALT+... keys at this level
		tui.Log("trace", fmt.Sprintf("Panel.inputHandler.2 %v %v %v %v %v %q %v", P.Id(), alt, ctrl, meta, shift, string(event.Rune()), event.Key()))

		switch event.Key() {

		case tcell.KeyESC:
			p := P.GetMostFocusedPanel()
			child := p.ChildFocus()
			count := p.GetItemCount()
			pp := "-1"
			if p._parent != nil {
				pp = p._parent.Id()
				tui.SetFocus(p._parent)
			}
			tui.Log("trace", fmt.Sprintf("Panel want's to blur at %d of %d | %v | %v", child, count, p.Id(), pp))

			//if child < 0 {
			//  tui.SetFocus(P)
			//  return nil
			//}

		case tcell.KeyRune:
			handled := false
			if alt {
				handled = true
				switch event.Rune() {
				// lowercase = nav
				// upsercase = move/insert

				// left, prev
				case 'h':
					ctx["action"] = "nav"
					ctx["where"] = "left"
				case 'H':
					ctx["action"] = "move"
					ctx["where"] = "prev"

				// down, prev
				case 'j':
					ctx["action"] = "nav"
					ctx["where"] = "down"
				case 'J':
					ctx["action"] = "insert"
					ctx["where"] = "prev"

				// up, next
				case 'k':
					ctx["action"] = "nav"
					ctx["where"] = "up"
				case 'K':
					ctx["action"] = "insert"
					ctx["where"] = "next"

				// up, right
				case 'l':
					ctx["action"] = "nav"
					ctx["where"] = "right"
				case 'L':
					ctx["action"] = "move"
					ctx["where"] = "next"

				default:
					handled = false
				}
			}

			mid := mostFocusedPanel.Id()

			if !handled && alt {
				switch event.Rune() {

				case 'T':
					ctx["action"] = "split"
					ctx["id"] = mid
					ctx["panel"] = mostFocusedPanel

				case 'D':
					ctx["action"] = "delete" // DELETE
					ctx["id"] = mid
					ctx["panel"] = mostFocusedPanel

				// flip flex orientation
				case 'F':
					mostFocusedPanel.FlipFlexDirection()
					return nil

				// dev stuff
				case 'v':
					focus := mostFocusedPanel.ChildFocus()
					count := mostFocusedPanel.GetItemCount()
					tui.Log("trace", fmt.Sprintf("Panel(%s).focus at %v of %v", mid, focus, count))
					return nil

				default:
					return event

				}
				handled = true
			}

			if handled {
				// ???
				ctx["index"] = P.Flex.ChildFocus()

				P.Refresh(ctx)
				return nil
			}

		}

		return event
	})
}
