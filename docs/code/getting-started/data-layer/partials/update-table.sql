{{ $DM := .DM }}
{{ $M := lookup .Name .DM.Models }}
{{ $fields := lookup "+" .Model.Fields }}
{{ if eq (gokind $fields) "map" }}
ALTER TABLE {{ snake .Name }}
	{{ range $l, $f := $fields }}
	{{ $F := lookup $l $M.Fields }}
	{{ template "add-field.sql" (dict "Field" $F "DM" .DM) }},
	{{ end }}
;
{{ end }}
