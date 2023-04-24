{{ $DM := .DM }}

/* ============= */

{{ range $K, $M := $DM.Models }}
{{ if ne $K "$hof" }}
/* This is a new table for {{$M.Name}} */
CREATE TABLE {{ snake $M.Name }} (
{{ range $K, $F := $M.Fields }}
{{ if ne $K "$hof" }}{{ template "create-field.sql" (dict "Field" $F "DM" $DM) }}{{ end }}
{{ end }}
)
{{ end }}

{{ end }}

