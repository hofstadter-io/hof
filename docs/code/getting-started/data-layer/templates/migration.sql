{{ $SNAP := .Snapshot }}
{{ $DM := .DM }}

{{/* 5 cases here */}}
{{ if $SNAP }}
/* has snapshot {{ $SNAP.Pos }} - {{ $SNAP.Timestamp }}*/
{{ if $SNAP.Pos }}
{{ template "update-table.sql" (dict "DM" $DM "SNAP" $SNAP) }}
{{ else }}
/* first datamodel, first snapshot */
{{ template "create-table.sql" (dict "DM" $DM "SNAP" $SNAP) }}
{{ end }}
{{ else if $DM.History }}
{{ if $DM.CurrDiff }}
/* has currdiff */
{{ template "update-table.sql" (dict "DM" $DM "SNAP" $SNAP) }}
{{ else }}
/* nothing to do */
{{ end }}
{{ else }}
/* first datamodel, no checkpoints */
{{ template "create-table.sql" (dict "DM" $DM "SNAP" $SNAP) }}
{{ end }}

