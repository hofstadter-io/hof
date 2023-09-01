// Copyright 2012 The Gorilla Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mux

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// Route stores information to match a request and build URLs.
type Route struct {
	// Parent where the route was registered (a Router).
	parent parentRoute
	// Request handler for the route.
	handler Handler
	// List of matchers.
	matchers []matcher
	// Manager for the variables from host and path.
	regexp *routeRegexpGroup
	// If true, this route never matches: it is only used to build URLs.
	buildOnly bool
	// The name used to build URLs.
	name string
	// Error resulted from building a route.
	err error

	buildVarsFunc BuildVarsFunc
}

// Match matches the route against the request.
func (r *Route) Match(req *Request, match *RouteMatch) bool {
	if r.buildOnly || r.err != nil {
		return false
	}

	var matchErr error

	// Match everything.
	for _, m := range r.matchers {
		if matched := m.Match(req, match); !matched {
			matchErr = nil
			return false
		}
	}

	if matchErr != nil {
		match.MatchErr = matchErr
		return false
	}

	if match.MatchErr == ErrMethodMismatch {
		// We found a route which matches request method, clear MatchErr
		match.MatchErr = nil
		// Then override the mis-matched handler
		match.Handler = r.handler
	}

	// Yay, we have a match. Let's collect some info about it.
	if match.Route == nil {
		match.Route = r
	}
	if match.Handler == nil {
		match.Handler = r.handler
	}
	if match.Vars == nil {
		match.Vars = make(map[string]string)
	}

	// Set variables.
	if r.regexp != nil {
		r.regexp.setMatch(req, match, r)
	}
	return true
}

// ----------------------------------------------------------------------------
// Route attributes
// ----------------------------------------------------------------------------

// GetError returns an error resulted from building the route, if any.
func (r *Route) GetError() error {
	return r.err
}

// BuildOnly sets the route to never match: it is only used to build URLs.
func (r *Route) BuildOnly() *Route {
	r.buildOnly = true
	return r
}

// Handler --------------------------------------------------------------------

func (R *Route) HandlerFunc(HandlerFunc) {

}

// Handler sets a handler for the route.
func (r *Route) Handler(handler Handler) *Route {
	if r.err == nil {
		r.handler = handler
	}
	return r
}

// GetHandler returns the handler for the route, if any.
func (r *Route) GetHandler() Handler {
	return r.handler
}

// Name -----------------------------------------------------------------------

// Name sets the name for the route, used to build URLs.
// If the name was registered already it will be overwritten.
func (r *Route) Name(name string) *Route {
	if r.name != "" {
		r.err = fmt.Errorf("mux: route already has name %q, can't set %q",
			r.name, name)
	}
	if r.err == nil {
		r.name = name
		r.getNamedRoutes()[name] = r
	}
	return r
}

// GetName returns the name for the route, if any.
func (r *Route) GetName() string {
	return r.name
}

// ----------------------------------------------------------------------------
// Matchers
// ----------------------------------------------------------------------------

// matcher types try to match a request.
type matcher interface {
	Match(*Request, *RouteMatch) bool
}

// addMatcher adds a matcher to the route.
func (r *Route) addMatcher(m matcher) *Route {
	if r.err == nil {
		r.matchers = append(r.matchers, m)
	}
	return r
}

// addRegexpMatcher adds a host or path matcher and builder to a route.
func (r *Route) addRegexpMatcher(tpl string, typ regexpType) error {
	if r.err != nil {
		return r.err
	}
	r.regexp = r.getRegexpGroup()
	if typ == regexpTypePath || typ == regexpTypePrefix {
		if len(tpl) > 0 && tpl[0] != '/' {
			return fmt.Errorf("mux: path must start with a slash, got %q", tpl)
		}
		if r.regexp.path != nil {
			tpl = strings.TrimRight(r.regexp.path.template, "/") + tpl
		}
	}
	rr, err := newRouteRegexp(tpl, typ, routeRegexpOptions{})
	if err != nil {
		return err
	}
	for _, q := range r.regexp.queries {
		if err = uniqueVars(rr.varsN, q.varsN); err != nil {
			return err
		}
	}
	if typ == regexpTypeQuery {
		r.regexp.queries = append(r.regexp.queries, rr)
	} else {
		r.regexp.path = rr
	}

	r.addMatcher(rr)
	return nil
}

// MatcherFunc ----------------------------------------------------------------

// MatcherFunc is the function signature used by custom matchers.
type MatcherFunc func(*Request, *RouteMatch) bool

// Match returns the match for a given request.
func (m MatcherFunc) Match(r *Request, match *RouteMatch) bool {
	return m(r, match)
}

