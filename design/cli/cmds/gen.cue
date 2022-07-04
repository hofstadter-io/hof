package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

#GenCommand: schema.#Command & {
	// TBD:   "✓"
	Name:  "gen"
	Usage: "gen [files...]"
	Aliases: ["G"]
	Short: "create arbitrary files from data with templates and generators"
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
			Short:   "G"
		},
		{
			Name:    "template"
			Type:    "[]string"
			Default: "nil"
			Help:    "Template mappings to render as '<filepath>;<?cuepath>;<?outpath>'"
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
			Help:    "enable diff3 support for adhoc render, generators are configured in code"
			Long:    "diff3"
			Short:   "D"
		},
	]
}

GenLongHelp: """
hof gen joins CUE with Go's text/template system and diff3
  create on-liners to generate any file from any data
  build reusable and modular generators
  edit and regenerate those files while keeping changes

If no generator is specified, hof gen runs in adhoc mode.

# Render a template
hof gen data.cue -T template.txt
hof gen data.yaml schema.cue -T template.txt > output.txt

# Add partials to the template context
hof gen data.cue -T template.txt -P partial.txt

# The template flag as code gen mappings
hof gen data.cue ...

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
  hof gen app.cue -G frontend -G backend -G migrations
  https://docs.hofstadter.io/first-example/

# You can extend or override a generator by using
# both the -G and -T/-P flags
"""
