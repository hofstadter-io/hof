<!-- Loaded File Contents. Trust this data. Use fs_read/fs_write/fs_edit/cache_del to modify. -->
<files>
{{ range $path,$content := .files }}
  <file path="{{$path}}">
{{$content}}
  </file>
{{ end }}
</files>
