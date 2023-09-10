package eval

import (
	"fmt"

	"github.com/gdamore/tcell/v2"

	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/events"
)

// This is probably now what Wrap*Handlers is in tview, and Panel might now be a different kind of component, since others embed and extend it
func (M *Eval) setupEventHandlers() {

	//
	// Our message bus system (which probably needs some updating for nested handling
	//

	// handle border display
	tui.AddWidgetHandler(M.Panel, "/sys/key/A-P", func(e events.Event) {
		if M.HasFocus() {
			M.showPanel = !M.showPanel
			M.SetShowBordersR(M.showPanel, M.showOther)
		}
	})

	tui.AddWidgetHandler(M.Panel, "/sys/key/A-O", func(e events.Event) {
		if M.HasFocus() {
			M.showOther = !M.showOther
			M.SetShowBordersR(M.showPanel, M.showOther)
		}
	})

	//
	// tview.Primitive scoped key input handling
	//
	M.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		alt := event.Modifiers() & tcell.ModAlt == tcell.ModAlt
		ctrl := event.Modifiers() & tcell.ModCtrl == tcell.ModCtrl
		meta := event.Modifiers() & tcell.ModMeta == tcell.ModMeta
		shift := event.Modifiers() & tcell.ModShift == tcell.ModShift

		ctx := make(map[string]any)

		// we only care about ALT+... keys at this level
		tui.Log("trace", fmt.Sprintf("Eval.inputHandler %v %v %v %v %q %v", alt, ctrl, meta, shift, string(event.Rune()), event.Key()))
		// tui.Log("warn", fmt.Sprintf("Eval.keyInput %v %v %v", alt, event.Key(), string(event.Rune())))

		panel := M.GetMostFocusedPanel()
		if panel != nil {
			ctx["panel"] = panel 
			ctx["panel-id"] = panel.Id()
			ctx["child-focus-index"] = panel.ChildFocus()
		}

		handled := false
		eRefresh := false
		switch event.Key() {

		// same comment about items & widgets with multiple things (also applies to the nav.* options under Alt-<rune>
		case tcell.KeyUp:
			if ctrl {
				ctx["action"] = "nav.up"
				handled = true
				eRefresh = true
			}
		case tcell.KeyDown:
			if ctrl {
				ctx["action"] = "nav.down"
				handled = true
				eRefresh = true
			}
		case tcell.KeyLeft:
			if ctrl {
				ctx["action"] = "nav.left"
				handled = true
				eRefresh = true
			}
		case tcell.KeyRight:
			if ctrl {
				ctx["action"] = "nav.right"
				handled = true
				eRefresh = true
			}


		case tcell.KeyRune:
			if alt {
				handled = true
				switch event.Rune() {
				// lowercase = nav
				// upsercase = move/insert

				// left, prev
				case 'h':
					ctx["action"] = "nav.left"
					eRefresh = true
				case 'H':
					ctx["action"] = "move"
					ctx["where"] = "prev"

				// down, prev
				case 'j':
					ctx["action"] = "nav.down"
					eRefresh = true
				case 'J':
					ctx["action"] = "insert"
					ctx["where"] = "prev"
					ctx["item"] = "help"

				// up, next
				case 'k':
					ctx["action"] = "nav.up"
					eRefresh = true
				case 'K':
					ctx["action"] = "insert"
					ctx["where"] = "next"
					ctx["item"] = "help"

				// up, right
				case 'l':
					ctx["action"] = "nav.right"
					eRefresh = true
				case 'L':
					ctx["action"] = "move"
					ctx["where"] = "next"

				default:
					handled = false
				}
			}

			// mid := panel.Id()

			if !handled && alt {
				switch event.Rune() {

				case 'T':
					ctx["action"] = "split"
					ctx["item"] = "help"

				case 'D':
					ctx["action"] = "delete" // DELETE

				// flip flex orientation
				case 'F':
					panel.FlipFlexDirection()
					return nil

				// dev stuff
				case 'v':
					focus := panel.ChildFocus()
					count := panel.GetItemCount()
					tui.Log("trace", fmt.Sprintf("Panel(%s).focus at %v of %v", panel.Id(), focus, count))
					return nil

				default:
					return event

				}
				handled = true
			}
		} // end key / rune switch

		if handled {
			if eRefresh {
				M.Refresh(ctx)
			} else {
				panel.Refresh(ctx)
			}
			return nil
		}

		return event
	})
}
