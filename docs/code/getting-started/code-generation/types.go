package types

{{ range .Input }}
type {{ .Name }} struct {
	{{ range .Fields -}}
	{{ camelT .Name }} {{ .Type }}
	{{ end }}
}
{{ end }}
