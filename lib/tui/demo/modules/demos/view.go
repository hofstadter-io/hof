package demos

import (
	"fmt"

	"github.com/hofstadter-io/hof/lib/connector"
	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/tview"
	"github.com/hofstadter-io/hof/lib/tui/hoc/router"
)

type Demo interface {
	DemoName() string
	DemoItem() tview.Primitive
}

// Both a Module and a Layout and a Switcher.SubLayout
type Demos struct {
	*tview.TextView
}

func NewDemos() *Demos {
	view := tview.NewTextView()
	view.
		SetTitle("  Demos  ").
		SetBorder(true).
		SetBorderPadding(1, 1, 2, 2)
	view.
		SetWrap(false).
		SetScrollable(true).
		SetDynamicColors(true).
		SetRegions(true)

	D := &Demos{
		TextView: view,
	}

	return D
}

func (D *Demos) Id() string {
	return "demos"
}

func (D *Demos) Name() string {
	return "Demos"
}

func (D *Demos) Connect(connector.Connector) {

}

func (D *Demos) Mount(context map[string]interface{}) error {
	tui.SendCustomEvent("/console/warn", "demos mount")

	D.Clear()
	fmt.Fprintln(D, demosContent())

	return nil
}

func (D *Demos) Refresh(context map[string]interface{}) error {
	D.Clear()
	fmt.Fprintln(D, demosContent())

	return nil
}

func (D *Demos) Routes() []router.RoutePair {
	return []router.RoutePair{
		router.RoutePair{"/demos", D},
		router.RoutePair{"/demos/{demo}", D},
	}
}

func (D *Demos) CommandName() string {
	return "demos"
}

func (D *Demos) CommandUsage() string {
	return "demos <demo>"
}

func (D *Demos) CommandHelp() string {
	return "displays the demos list"
}
func (D *Demos) CommandCallback(args []string, context map[string]interface{}) {
	demo := ""
	if len(args) > 0 {
		demo = args[0]
	}

	if D.IsMounted() {
		if context == nil {
			context = make(map[string]interface{})
		}

		vars := map[string]string{}
		context["demos-args"] = args
		vars["demo"] = demo
		context["vars"] = vars

		D.Refresh(context)
	} else {
		path := "/demos"
		if demo != "" {
			path += "/" + demo
		}

		go tui.SendCustomEvent("/router/dispatch", path)
	}
}