// MatcherFunc adds a custom function to be used as request matcher.
func (r *Route) MatcherFunc(f MatcherFunc) *Route {
	return r.addMatcher(f)
}

// Path -----------------------------------------------------------------------

// Path adds a matcher for the URL path.
// It accepts a template with zero or more URL variables enclosed by {}. The
// template must start with a "/".
// Variables can define an optional regexp pattern to be matched:
//
// - {name} matches anything until the next slash.
//
// - {name:pattern} matches the given regexp pattern.
//
// For example:
//
//     r := mux.NewRouter()
//     r.Path("/products/").Handler(ProductsHandler)
//     r.Path("/products/{key}").Handler(ProductsHandler)
//     r.Path("/articles/{category}/{id:[0-9]+}").
//       Handler(ArticleHandler)
//
// Variable names must be unique in a given route. They can be retrieved
// calling mux.Vars(request).
func (r *Route) Path(tpl string) *Route {
	r.err = r.addRegexpMatcher(tpl, regexpTypePath)
	return r
}

// PathPrefix -----------------------------------------------------------------

// PathPrefix adds a matcher for the URL path prefix. This matches if the given
// template is a prefix of the full URL path. See Route.Path() for details on
// the tpl argument.
//
// Note that it does not treat slashes specially ("/foobar/" will be matched by
// the prefix "/foo") so you may want to use a trailing slash here.
//
// Also note that the setting of Router.StrictSlash() has no effect on routes
// with a PathPrefix matcher.
func (r *Route) PathPrefix(tpl string) *Route {
	r.err = r.addRegexpMatcher(tpl, regexpTypePrefix)
	return r
}

// Query ----------------------------------------------------------------------

// Queries adds a matcher for URL query values.
// It accepts a sequence of key/value pairs. Values may define variables.
// For example:
//
//     r := mux.NewRouter()
//     r.Queries("foo", "bar", "id", "{id:[0-9]+}")
//
// The above route will only match if the URL contains the defined queries
// values, e.g.: ?foo=bar&id=42.
//
// It the value is an empty string, it will match any value if the key is set.
//
// Variables can define an optional regexp pattern to be matched:
//
// - {name} matches anything until the next slash.
//
// - {name:pattern} matches the given regexp pattern.
func (r *Route) Queries(pairs ...string) *Route {
	length := len(pairs)
	if length%2 != 0 {
		r.err = fmt.Errorf(
			"mux: number of parameters must be multiple of 2, got %v", pairs)
		return nil
	}
	for i := 0; i < length; i += 2 {
		if r.err = r.addRegexpMatcher(pairs[i]+"="+pairs[i+1], regexpTypeQuery); r.err != nil {
			return r
		}
	}

	return r
}

// BuildVarsFunc --------------------------------------------------------------

// BuildVarsFunc is the function signature used by custom build variable
// functions (which can modify route variables before a route's URL is built).
type BuildVarsFunc func(map[string]string) map[string]string

// BuildVarsFunc adds a custom function to be used to modify build variables
// before a route's URL is built.
func (r *Route) BuildVarsFunc(f BuildVarsFunc) *Route {
	r.buildVarsFunc = f
	return r
}

// Subrouter ------------------------------------------------------------------

// Subrouter creates a subrouter for the route.
//
// It will test the inner routes only if the parent route matched. For example:
//
//     r := mux.NewRouter()
//     s := r.Host("www.example.com").Subrouter()
//     s.HandleFunc("/products/", ProductsHandler)
//     s.HandleFunc("/products/{key}", ProductHandler)
//     s.HandleFunc("/articles/{category}/{id:[0-9]+}"), ArticleHandler)
//
// Here, the routes registered in the subrouter won't be tested if the host
// doesn't match.
func (r *Route) Subrouter() *Router {
	router := &Router{parent: r}
	r.addMatcher(router)
	return router
}

// ----------------------------------------------------------------------------
// URL building
// ----------------------------------------------------------------------------

