ADD {{ snake .Field.Name }} {{ if .Field.sqlType -}}
{{ .Field.sqlType }}{{ else }}{{ .Field.Type -}}{{end}}
{{ if .Field.Relation }}
{{ end }}
