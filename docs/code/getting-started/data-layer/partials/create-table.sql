{{ $DM := .DM }}
{{ if .SNAP }}{{ $DM = .SNAP.Data }}{{ end }}

{{ range $K, $M := $DM.Models }}
{{ if ne $K "$hof" }}
CREATE TABLE {{ snake $M.Name }} (
{{ range $K, $F := $M.Fields }}
{{ if ne $K "$hof" }}{{ template "create-field.sql" (dict "Field" $F "DM" $DM) }}{{ end }}
{{ end }}
);
{{ end }}
{{ end }}
