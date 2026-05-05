# == CURRENT SYSTEM STATE ==

CONTEXT SIZE: {{ .contextSize }}

<!-- Your Working Memory. Use cache_put/cache_del to modify. -->
<cache>
{{ range $key,$val := .cache }}
  <entry key="{{$key}}">
{{$val}}
  </entry>
{{ end }}
</cache>

<!-- Loaded File Contents. Trust this data. -->
<files>
{{ range $path,$content := .files }}
  <file path="{{$path}}">
{{$content}}
  </file>
{{ end }}
</files>

<!-- Current Strategic Plan (Derived from cache key 'planning') -->
<planning>
{{ .planning }} 
</planning>

<!-- Environment Info -->
<env>
{{ yaml .env }}
</env>
