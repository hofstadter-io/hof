package home

import (
	"fmt"

	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/tview"
	"github.com/hofstadter-io/hof/lib/tui/hoc/router"
)

// Both a Module and a Layout and a Switcher.SubLayout
type Home struct {
	*tview.TextView
}

func NewHome() *Home {
	view := tview.NewTextView()
	view.
		SetTitle("  Home  ").
		SetBorder(true).
		SetBorderPadding(1, 1, 2, 2)
	view.
		SetWrap(false).
		SetScrollable(true).
		SetDynamicColors(true).
		SetRegions(true)

	fmt.Fprintln(view, homeContent)

	h := &Home{
		TextView: view,
	}

	return h
}

func (H *Home) Id() string {
	return "home"
}

func (H *Home) Routes() []router.RoutePair {
	return []router.RoutePair{
		router.RoutePair{"/", H},
		router.RoutePair{"/home", H},
	}
}

func (H *Home) Name() string {
	return "Home"
}

func (H *Home) HotKey() string {
	return ""
}

func (H *Home) CommandName() string {
	return "home"
}

func (H *Home) CommandUsage() string {
	return "home"
}

func (H *Home) CommandHelp() string {
	return "displays the home view"
}
func (H *Home) CommandCallback(args []string, context map[string]interface{}) {
	go tui.SendCustomEvent("/router/dispatch", "/home")
}
