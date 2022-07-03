package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

#GenCommand: schema.#Command & {
	// TBD:   "✓"
	Name:  "gen"
	Usage: "gen [files...]"
	Aliases: ["G"]
	Short: "render directories of code using modular generators"
	Long: GenLongHelp

	Flags: [...schema.#Flag] & [
		{
			Name:    "stats"
			Type:    "bool"
			Default: "false"
			Help:    "Print generator statistics"
			Long:    "stats"
			Short:   "s"
		},
		{
			Name:    "generator"
			Type:    "[]string"
			Default: "nil"
			Help:    "Generators to run, default is all discovered"
			Long:    "generator"
			Short:   "g"
		},
	]
}

#RenderCommand: schema.#Command & {
	// TBD:   "✓"
	Name:  "render"
	Usage: "render [flags] [entrypoints...]"
	Aliases: ["R"]
	Short: "generate arbitrary files from data and CUE entrypoints"
	Long: RenderLongHelp

	Flags: [...schema.#Flag] & [
		{
			Name:    "template"
			Type:    "[]string"
			Default: "nil"
			Help:    "Template mappings to render with data from entrypoint as: <filepath>;<?cuepath>;<?outpath>"
			Long:    "template"
			Short:   "T"
		},
		{
			Name:    "partial"
			Type:    "[]string"
			Default: "nil"
			Help:    "file globs to partial templates to register with the templates"
			Long:    "partial"
			Short:   "P"
		},
		{
			Name:    "diff3"
			Type:    "bool"
			Default: "false"
			Help:    "enable diff3 support, requires the .hof shadow directory"
			Long:    "diff3"
			Short:   "D"
		},
	]
}

GenLongHelp: """
render directories of code with reusable, modular generators

  hof gen app.cue -g frontend -g backend -g migrations

  https://docs.hofstadter.io/first-example/
"""

RenderLongHelp: """
hof render joins CUE with Go's text/template system and diff3
  create on-liners to generate any file from any data
  edit and regenerate those files while keeping changes

# Render a template
hof render data.cue -T template.txt
hof render data.yaml schema.cue -T template.txt > output.txt

# Add partials to the template context
hof render data.cue -T template.txt -P partial.txt

# The template flag as code gen mappings
hof render data.cue ...

  # Generate multiple templates at once
  -T templateA.txt -T templateB.txt

  # Select a sub-input value by CUEpath
  -T 'templateA.txt:foo'
  -T 'templateB.txt:sub.val'

  # Choose a schema with @
  -T 'templateA.txt:foo@#foo'
  -T 'templateB.txt:sub.val@schemas.val'

  # Writing to file with ; (semicolon)
  -T 'templateA.txt;a.txt'
  -T 'templateB.txt:sub.val@schema;b.txt'

  # Templated output path 
  -T 'templateA.txt:;{{ .name | lower }}.txt'

  # Repeated templates are used when
  # 1. the output has a '[]' prefix
  # 2. the input is a list or array
  #   The template will be processed per entry
  #   This also requires using a templated outpath
  -T 'template.txt:items;[]out/{{ .filepath }}.txt'

# Learn about writing templates, with extra functions and helpers
  https://docs.hofstadter.io/code-generation/template-writing/

# Check the tests for complete examples
  https://github.com/hofstadter-io/hof/tree/_dev/test/render

# Compose code gen mappings into reusable modules with
  hof gen app.cue -g frontend -g backend -g migrations
  https://docs.hofstadter.io/first-example/
"""
