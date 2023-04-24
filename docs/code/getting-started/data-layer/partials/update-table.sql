/* updating table */
{{ $DM := .DM }}
{{ $DIFF := .DM.CurrDiff }}

{{ if .SNAP }}
{{ $DM = .SNAP.Data }}
{{ $DIFF = .SNAP.CurrDiff }}
{{ end }}

{{ range $K, $m := $DIFF.Models }}
{{ $M := lookup $K $DM.Models }}
{{ if ne $K "$hof" }}
{{ range $K, $fields := $m.Fields }}
{{ if eq $K "+"}}
	{{ range $l, $f := $fields }}
	{{ $F := lookup $l $M.Fields }}
ALTER TABLE {{ snake $M.Name }}
	{{ template "add-field.sql" (dict "Field" $F "DM" $DM) }}
	{{ end }}
{{ end }}
{{ end }}
;
{{ end }}
{{ end }}
