package partials

FlagVar : RealFlagVar

RealFlagVar : """
{{ define "flag-vars" }}
// FlagVars

{{ end }}


"""

OrigFlagVar : """
{{#with . as |Cmd| }}
{{#if Cmd.pflags}}
var (
{{#each Cmd.pflags}}
{{#if Cmd.parent}}
	{{camelT Cmd.name}}{{camelT name }}PFlag {{> cli/golang/go-type.go type }}
{{else}}
	Root{{camelT name }}PFlag {{> cli/golang/go-type.go type }}
{{/if}}
	{{/each}}
)
{{/if}}

{{#if Cmd.flags}}
var (
{{#each Cmd.flags}}
{{#if Cmd.parent}}
	{{camelT Cmd.name}}{{camelT name }}Flag {{> cli/golang/go-type.go type }}
{{else}}
	Root{{camelT name }}Flag {{> cli/golang/go-type.go type }}
{{/if}}
	{{/each}}
)
{{/if}}
{{/with}}
"""
