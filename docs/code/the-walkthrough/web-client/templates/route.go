package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// {{ .ROUTE.Name }}Routes sets up the routes in a router Group
func {{ .ROUTE.Name }}Routes(G *echo.Group) {
	g := G.Group("{{ .ROUTE.Path }}{{ range $PATH := .ROUTE.Params }}/:{{$PATH}}{{ end }}")
	g.{{.ROUTE.Method}}( "", {{.ROUTE.Name}}Handler)

	// we'll handle sub-routes here in "full-example"
}

// Handler implementation is in a partial template
{{ template "handler.go" .ROUTE }}
