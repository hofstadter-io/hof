package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	{{ if gt (len .SERVER.Routes ) 1 }}
	"{{ .Module }}/routes"
	{{ end }}
	{{ if gt (len .Resources ) 1 }}
	"{{ .Module }}/resources"
	{{ end }}
)

func setupRouter(e *echo.Echo) error {

	// Internal routes
	e.GET("/internal/alive", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	{{ if .SERVER.Prometheus }}
	h := promhttp.Handler()
	e.GET("/internal/metrics", func(c echo.Context) error {
		h.ServeHTTP(c.Response(), c.Request())
		return nil
	})
	{{ end }}

	// Explicit routes
	g := e.Group("")

	{{ range $R := .SERVER.Routes -}}
	routes.{{ $R.Name }}Routes(g)
	{{ end }}

	{{ range $R := .Resources -}}
	resources.{{ $R.Name }}Routes(g)
	{{ end }}

	return nil
}
