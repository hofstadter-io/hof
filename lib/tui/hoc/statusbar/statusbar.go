package statusbar

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"

	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/events"
	"github.com/hofstadter-io/hof/lib/tui/tview"
)

const emptyMsg = "press 'Ctrl-<space>' to enter a command or '/path/to/something' to navigate"

type StatusBar struct {
	*tview.TextView

	curr    string   // current input (potentially partial)
	hIdx    int      // where we are in history
	history []string // command history
}

func New() *StatusBar {
	S := &StatusBar{
		history: []string{},
	}

	status := tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignLeft)
	status.SetTitle(" Status ").SetTitleAlign(tview.AlignRight).SetBorder(true).SetBorderPadding(0, 0, 1, 0)

	S.TextView = status

	return S
}

func (S *StatusBar) Mount(context map[string]interface{}) error {
	//tui.AddWidgetHandler(S, "/sys/key/C-s", func(e events.Event) {
	//  S.SetBorderColor(tcell.ColorFuchsia)
	//  tui.SetFocus(S.TextView)
	//})
	S.SetDoneFunc(func(key tcell.Key) {
		S.SetBorderColor(tcell.ColorWhite)
		tui.Unfocus()
	})

	tui.AddWidgetHandler(S, "/user/error", func(evt events.Event) {
		str := fmt.Sprintf("[red]%v[white]", evt.Data.(*events.EventCustom).Data())

		S.Clear()
		fmt.Fprint(S, str)
		tui.Draw()

		go func() {
			time.Sleep(time.Second * 6)
			text := S.GetText(false)
			if text == str {
				S.Clear()
				fmt.Fprint(S, "[lime]ok[white]")
				tui.Draw()
			}
		}()
	})

	tui.AddWidgetHandler(S, "/status/message", func(evt events.Event) {
		str := evt.Data.(*events.EventCustom).Data().(string)
		S.history = append(S.history, str)

		S.Clear()
		fmt.Fprint(S, str)
		tui.Draw()

		go func() {
			time.Sleep(time.Second * 6)
			text := S.GetText(false)
			if text == str {
				S.Clear()
				fmt.Fprint(S, "[lime]ok[white]")
				tui.Draw()
			}
		}()
	})

	return nil
}
func (S *StatusBar) Unmount() error {

	tui.RemoveWidgetHandler(S, "/sys/key/C-s")
	tui.RemoveWidgetHandler(S, "/user/error")
	tui.RemoveWidgetHandler(S, "/status/message")

	return nil
}

// InputHandler returns the handler for this primitive.
func (S *StatusBar) InputHandler() func(*tcell.EventKey, func(tview.Primitive)) {
	return S.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {

		dist := 1

		// Process key evt.
		switch key := event.Key(); key {

		// Upwards, back in history
		case tcell.KeyHome:
			dist = len(S.history)
			fallthrough
		case tcell.KeyPgUp:
			dist += 4
			fallthrough
		case tcell.KeyUp: // Regular character.
			S.hIdx -= dist
			if S.hIdx < 0 {
				S.hIdx = 0
			}

		// Downwards, more recent in history
		case tcell.KeyEnd:
			dist = len(S.history)
			fallthrough
		case tcell.KeyPgDn:
			dist += 4
			fallthrough
		case tcell.KeyDown:
			S.hIdx += dist
			if S.hIdx >= len(S.history) {
				S.hIdx = len(S.history) - 1
			}

		}

		str := ""
		if len(S.history) > 0 {
			str = S.history[S.hIdx]
		}
		S.Clear()
		fmt.Fprint(S, str)
		tui.Draw()
	})
}
