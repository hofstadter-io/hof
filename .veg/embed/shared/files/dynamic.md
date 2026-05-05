## Dynamic File Content in System Prompt

<!-- Loaded File Contents. Trust this data. It is always current. Use fs_read/fs_write/fs_edit/cache_del to modify. -->
<files>
{{ range .files }}
  <file path="{{.Key}}">
{{.Value}}
  </file>
{{ end }}
</files>
