// Copyright 2012 The Gorilla Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mux

import (
	"fmt"
	"net/url"
	"path"
	"regexp"

	"github.com/pkg/errors"
	"github.com/hofstadter-io/hof/lib/tui/tview"
)

var (
	ErrMethodMismatch = errors.New("method is not allowed")
	ErrNotFound       = errors.New("no matching route was found")
)

// NewRouter returns a new router instance.
func NewRouter() *Router {
	return &Router{namedRoutes: make(map[string]*Route)}
}

// Router registers routes to be matched and dispatches a handler.
//
// It implements the http.Handler interface, so it can be registered to serve
// requests:
//
//     var router = mux.NewRouter()
//
//     func main() {
//         http.Handle("/", router)
//     }
//
// Or, for Google App Engine, register it in a init() function:
//
//     func init() {
//         http.Handle("/", router)
//     }
//
// This will send all incoming requests to the router.
type Router struct {
	// Configurable Handler to be used when no route matches.
	NotFoundHandler Handler

	// Parent route, if this is a subrouter.
	parent parentRoute
	// Routes to be matched, in order.
	routes []*Route
	// Routes by name for URL building.
	namedRoutes map[string]*Route
	// Slice of middlewares to be called after a match is found
	middlewares []middleware
}

// Match attempts to match the given request against the router's registered routes.
//
// If the request matches a route of this router or one of its subrouters the Route,
// Handler, and Vars fields of the the match argument are filled and this function
// returns true.
//
// If the request does not match any of this router's or its subrouters' routes
// then this function returns false. If available, a reason for the match failure
// will be filled in the match argument's MatchErr field. If the match failure type
// (eg: not found) has a registered handler, the handler is assigned to the Handler
// field of the match argument.
func (r *Router) Match(req *Request, match *RouteMatch) bool {
	for _, route := range r.routes {
		if route.Match(req, match) {
			return true
		}
	}

	// Closest match for a router (includes sub-routers)
	if r.NotFoundHandler != nil {
		match.Handler = r.NotFoundHandler
		match.MatchErr = ErrNotFound
		return true
	}

	match.MatchErr = ErrNotFound
	return false
}

// Dispatches the handler registered in the matched route.
//
// When there is a match, the route variables can be retrieved calling
// mux.Vars(request).
func (r *Router) Dispatch(fullpath string, context map[string]interface{}) (tview.Primitive, *Request, error) {
	u, err := url.Parse(fullpath)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "in Router.Dispatch(path): %q\n", fullpath)
	}
	if context == nil {
		context = make(map[string]interface{})
	}
	req := &Request{
		Path:    u.Path,
		Queries: u.Query(),
		Context: context,
	}

	// Clean path to canonical form and redirect.
	if p := cleanPath(req.Path); p != req.Path {
		req.Path = p
	}

	return r.Serve(req)
}

func (r *Router) Serve(req *Request) (tview.Primitive, *Request, error) {
	var match RouteMatch
	var handler Handler
	if r.Match(req, &match) {
		handler = match.Handler
		req = setVars(req, match.Vars)
		req = setCurrentRoute(req, match.Route)
	}
	// fmt.Printf("\n\n\n\nMatch !!!\n\n%+v\n\n%+v\n\n", req, match)

	if handler == nil {
		handler = r.NotFoundHandler
	}

	// finalize dispatch by calling the handler
	return handler.Serve(req)
}

// Get returns a route registered with the given name.
func (r *Router) Get(name string) *Route {
	return r.getNamedRoutes()[name]
}

// GetRoute returns a route registered with the given name. This method
// was renamed to Get() and remains here for backwards compatibility.
func (r *Router) GetRoute(name string) *Route {
	return r.getNamedRoutes()[name]
}

// ----------------------------------------------------------------------------
// parentRoute
// ----------------------------------------------------------------------------

// getNamedRoutes returns the map where named routes are registered.
func (r *Router) getNamedRoutes() map[string]*Route {
	if r.namedRoutes == nil {
		if r.parent != nil {
			r.namedRoutes = r.parent.getNamedRoutes()
		} else {
			r.namedRoutes = make(map[string]*Route)
		}
	}
	return r.namedRoutes
}

