{{ snake .Field.Name }} {{ if .Field.sqlType -}}
{{ .Field.sqlType }}{{ else }}{{ .Field.Type -}}{{end}}{{ with .Field.Default }} DEFAULT {{.}}{{ end }}
{{ if .Field.Relation }}
{{ end }}
