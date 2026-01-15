---
title: "Spacing & Indentation"
weight: 20
---


{{<lead>}}
Formatting, spacing, and indentation are
important for high-quality code.
`hof` assists you in making your
generated code share these qualities.
{{</lead>}}


## Do not template data files

Before we get into hof's helpers for whitespace control,
we want to reiterate that you should not
template data files like Helm has you doing.
This is overly difficult to get right
and manipulation in CUE is much safer.
Hof has dedicated capabilities for generating
data files in several formats along side your code
by leveraging the CUE engine.

See the [generating data files](/code-generation/data-files/) section to see how.


## Code formatting

Hof runs code formatters on the generated files.
This ensures consistent indentation and formatting.

- Data files are handled by the CUE engine
- Go formatting is handled in hof, via Go's parsing libraries
- Other languages are handled by containers

The container based setup is under `hof fmt`,
can be controlled per generator or file,
and is fully customizable and extensible.
The defaults use the most popular tools and rules per-language.

To learn more about the formatters,
see [code-gen/formatting](/code-generation/formatting/).


## Indent Helper

Hof's `indent` helper mirros the behavior
of Helm's while also supporting custom indentation.
Note, the first line is not indented so that
the `{{ ... }}` appears where it should in the template.
This helps you keep your template code visual aligned.

The `indent` helper supports two variations

- the number of spaces to indent
- a string with the indentation to use


{{<codeInner title="indentation variations">}}
# helm style
  {{ indent 2 .content }}

# spaces
    {{ indent "    " .content }}

# tabs
	{{ indent "\t" .content }}

# any string is supported
-> {{ indent "-> " .content }}
{{</codeInner>}}



## Spacing Helpers

If you want to remove whitespace in templates,
use a hyphen inside the curly braces.
This makes it easier to make your template code
understandable while still having correctly formatted output.

{{<codeInner>}}
// removes all whitespace preceeding the {{
{{- .name }}

// removes all whitespace after the }}
{{ .name -}}

// use template comments to ensure exact spacing
// here we are keeping two spaces on both sides of .name
{{- /* ensure before */}}  {{ .name }}  {{/* ensure after */ -}}
{{</codeInner>}}


examples:

- loops
- group spacing
- how to stop

(link to go playground to see how it works)

(lins to examples)



## Partial templates

Can help with both



