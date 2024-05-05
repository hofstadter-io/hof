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


{{/*
need to loop over models here 
where to interleave?
(need to know if at snapshot 0 or not for a particular model)
( while looping over root datamodel snapshots, which has every checkpoint )
( models will not typically all start or change at the same time,
  but we have to apply them in the overall timestamp order for the database )
*/}}

{{ if .DM.History }}
{{ range $K, $M := .DM.Models }}
{{ if ne $K "$hof" }}


CREATE TABLE {{ snake $M.Name }} (

{{ range $K, $F := $M.Fields }}
{{ if ne $K "$hof" }}{{ template "create-field.sql" (dict "Field" $F "DM" $DM) }}{{ end }}
{{ end }}

);
{{ end }}
{{ end }}


{{ if $DM.CurrDiff }}
/* has currdiff */
{{ template "update-table.sql" (dict "DM" $DM "SNAP" $SNAP) }}
{{ else }}
	{{ /* nothing to do */ }}
{{ end }}
{{ else }}
/* first datamodel, no checkpoints */
{{ template "create-table.sql" (dict "DM" $DM "SNAP" $SNAP) }}
{{ end }}

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
