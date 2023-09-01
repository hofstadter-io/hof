package router

import (
	"fmt"
	"errors"


	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/events"
	"github.com/hofstadter-io/hof/lib/tui/mux"
	"github.com/hofstadter-io/hof/lib/tui/tview"
)

type RoutePair struct {
	Path  string
	Thing interface{}
}

type Routable interface {
	Routings() []RoutePair
}

type Router struct {
	*tview.Pages

	// internal router
	iRouter *mux.Router
}

func New() *Router {
	r := &Router{
		Pages:   tview.NewPages(),
		iRouter: mux.NewRouter(),
	}

	tui.AddWidgetHandler(r.Pages, "/router/dispatch", func(ev events.Event) {
		d := ev.Data.(*events.EventCustom).Data()
		path := ""

		var context map[string]any

		switch d := d.(type) {
		// just the path as a string
		case string:
		  path = d
			context = map[string]any{
				"path": path,
			}

		// a context with the path set
		case map[string]any:
			context = d
			path = context ["path"].(string)
		}

		go tui.SendCustomEvent("/console/info", fmt.Sprintf("router dispatching: %q %#v", path, context))
		r.SetActive(path, context)
	})

	return r
}

func (R *Router) SetNotFound(layout tview.Primitive) {
	handler := func(req *mux.Request) (tview.Primitive, *mux.Request, error) {
		return layout, req, nil
	}
	R.iRouter.NotFoundHandler = mux.NewDefaultHandler(handler)
	R.AddPage(layout.Id(), layout, true, false)
}

func (R *Router) AddRoute(path string, thing interface{}) error {

	switch t := thing.(type) {
	case tview.Primitive:
		R.AddRouteLayout(path, t)

	case mux.HandlerFunc:
		R.AddRouteHandlerFunc(path, t)

	case mux.Handler:
		R.AddRouteHandler(path, t)

	default:
		return errors.New("Unknown thing to be routed to...")
	}

	return nil
}

func (R *Router) AddRouteLayout(path string, layout tview.Primitive) error {
	R.AddPage(layout.Id(), layout, true, false)
	handler := func(req *mux.Request) (tview.Primitive, *mux.Request, error) {
		// go tui.SendCustomEvent("/console/error", fmt.Sprintf("handle %#v", req))
		//ctx := make(map[string]any)
		//layout.Mount(ctx)
		return layout, req, nil
	}
	R.iRouter.Handle(path, mux.NewDefaultHandler(handler))
	return nil
}

func (R *Router) AddRouteHandlerFunc(path string, handler mux.HandlerFunc) error {
	R.iRouter.Handle(path, mux.NewDefaultHandler(handler))
	return nil
}

func (R *Router) AddRouteHandler(path string, handler mux.Handler) error {
	R.iRouter.Handle(path, handler)
	return nil
}

func (R *Router) SetActive(path string, context map[string]interface{}) {
	layout, _, err := R.iRouter.Dispatch(path, context)
	if err != nil {
		go tui.SendCustomEvent("/console/error", fmt.Errorf("in dispatch handler: %w", err))
	}
	if layout != nil {
		R.setActive(layout, context)
	} else {
		go tui.SendCustomEvent("/console/error", "nil layout in dispatch handler")
	}
}

func (R *Router) setActive(layout tview.Primitive, context map[string]interface{}) {
	go tui.SendCustomEvent("/console/warn", fmt.Sprintf("setActive: %v", context))
	R.Pages.SwitchToPage(layout.Id(), context)
	tui.Draw()
}