// getRegexpGroup returns regexp definitions from the parent route, if any.
func (r *Router) getRegexpGroup() *routeRegexpGroup {
	if r.parent != nil {
		return r.parent.getRegexpGroup()
	}
	return nil
}

func (r *Router) buildVars(m map[string]string) map[string]string {
	if r.parent != nil {
		m = r.parent.buildVars(m)
	}
	return m
}

// ----------------------------------------------------------------------------
// Route factories
// ----------------------------------------------------------------------------

// NewRoute registers an empty route.
func (r *Router) NewRoute() *Route {
	route := &Route{parent: r}
	r.routes = append(r.routes, route)
	return route
}

// Handle registers a new route with a matcher for the URL path.
// See Route.Path() and Route.Handler().
func (r *Router) Handle(path string, handler Handler) *Route {
	return r.NewRoute().Path(path).Handler(handler)
}

// MatcherFunc registers a new route with a custom matcher function.
// See Route.MatcherFunc().
func (r *Router) MatcherFunc(f MatcherFunc) *Route {
	return r.NewRoute().MatcherFunc(f)
}

// Path registers a new route with a matcher for the URL path.
// See Route.Path().
func (r *Router) Path(tpl string) *Route {
	return r.NewRoute().Path(tpl)
}

// PathPrefix registers a new route with a matcher for the URL path prefix.
// See Route.PathPrefix().
func (r *Router) PathPrefix(tpl string) *Route {
	return r.NewRoute().PathPrefix(tpl)
}

// Queries registers a new route with a matcher for URL query values.
// See Route.Queries().
func (r *Router) Queries(pairs ...string) *Route {
	return r.NewRoute().Queries(pairs...)
}

// BuildVarsFunc registers a new route with a custom function for modifying
// route variables before building a URL.
func (r *Router) BuildVarsFunc(f BuildVarsFunc) *Route {
	return r.NewRoute().BuildVarsFunc(f)
}

// Walk walks the router and all its sub-routers, calling walkFn for each route
// in the tree. The routes are walked in the order they were added. Sub-routers
// are explored depth-first.
func (r *Router) Walk(walkFn WalkFunc) error {
	return r.walk(walkFn, []*Route{})
}

// SkipRouter is used as a return value from WalkFuncs to indicate that the
// router that walk is about to descend down to should be skipped.
var SkipRouter = errors.New("skip this router")

// WalkFunc is the type of the function called for each route visited by Walk.
// At every invocation, it is given the current route, and the current router,
// and a list of ancestor routes that lead to the current route.
type WalkFunc func(route *Route, router *Router, ancestors []*Route) error

func (r *Router) walk(walkFn WalkFunc, ancestors []*Route) error {
	for _, t := range r.routes {
		err := walkFn(t, r, ancestors)
		if err == SkipRouter {
			continue
		}
		if err != nil {
			return err
		}
		for _, sr := range t.matchers {
			if h, ok := sr.(*Router); ok {
				ancestors = append(ancestors, t)
				err := h.walk(walkFn, ancestors)
				if err != nil {
					return err
				}
				ancestors = ancestors[:len(ancestors)-1]
			}
		}
		if h, ok := t.handler.(*Router); ok {
			ancestors = append(ancestors, t)
			err := h.walk(walkFn, ancestors)
			if err != nil {
				return err
			}
			ancestors = ancestors[:len(ancestors)-1]
		}
	}
	return nil
}

// ----------------------------------------------------------------------------
// Context
// ----------------------------------------------------------------------------

// RouteMatch stores information about a matched route.
type RouteMatch struct {
	Route   *Route
	Handler Handler
	Vars    map[string]string

	// MatchErr is set to appropriate matching error
	// It is set to ErrMethodMismatch if there is a mismatch in
	// the request method and route method
	MatchErr error
}

type contextKey int

const (
	varsKey  = "vars"
	routeKey = "route"
)

