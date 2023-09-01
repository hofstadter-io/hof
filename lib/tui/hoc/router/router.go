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
		// this should always be a map[string]any (i.e. the context)
		d := ev.Data.(*events.EventCustom).Data()
		path := ""

		switch d := d.(type) {
		case string:
		  path = d
		case map[string]any:
			path = d["path"].(string)
		}

		go tui.SendCustomEvent("/console/info", fmt.Sprintf("router dispatching: %q %#v", path, ev.Data))

		context := map[string]interface{}{
			"activation": "dispatch",
			"path":       path,
			"data":       ev.Data,
			"event":      ev,
		}
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
		go tui.SendCustomEvent("/console/error", fmt.Sprintf("handle %#v", req))
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
	layout, req, err := R.iRouter.Dispatch(path, context)
	if err != nil {
		go tui.SendCustomEvent("/console/error", fmt.Errorf("in dispatch handler: %w", err))
	}
	if layout != nil {
		ctx := req.Context
		req.Context = nil
		ctx["req"] = req
		R.setActive(layout, ctx)
	} else {
		go tui.SendCustomEvent("/console/error", "nil layout in dispatch handler")
	}
}

func (R *Router) setActive(layout tview.Primitive, context map[string]interface{}) {
	// R.Pages.SwitchToPage(layout.Id())
	R.Pages.SwitchToPage(layout.Id(), context)
	tui.Draw()
}
