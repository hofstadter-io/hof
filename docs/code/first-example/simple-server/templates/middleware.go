package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/echo-contrib/prometheus"
)

func setupMiddleware(e *echo.Echo) error {
	// setup recovery middleware
	e.Use(middleware.Recover())

	// setup logging middleware
	e.Use(middleware.Logger())

	{{ if .SERVER.Prometheus }}
	// Setup metrics middleware
	p := prometheus.NewPrometheus("{{ .Server.Name }}", nil)
	e.Use(p.HandlerFunc)
	{{ end }}

	return nil
}

