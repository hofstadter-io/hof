{{ if .Snapshot }}
{{ template "snapshot-table.sql" (dict "DM" .DM "Snapshot" .Snapshot) }}
{{ else }}
{{ template "latest-table.sql" (dict "DM" .DM) }}
{{ end }}

