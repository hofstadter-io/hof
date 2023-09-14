package flower

import (
	"github.com/gdamore/tcell/v2"

	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/tview"
)

func (F *Flower) setupKeybinds() {
	// events (hotkeys)
	F.SetInputCapture(func(evt *tcell.EventKey) *tcell.EventKey {
		switch evt.Key() {
		case tcell.KeyRune:
			if (evt.Modifiers() & tcell.ModAlt) == tcell.ModAlt {
				switch evt.Rune() {
				case 'f':
					flexDir := F.GetDirection()
					if flexDir == tview.FlexRow {
						F.SetDirection(tview.FlexColumn)
					} else {
						F.SetDirection(tview.FlexRow)
					}
					tui.Draw()

				case 's':
					F.Rebuild()

				case 'S':
					F.Rebuild()

				case 'R':
					F.Rebuild()

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