// Vars returns the route variables for the current request, if any.
func Vars(r *Request) map[string]string {
	if rv := r.Context[varsKey]; rv != nil {
		return rv.(map[string]string)
	}
	return nil
}

// CurrentRoute returns the matched route for the current request, if any.
// This only works when called inside the handler of the matched route
// because the matched route is stored in the request context which is cleared
// after the handler returns, unless the KeepContext option is set on the
// Router.
func CurrentRoute(r *Request) *Route {
	if rv := r.Context[routeKey]; rv != nil {
		return rv.(*Route)
	}
	return nil
}

func setVars(r *Request, val interface{}) *Request {
	r.Context[varsKey] = val
	return r
}

func setCurrentRoute(r *Request, val interface{}) *Request {
	r.Context[routeKey] = val
	return r
}

// ----------------------------------------------------------------------------
// Helpers
// ----------------------------------------------------------------------------

// cleanPath returns the canonical path for p, eliminating . and .. elements.
// Borrowed from the net/http package.
func cleanPath(p string) string {
	if p == "" {
		return "/"
	}
	if p[0] != '/' {
		p = "/" + p
	}
	np := path.Clean(p)
	// path.Clean removes trailing slash except for root;
	// put the trailing slash back if necessary.
	if p[len(p)-1] == '/' && np != "/" {
		np += "/"
	}

	return np
}

// uniqueVars returns an error if two slices contain duplicated strings.
func uniqueVars(s1, s2 []string) error {
	for _, v1 := range s1 {
		for _, v2 := range s2 {
			if v1 == v2 {
				return fmt.Errorf("mux: duplicated route variable %q", v2)
			}
		}
	}
	return nil
}

// checkPairs returns the count of strings passed in, and an error if
// the count is not an even number.
func checkPairs(pairs ...string) (int, error) {
	length := len(pairs)
	if length%2 != 0 {
		return length, fmt.Errorf(
			"mux: number of parameters must be multiple of 2, got %v", pairs)
	}
	return length, nil
}

// mapFromPairsToString converts variadic string parameters to a
// string to string map.
func mapFromPairsToString(pairs ...string) (map[string]string, error) {
	length, err := checkPairs(pairs...)
	if err != nil {
		return nil, err
	}
	m := make(map[string]string, length/2)
	for i := 0; i < length; i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m, nil
}

// mapFromPairsToRegex converts variadic string parameters to a
// string to regex map.
func mapFromPairsToRegex(pairs ...string) (map[string]*regexp.Regexp, error) {
	length, err := checkPairs(pairs...)
	if err != nil {
		return nil, err
	}
	m := make(map[string]*regexp.Regexp, length/2)
	for i := 0; i < length; i += 2 {
		regex, err := regexp.Compile(pairs[i+1])
		if err != nil {
			return nil, err
		}
		m[pairs[i]] = regex
	}
	return m, nil
}

// matchInArray returns true if the given string value is in the array.
func matchInArray(arr []string, value string) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}

// matchMapWithString returns true if the given key/value pairs exist in a given map.
func matchMapWithString(toCheck map[string]string, toMatch map[string][]string, canonicalKey bool) bool {
	for k, v := range toCheck {

		if values := toMatch[k]; values == nil {
			return false
		} else if v != "" {
			// If value was defined as an empty string we only check that the
			// key exists. Otherwise we also check for equality.
			valueExists := false
			for _, value := range values {
				if v == value {
					valueExists = true
					break
				}
			}
			if !valueExists {
				return false
			}
		}
	}
	return true
}

// matchMapWithRegex returns true if the given key/value pairs exist in a given map compiled against
// the given regex
func matchMapWithRegex(toCheck map[string]*regexp.Regexp, toMatch map[string][]string, canonicalKey bool) bool {
	for k, v := range toCheck {
		if values := toMatch[k]; values == nil {
			return false
		} else if v != nil {
			// If value was defined as an empty string we only check that the
			// key exists. Otherwise we also check for equality.
			valueExists := false
			for _, value := range values {
				if v.MatchString(value) {
					valueExists = true
					break
				}
			}
			if !valueExists {
				return false
			}
		}
	}
	return true
}
