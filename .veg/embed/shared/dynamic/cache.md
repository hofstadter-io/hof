This is the your key/value cache:
<cache>
{{ range $key,$val := .cache }}
<{{$key}}>
{{$val}}
</{{$key}}>
{{ end}}
</cache>

REMEMBER: be mindful to not let your cache size get too big

good < 50000
ok   > 50000
hmm  > 100000
bad  > 200000

balancing the decision based on complexity and length on conversation
- long conversation? see if anything can be removed and use `cache_put` '' to zero it out, consider summarizing or consolidating several cache entries too.
- complex problem? the cache size ratings are are guidance and not strict rules. Keep a complexity value up-to-date in your <planning>

CONTEXT SIZE: {{ .contextSize }}
