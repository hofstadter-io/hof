package console

import (
	"fmt"
	"strings"

	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/events"
	"github.com/hofstadter-io/hof/lib/tui/tview"
)

type DevConsoleWidget struct {
	*tview.TextView
}

func NewDevConsoleWidget() *DevConsoleWidget {
	textView := tview.NewTextView().
		SetScrollable(true).
		SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(func() {
			tui.Draw()
		})

	textView.SetTitle(" console ").SetBorder(true)

	C := &DevConsoleWidget{
		TextView: textView,
	}

	return C
}

func (C *DevConsoleWidget) Mount(context map[string]interface{}) error {
	tui.AddGlobalHandler("/console", func(ev events.Event) {
		d := ev.Data
		switch t := ev.Data.(type) {
		case *events.EventCustom:
			d = t.Data()
		case *events.EventKey:
			d = t.KeyStr
		}
		line := fmt.Sprintf("[%s] %v", ev.When().Format("2006-01-02 15:04:05"), d)

		level := strings.TrimPrefix(ev.Path, "/console/")
		if len(level) > 6 && level[:6] == "color-" {
			color := level[6:]
			line = fmt.Sprintf("[%s]%.5s  %s[white]", color, color, line)
		} else {
			switch level {
			case "crit":
				line = fmt.Sprintf("[red]CRIT   %s[white]", line)
			case "error":
				line = fmt.Sprintf("[orange]ERROR  %s[white]", line)
			case "warn":
				line = fmt.Sprintf("[yellow]WARN   %s[white]", line)
			case "info":
				line = fmt.Sprintf("INFO  %s", line)
			case "debug":
				line = fmt.Sprintf("[green]DEBUG  %s[white]", line)
			case "trace":
				line = fmt.Sprintf("[aqua]TRACE  %s[white]", line)
			}
		}

		fmt.Fprintln(C, line)
		C.ScrollToEnd()
	})

	return nil
}

func (C *DevConsoleWidget) Unmount() error {
	tui.RemoveWidgetHandler(C, "/console")
	return nil
}
