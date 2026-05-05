package types

{{ range .Input }}
// use a template fragment
{{ template "struct" .}}
{{ end }}

// define a template fragment
{{ define "struct" }}
type {{ .Name }} struct {
	{{ range .Fields }}
	// template from -P flag
	{{ template "field.go" . }}
	{{ end }}
}
{{ end }}
