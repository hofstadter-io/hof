package partials

GoType : """
{{#with . as |T| }}
{{#if (contains T ":")}}
	{{#if (hasprefix T "array")}}
[]{{trimprefix T "array:"}}
	{{/if}}
	{{#if (hasprefix T "map")}}
map[string]{{trimprefix T "map:"}}
	{{/if}}
{{else}}
{{T}}
{{/if}}
{{/with}}
"""
