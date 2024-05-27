{{ $ROUTE := . }}
func {{ $ROUTE.Name }}Handler(c echo.Context) (err error) {

	// process any path and query params
	{{ range $P := $ROUTE.Params }}
	{{ $P }} := c.Param("{{ $P }}")
	{{ end }}

	{{ range $Q := $ROUTE.Query }}
	{{ $Q }} := c.QueryParam("{{ $Q }}")
	{{ end }}

	{{ if $ROUTE.Body }}
	// custom body
	{{ $ROUTE.Body }}
	{{ else }}
	// default body
	c.String(http.StatusNotImplemented, "Not Implemented")
	{{ end }}

	return nil
}
