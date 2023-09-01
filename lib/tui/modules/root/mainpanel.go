package root

import (
	"fmt"

	"github.com/hofstadter-io/hof/lib/connector"
	"github.com/hofstadter-io/hof/lib/tui/tview"
	"github.com/hofstadter-io/hof/lib/tui/hoc/layouts/panels"
	"github.com/hofstadter-io/hof/lib/tui/hoc/router"
)

type Routables interface {
	Routes() []router.RoutePair
}

type MainPanel struct {
	*panels.Layout

	//
	// Left Panel
	//

	//
	// Main Panel
	//
	mainView *router.Router

	//
	// Right Panel
	//

}

func NewMainPanel() *MainPanel {
	M := &MainPanel{
		Layout: panels.New(),
	}

	M.setupRouter()

	return M
}

func (M *MainPanel) Connect(C connector.Connector) {
	// Get the Routable modules
	rtbls := C.Get((*Routables)(nil))
	for _, Rtbl := range rtbls {
		rtbl := Rtbl.(Routables)
		for _, pair := range rtbl.Routes() {
			M.mainView.AddRoute(pair.Path, pair.Data)
		}
	}
}

func (M *MainPanel) setupRouter() {

	// Set a NotFound View (aka 404 w/o the internet)
	nfv := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	nfv.SetTitle("  Not Found  ").SetBorder(true)
	fmt.Fprint(nfv, "\n\nThe requested path or view does not exist.\n\n")

	M.mainView = router.New()
	M.mainView.SetNotFound(nfv)

	M.SetMainPanel("main-content", M.mainView, 0, 1, 1, "C-m")

}
