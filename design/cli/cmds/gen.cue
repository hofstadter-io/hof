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
			Help:    "print generator statistics"
			Long:    "stats"
			Short:   "s"
		},
		{
			Name:    "generator"
			Type:    "[]string"
			Default: "nil"
			Help:    "generator tags to run, default is all"
			Long:    "generator"
			Short:   "G"
		},
		{
			Name:    "template"
			Type:    "[]string"
			Default: "nil"
			Help:    "template mappings to render as '<filepath>;<?cuepath>;<?outpath>'"
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
		{
			Name:    "watch"
			Type:    "bool"
			Default: "false"
			Help:    "run in watch mode, regenerating when files change"
			Long:    "watch"
			Short:   "w"
		},
		{
			Name:    "WatchGlobs"
			Type:    "[]string"
			Default: "nil"
			Help:    "filepath globs to watch for changes and regen"
			Long:    "watch-globs"
			Short:   "W"
		},
		{
			Name:    "WatchXcue"
			Type:    "[]string"
			Default: "nil"
			Help:    "like watch, but skips CUE reload, useful when working on templates, can be used with watch"
			Long:    "watch-xcue"
			Short:   "X"
		},
		{
			Name:    "CreateModule"
			Type:    "string"
			Default: ""
			Help:    "output path to write a generator module for the given flags"
			Long:    "create-module"
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

# You can mix adhof with generators by using
# both the -G and -T/-P flags
"""
