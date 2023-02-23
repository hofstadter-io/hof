module {{ .SERVER.GoModule }}

go 1.17

require (
	github.com/labstack/echo-contrib v0.11.0
	github.com/labstack/echo/v4 v4.6.1
	github.com/prometheus/client_golang v1.11.0
)
