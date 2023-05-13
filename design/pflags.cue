package design

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
		Name:    "tags"
		Long:    "tags"
		Short:   "t"
		Type:    "[]string"
		Default: "nil"
		Help:    "@tags() to be injected into CUE code"
	},
	{
		Name:    "verbosity"
		Long:    "verbosity"
		Short:   "v"
		Type:    "int"
		Default: ""
		Help:    "set the verbosity of output"
	},
	{
		Name:    "IgnoreData"
		Long:    "ignore-data"
		Short:   ""
		Type:    "bool"
		Default: ""
		Help:    "ignore all data files unless explicitly supplied"
	},
	{
		Name:    "InjectEnv"
		Long:    "inject-env"
		Short:   ""
		Type:    "bool"
		Default: ""
		Help:    "inject all ENV VARs as default tag vars"
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
