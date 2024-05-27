package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo-contrib/echoprometheus"

	{{ if gt (len .SERVER.Routes ) 1 }}
	"{{ .SERVER.GoModule }}/routes"
	{{ end }}
)

func setupRouter(e *echo.Echo) error {

	// Internal routes
	e.GET("/internal/alive", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	{{ if .SERVER.Prometheus }}
	e.GET("/internal/metrics", echoprometheus.NewHandler())
	{{ end }}

	// Application routes group
	g := e.Group("")

	// Register the routes
	{{ range $R := .SERVER.Routes -}}
	routes.{{ $R.Name }}Routes(g)
	{{ end }}

	return nil
}
