package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/echo-contrib/echoprometheus"
)

func setupMiddleware(e *echo.Echo) error {
	// ensure request IDs
	e.Use(middleware.RequestID())

	// setup logging middleware
	e.Use(middleware.Logger())

	// setup recovery middleware
	e.Use(middleware.Recover())

	{{ if .SERVER.Auth }}
	// setup auth middleware
	setupAuth(e)
	{{ end }}

	{{ if .SERVER.Prometheus }}
	// setup metrics middleware
	e.Use(echoprometheus.NewMiddleware("{{ .SERVER.Name }}"))
	{{ end }}

	return nil
}

