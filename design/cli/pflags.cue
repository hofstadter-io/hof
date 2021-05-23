package cli

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

#CliPflags: [...schema.#Flag] & [
	{
		Name:    "package"
		Long:    "package"
		Short:   "p"
		Type:    "string"
		Default: ""
		Help:    "the Cue package context to use during execution"
	},
	{
		Name:    "verbose"
		Long:    "verbose"
		Short:   "v"
		Type:    "int"
		Default: ""
		Help:    "set the verbosity of output"
	},
	{
		Name:    "quiet"
		Long:    "quiet"
		Short:   "q"
		Type:    "bool"
		Default: ""
		Help:    "turn off output and assume defaults at prompts"
	},
]
