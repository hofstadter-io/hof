These are the current files:
<files>
{{ range $path,$content := .files }}
<file path="{{$path}}">
{{$content}}
</file>
{{ end}}
</files>