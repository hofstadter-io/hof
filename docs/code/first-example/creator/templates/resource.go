package resources

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// {{ .RESOURCE.Name }}Routes sets up the routes in a router Group
func {{ .RESOURCE.Name }}Routes(G *echo.Group) {
	fmt.Println("adding {{.RESOURCE.Name}} routes")
	g := G.Group("/{{ kebab .RESOURCE.Name }}")

	// wire up CRUD routes
	{{ range $R := .RESOURCE.Routes }}
	g.{{ $R.Method }}( "{{ range $PATH := $R.Params }}/:{{$PATH}}{{ end }}", {{ $R.Name }}Handler)
	{{- end }}
}

{{ range $R := .RESOURCE.Routes }}
{{ template "handler.go" $R }}
{{ end }}
