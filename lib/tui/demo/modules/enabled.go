package modules

import (
	"github.com/hofstadter-io/hof/lib/tui/tview"

	"github.com/hofstadter-io/hof/lib/connector"

	"github.com/hofstadter-io/hof/lib/tui/demo/modules/demos"
	"github.com/hofstadter-io/hof/lib/tui/demo/modules/help"
	"github.com/hofstadter-io/hof/lib/tui/demo/modules/home"
	"github.com/hofstadter-io/hof/lib/tui/modules/root"
)

var (
	Module   connector.Connector
	rootView tview.Primitive
)

func Init() {
	rootView = root.New()

	items := []interface{}{
		// primary layout components
		rootView,

		// routable pages
		home.New(),
		help.New(),
		demos.New(),
	}

	conn := connector.New("root")
	conn.Add(items)
	Module = conn

	Module.Connect(Module)
}

func RootView() tview.Primitive {
	return rootView
}
