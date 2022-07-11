package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

#GenCommand: schema.#Command & {
	// TBD:   "✓"
	Name:  "gen"
	Usage: "gen [files...]"
	Aliases: ["G"]
	Short: "create with modular code gen...  CUE & data + templates = _"
	Long: GenLongHelp

	Flags: [...schema.#Flag] & [
		{
			Name:    "list"
			Type:    "bool"
			Default: "false"
			Help:    "list available generators"
			Long:    "list"
			Short:   "l"
		},
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
			Help:    "generator tags to run, default is all, or none if -T is used"
			Long:    "generator"
			Short:   "G"
		},
		{
			Name:    "template"
			Type:    "[]string"
			Default: "nil"
			Help:    "template mapping to render, see help for format"
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
			Help:    "enable diff3 support for custom code"
			Long:    "diff3"
			Short:   "D"
		},
		{
			Name:    "watch"
			Type:    "bool"
			Default: "false"
			Help:    "run in watch mode, regenerating when files change, implied by -W/X"
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
			Help:    "like watch, but skips CUE reload, good for template authoring"
			Long:    "watch-xcue"
			Short:   "X"
		},
		{
			Name:    "AsModule"
			Type:    "string"
			Default: ""
			Help:    "<name> for the generator module made from the given flags"
			Long:    "as-module"
		},
		{
			Name:    "InitModule"
			Type:    "string"
			Default: ""
			Help:    "<name> to bootstrap a new genarator module"
			Long:    "init"
		},
		{
			Name:    "Outdir"
			Type:    "string"
			Default: ""
			Help:    "base directory to write all output u"
			Long:    "outdir"
			Short:   "O"
		},
	]
}

GenLongHelp: """
hof unifies CUE with Go's text/template system and diff3
  create on-liners to generate any file from any 
  build reusable and modular generators
  edit and regenerate those files while keeping changes

# Render a template
  hof gen input.cue -T template.txt
  hof gen input.yaml schema.cue -T template.txt > output.txt

# Add partials to the template context
  hof gen input.cue -T template.txt -P partial.txt

# The template flag as code gen mappings
  hof gen input.cue -T ... -T ...

  # Generate multiple templates at once
  -T templateA.txt -T templateB.txt

  # Select a sub-input value by CUEpath
  -T templateA.txt:foo
  -T templateB.txt:sub.val

  # Choose a schema with @
  -T templateA.txt:foo@Foo
  -T templateB.txt:sub.val@schemas.val

  # Writing to file with = (semicolon)
  -T templateA.txt=a.txt
  -T templateB.txt:sub.val@schema=b.txt

  # Templated output path, braces need quotes
  -T templateA.txt:='{{ .name | lower }}.txt'

  # Repeated templates are used when
  # 1. the output has a '[]' prefix
  # 2. the input is a list or array
  #   The template will be processed per entry
  #   This also requires using a templated outpath
  -T template.txt:items='[]out/{{ .filepath }}.txt'

  # Output everything to a directory (out name is the same)
  -O out -T types.go -T handlers.go

  # Watch files and directories, doing full or Xcue-less reloads
  -W *.cue -X *.go -O out -T types.go -T handlers.go

# Turn any hof gen flags into a reusable generator module
  hof gen [entrypoints] flags... --as-module [name]
  hof gen [entrypoints] -G [name]

# Bootstrap a new generator module
  hof gen --init github.com/hofstadter-io/demos

# Learn about writing templates, with extra functions and helpers
  https://docs.hofstadter.io/code-generation/template-writing/

# Check the tests for complete examples
  https://github.com/hofstadter-io/hof/tree/_dev/test/render

# Compose code gen mappings into reusable modules with
  hof gen app.cue -G frontend -G backend -G migrations -T ...
  https://docs.hofstadter.io/first-example/
"""
