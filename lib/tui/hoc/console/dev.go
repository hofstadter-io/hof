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
	tui.AddWidgetHandler(C, "/sys/key/A-x", func (e events.Event) {
		C.Clear()
	})
	tui.AddGlobalHandler("/console", func(ev events.Event) {
		d := ev.Data
		switch t := ev.Data.(type) {
		case *events.EventCustom:
			d = t.Data()
		case *events.EventKey:
			d = t.KeyStr
		}

		// consider using regions here, and coloring afterwards, so output doesn't get screwed by the stray [] from the log messages
		line := fmt.Sprintf("[%s] %v", ev.When().Format("2006-01-02 15:04:05"), d)

		level := strings.TrimPrefix(ev.Path, "/console/")
		if len(level) > 6 && level[:6] == "color-" {
			color := level[6:]
			line = fmt.Sprintf("[%s]%.5s  %s[ivory]", color, color, line)
		} else {
			switch level {
			case "crit":
				line = fmt.Sprintf("[#FF00FF]CRIT   %s[ivory]", line)
			case "error":
				line = fmt.Sprintf("[red]ERROR  %s[ivory]", line)
			case "warn":
				line = fmt.Sprintf("[gold]WARN   %s[ivory]", line)
			case "info":
				line = fmt.Sprintf("[ivory]INFO   %s[ivory]", line)
			case "extra":
				line = fmt.Sprintf("[lightskyblue]EXTRA  %s[ivory]", line)
			case "debug":
				line = fmt.Sprintf("[deepskyblue]DEBUG  %s[ivory]", line)
			case "trace":
				line = fmt.Sprintf("[lawngreen]TRACE  %s[ivory]", line)
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
