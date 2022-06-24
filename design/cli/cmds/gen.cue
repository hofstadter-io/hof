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
	]
}

GenLongHelp: """
render directories of code using modular generators

Doc: https://docs.hofstadter.io/first-example/

hof gen app.cue -g frontend -g backend -g migrations
"""

RenderLongHelp: """
hof render joins CUE with Go's text/template system

# Learn about writing templates, with extra functions and helpers
https://docs.hofstadter.io/code-generation/template-writing/

hof gen is a modular and composable version of this command

# Render a template
hof render data.cue -T template.txt
hof render data.yaml schema.cue -T template.txt > output.txt

# Add partials to the template context
hof render data.cue -T template.txt -P partial.txt

# The template flag
hof render data.cue ...

  # Multiple templates
  -T templateA.txt -T templateB.txt

  # Cuepath to select sub-input values
  -T 'templateA.txt;foo'
  -T 'templateB.txt;sub.val'

  # Writing to file
  -T 'templateA.txt;;a.txt'
  -T 'templateB.txt;sub.val;b.txt'

  # Templated output path 
  -T 'templateA.txt;;{{ .name | ToLower }}.txt'

  # Repeated templates when input is a list
  #   The template will be processed per item
  #   This also requires using a templated outpath
  -T 'template.txt;items;out/{{ .filepath }}.txt'

"""
