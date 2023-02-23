package types

type {{ .Name }} struct {
	{{ range .Fields }}
	{{ camelT .Name }} {{ .Type }}
	{{ end }}
}
