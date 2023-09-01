package modules

import (
	"github.com/hofstadter-io/hof/lib/tui/tview"

	"github.com/hofstadter-io/hof/lib/connector"

	// base modules
	"github.com/hofstadter-io/hof/lib/tui/modules/root"
	"github.com/hofstadter-io/hof/lib/tui/modules/home"
	"github.com/hofstadter-io/hof/lib/tui/modules/help"

	// core modules
	"github.com/hofstadter-io/hof/lib/tui/modules/eval"

	// extra modules
	"github.com/hofstadter-io/hof/lib/tui/modules/ls"

)

var (
	Conn   connector.Connector
	rootView tview.Primitive
)

func Init() {
	rootView = root.New()

	items := []interface{}{
		// primary / root layout component
		rootView,

		// base modules
		home.New(),
		help.New(),

		// core modules
		eval.New(),

		// extra modules
		ls.New(),
	}

	conn := connector.New("root")
	conn.Add(items)
	Conn = conn

	Conn.Connect(Conn)
}

func RootView() tview.Primitive {
	return rootView
}
