{{- range $i, $F := $.RESOURCE.Model.OrderedFields -}}
{{- if (gt $i 0) }}, {{ end }}{{ $F.Name -}}
{{- end -}}
