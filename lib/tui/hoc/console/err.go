package console

import (
	"fmt"

	"github.com/gdamore/tcell/v2"

	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/events"
	"github.com/hofstadter-io/hof/lib/tui/tview"
)

type ErrConsoleWidget struct {
	*tview.TextView
}

func NewErrConsoleWidget() *ErrConsoleWidget {
	textView := tview.NewTextView()
	textView.
		SetTextColor(tcell.ColorMaroon).
		SetScrollable(true).
		SetChangedFunc(func() {
			tui.Draw()
			textView.ScrollToEnd()
		})

	textView.SetTitle(" errors ").
		SetBorder(true).
		SetBorderColor(tcell.ColorRed)

	C := &ErrConsoleWidget{
		TextView: textView,
	}

	return C
}

func (C *ErrConsoleWidget) Mount(context map[string]interface{}) error {

	tui.AddGlobalHandler("/user/error", func(evt events.Event) {
		str := evt.Data.(*events.EventCustom).Data()
		text := fmt.Sprintf("[%s] %v\n", evt.When().Format("2006-01-02 15:04:05"), str)
		fmt.Fprintf(C, "%s", text)
	})

	tui.AddGlobalHandler("/sys/err", func(ev events.Event) {
		err := ev.Data.(*events.EventError)
		line := fmt.Sprintf("[%s] %v", ev.When().Format("2006-01-02 15:04:05"), err)
		fmt.Fprintf(C, "[red]SYSERR %v[white]\n", line)
	})

	return nil
}
func (C *ErrConsoleWidget) Unmount() error {
	tui.RemoveWidgetHandler(C, "/user/error")
	tui.RemoveWidgetHandler(C, "/sys/err")
	return nil
}
