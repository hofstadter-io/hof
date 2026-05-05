package types

// this range used to have .Input
{{ range . }}
type {{ .Name }} struct {
	{{ range .Fields -}}
	{{ camelT .Name }} {{ .Type }}
	{{ end }}
}
{{ end }}
