package main

import (
	"crypto/subtle"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func setupAuth(e *echo.Echo) {
	{{ if .SERVER.Auth.apikey }}
	e.Use(middleware.KeyAuth(func(key string, c echo.Context) (bool, error) {
		if subtle.ConstantTimeCompare([]byte(key), []byte("valid-key")) == 1 {
			return true, nil
		}
		return false, nil
	}))
	{{ end }}

	{{ if .SERVER.Auth.basic }}
	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
			// Be careful to use constant time comparison to prevent timing attacks
			if subtle.ConstantTimeCompare([]byte(username), []byte("joe")) == 1 &&
					subtle.ConstantTimeCompare([]byte(password), []byte("secret")) == 1 {
					return true, nil
			}
			return false, nil
	}))
	{{ end }}
}
