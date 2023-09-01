package help

import (
	"fmt"
	"strings"

	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/tview"
	"github.com/hofstadter-io/hof/lib/tui/hoc/router"
)

// Both a Module and a Layout and a Switcher.SubLayout
type Help struct {
	*tview.TextView
}

func NewHelp() *Help {
	view := tview.NewTextView()
	view.
		SetTitle("  Help  ").
		SetBorder(true).
		SetBorderPadding(1, 1, 2, 2)
	view.
		SetWrap(false).
		SetScrollable(true).
		SetDynamicColors(true).
		SetRegions(true)

	fmt.Fprintln(view, helpContent)

	h := &Help{
		TextView: view,
	}

	return h
}

func (H *Help) Id() string {
	return "help"
}

func (H *Help) Name() string {
	return "Help"
}

func (H *Help) Routes() []router.RoutePair {
	return []router.RoutePair{
		router.RoutePair{"/help", H},
		router.RoutePair{"/help/{topic}", H},
		router.RoutePair{"/help/{topic}/{subTopic}", H},
		router.RoutePair{"/help/{topic}/{subTopic}/{subSubTopic}", H},
		router.RoutePair{"/help/{topic}/{subTopic}/{subSubTopic}/{section}", H},
	}
}

func (H *Help) CommandName() string {
	return "help"
}

func (H *Help) CommandUsage() string {
	return "help <topic> [sub-topics...]"
}

func (H *Help) CommandHelp() string {
	return "displays the home view"
}
func (H *Help) CommandCallback(args []string, context map[string]interface{}) {
	helpPath := "/help"

	if len(args) > 0 {
		H.Clear()
		fmt.Fprintln(H, "Help -", args, "\n\n")
		helpPath += "/" + strings.Join(args, "/")
	} else {
		H.Clear()
		fmt.Fprintln(H, "Help - Main\n\n", helpContent)
	}

	go tui.SendCustomEvent("/router/dispatch", helpPath)
}
