{{ $ROUTE := . }}
func {{ $ROUTE.Name }}Handler(c echo.Context) (err error) {
	// hello

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
	fmt.Println(
		{{ range $ROUTE.Params }}{{ . }},{{ end }}
		{{ range $ROUTE.Query }}{{ . }},{{ end }}
	)
	c.String(http.StatusNotImplemented, "Not Implemented")
	{{ end }}

	return nil
}
