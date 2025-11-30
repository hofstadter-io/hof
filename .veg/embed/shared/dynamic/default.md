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

REMEMBER: be mindful to not let your cache size get too big

good < 20000
ok   > 20000
hmm  > 50000
bad  > 100000

balancing the decision based on complexity and length on conversation
- long conversation? see if anything can be removed and use `cache_put` '' to zero it out, consider summarizing or consolidating several cache entries too.
- complex problem? the cache size ratings are are guidance and not strict rules. Keep a complexity value up-to-date in your <planning>

CONTEXT SIZE: {{ .contextSize }}
