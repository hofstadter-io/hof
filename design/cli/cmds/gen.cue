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
	Long: """
  \(Short)

  Doc: https://docs.hofstadter.io/first-example/

  hof gen -g frontend -g backend design.cue
  """

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
  Long: """
  \(Short)

  hof render -t template.go data.cue > file.go
  """

	Flags: [...schema.#Flag] & [
		{
			Name:    "template"
			Type:    "[]string"
			Default: "nil"
			Help:    "Template mappings to render with data from entrypoint as: filepath|cuepath|outpath"
			Long:    "template"
			Short:   "t"
		},
	]
}