// Substitute builds the full path from the template for the route.
//
// It accepts a sequence of key/value pairs for the route variables. For
// example, given this route:
//
//     r := mux.NewRouter()
//     r.HandleFunc("/articles/{category}/{id:[0-9]+}", ArticleHandler).
//       Name("article")
//
// ...a Substitution for it can be built using:
//
//     url, err := r.Get("article").Substitute("category", "technology", "id", "42")
//
// ...which will return an url.URL with the following path:
//
//     "/articles/technology/42"
//
// This also works for host variables:
//
//     r := mux.NewRouter()
//     r.HandleFunc("/articles/{category}/{id:[0-9]+}?edit={edit}", ArticleHandler).
//       Name("article")
//
//     // url.String() will be "/articles/technology/42?edit=true"
//     url, err := r.Get("article")
//                  .Substitute("category", "technology",
//                              "id", "42",
//                              "edit", "true")
//
// All variables defined in the route are required, and their values must
// conform to the corresponding patterns.
func (r *Route) Substitute(pairs ...string) (string, error) {
	if r.err != nil {
		return "", r.err
	}
	if r.regexp == nil {
		return "", errors.New("mux: route doesn't have a path")
	}
	values, err := r.prepareVars(pairs...)
	if err != nil {
		return "", err
	}
	var path string
	if r.regexp.path != nil {
		if path, err = r.regexp.path.url(values); err != nil {
			return "", err
		}
	}
	queries := make([]string, 0, len(r.regexp.queries))
	for _, q := range r.regexp.queries {
		var query string
		if query, err = q.url(values); err != nil {
			return "", err
		}
		queries = append(queries, query)
	}

	queryStr := strings.Join(queries, "&")
	fullPath := path + "?" + queryStr

	return fullPath, nil
}

// GetPathTemplate returns the template used to build the
// route match.
// This is useful for building simple REST API documentation and for instrumentation
// against third-party services.
// An error will be returned if the route does not define a path.
func (r *Route) GetPathTemplate() (string, error) {
	if r.err != nil {
		return "", r.err
	}
	if r.regexp == nil || r.regexp.path == nil {
		return "", errors.New("mux: route doesn't have a path")
	}
	return r.regexp.path.template, nil
}

// GetPathRegexp returns the expanded regular expression used to match route path.
// This is useful for building simple REST API documentation and for instrumentation
// against third-party services.
// An error will be returned if the route does not define a path.
func (r *Route) GetPathRegexp() (string, error) {
	if r.err != nil {
		return "", r.err
	}
	if r.regexp == nil || r.regexp.path == nil {
		return "", errors.New("mux: route does not have a path")
	}
	return r.regexp.path.regexp.String(), nil
}

// GetQueriesRegexp returns the expanded regular expressions used to match the
// route queries.
// This is useful for building simple REST API documentation and for instrumentation
// against third-party services.
// An empty list will be returned if the route does not have queries.
func (r *Route) GetQueriesRegexp() ([]string, error) {
	if r.err != nil {
		return nil, r.err
	}
	if r.regexp == nil || r.regexp.queries == nil {
		return nil, errors.New("mux: route doesn't have queries")
	}
	var queries []string
	for _, query := range r.regexp.queries {
		queries = append(queries, query.regexp.String())
	}
	return queries, nil
}

// GetQueriesTemplates returns the templates used to build the
// query matching.
// This is useful for building simple REST API documentation and for instrumentation
// against third-party services.
// An empty list will be returned if the route does not define queries.
func (r *Route) GetQueriesTemplates() ([]string, error) {
	if r.err != nil {
		return nil, r.err
	}
	if r.regexp == nil || r.regexp.queries == nil {
		return nil, errors.New("mux: route doesn't have queries")
	}
	var queries []string
	for _, query := range r.regexp.queries {
		queries = append(queries, query.template)
	}
	return queries, nil
}

// prepareVars converts the route variable pairs into a map. If the route has a
// BuildVarsFunc, it is invoked.
func (r *Route) prepareVars(pairs ...string) (map[string]string, error) {
	m, err := mapFromPairsToString(pairs...)
	if err != nil {
		return nil, err
	}
	return r.buildVars(m), nil
}

// ----------------------------------------------------------------------------
// parentRoute
// ----------------------------------------------------------------------------

// parentRoute allows routes to know about parent host and path definitions.
type parentRoute interface {
	getNamedRoutes() map[string]*Route
	getRegexpGroup() *routeRegexpGroup
	buildVars(map[string]string) map[string]string
}

// getNamedRoutes returns the map where named routes are registered.
func (r *Route) getNamedRoutes() map[string]*Route {
	if r.parent == nil {
		// During tests router is not always set.
		r.parent = NewRouter()
	}
	return r.parent.getNamedRoutes()
}

// getRegexpGroup returns regexp definitions from this route.
func (r *Route) getRegexpGroup() *routeRegexpGroup {
	if r.regexp == nil {
		if r.parent == nil {
			// During tests router is not always set.
			r.parent = NewRouter()
		}
		regexp := r.parent.getRegexpGroup()
		if regexp == nil {
			r.regexp = new(routeRegexpGroup)
		} else {
			// Copy.
			r.regexp = &routeRegexpGroup{
				path:    regexp.path,
				queries: regexp.queries,
			}
		}
	}
	return r.regexp
}

func (r *Route) buildVars(m map[string]string) map[string]string {
	if r.parent != nil {
		m = r.parent.buildVars(m)
	}
	if r.buildVarsFunc != nil {
		m = r.buildVarsFunc(m)
	}
	return m
}
