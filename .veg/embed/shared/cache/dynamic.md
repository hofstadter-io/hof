<!-- Your Working Memory. Use cache_put/cache_del to modify. -->
<cache>
{{ range $key,$val := .cache }}
  <entry key="{{$key}}">
{{$val}}
  </entry>
{{ end }}
</cache>
