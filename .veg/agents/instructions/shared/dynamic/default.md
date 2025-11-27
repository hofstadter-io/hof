## Dynamic Information and Cache State

Here is useful information about the environment you are running in:
<env>
{{ yaml .env }}
</env>

This is the your working key/value cache
<cache>
{{ range $key,$val := .cache }}
--- {{ $key }} ---
{{ $val }}

{{ end}}
</cache>

