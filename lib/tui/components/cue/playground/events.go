package playground

import (
	"github.com/gdamore/tcell/v2"

	"github.com/hofstadter-io/hof/lib/tui/tview"
)
func (C *Playground) setupKeybinds() {
	// events (hotkeys)
	C.SetInputCapture(func(evt *tcell.EventKey) *tcell.EventKey {
		switch evt.Key() {
		case tcell.KeyRune:
			if (evt.Modifiers() & tcell.ModAlt) == tcell.ModAlt {
				switch evt.Rune() {
				case 'f':
					flexDir := C.GetDirection()
					if flexDir == tview.FlexRow {
						C.SetDirection(tview.FlexColumn)
					} else {
						C.SetDirection(tview.FlexRow)
					}

				case 'S':
					C.UseScope(!C.useScope)
					C.Rebuild(C.useScope)

				case 'R':
					C.Rebuild(true)

				default: 
					return evt
				}

				return nil
			}

			return evt

		default:
			return evt
		}
	})	
}
