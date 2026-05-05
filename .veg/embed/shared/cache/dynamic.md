<!-- Your Working Memory. Use cache_put/cache_del to modify. -->
<cache>
{{ range .cache }}
  <entry key="{{.Key}}">
{{.Value}}
  </entry>
{{ end }}
</cache>
