{{- range $i, $F := $.RESOURCE.Model.OrderedFields -}}
{{- if (gt $i 0) }}, {{ end }}form.{{ $F.Name -}}.value
{{- end -}}
