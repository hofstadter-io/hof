import (
	// ...
	
	{{ if gt (len .Resources ) 1 }}
	"{{ .Module }}/resources"
	{{ end }}
)

func setupRouter(e *echo.Echo) error {
	// ...
	{{ range $R := .Resources -}}
	resources.{{ $R.Name }}Routes(g)
	{{ end }}

	return nil
}
