{{ $DM := .DM }}
{{ $S := .Snapshot }}
{{ $DIFF := $S.Lense.CurrDiff }}
/* has snapshot {{ $S.Pos }} - {{ $S.Timestamp }}*/

{{ range $K, $M := $S.Data.Models }}
{{ if and (ne $K "$hof") (ne $K "+")}}
/* == {{ $K }} == */

{{ if $S.Pos }}

{{/* update/create as needed */}}
{{ $UM := lookup $K $DIFF.Models }}
{{ $PTMP := lookup "+" $DIFF.Models }}

{{ if eq (gokind $UM) "map" }}
{{/* is it a new field? */}}
{{ template "update-table.sql" (dict "DM" $DM "Model" $UM "Name" $M.Name) }}

{{ else if eq (gokind $PTMP) "map" }}
{{/* is it a new model? */}}
{{ $CM := lookup $K $PTMP }}
{{ if eq (gokind $CM) "map" }}
{{ template "create-table.sql" (dict "DM" $DM "Model" $M) }}
{{ else }}
/*    create - nothing to do */
{{ end }}

{{ else }}
/*    update - nothing to do */
{{ end }}

{{ else }}
{{/* first generation, so all create */}}
{{ template "create-table.sql" (dict "DM" $DM "Model" $M) }}
{{ end }}{{/* if $S.Pos */}}

{{ end }}
{{ end }}{{/* end range .DM.Models */}}
