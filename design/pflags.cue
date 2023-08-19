package design

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

CliPflags: [...schema.Flag] & [ {
	Name:    "package"
	Long:    "package"
	Short:   "p"
	Type:    "string"
	Default: ""
	Help:    "the Cue package context to use during execution"
}, {
	Name:    "tags"
	Long:    "tags"
	Short:   "t"
	Type:    "[]string"
	Default: "nil"
	Help:    "@tags() to be injected into CUE code"
}, {
	Name:    "path"
	Long:    "path"
	Short:   "l"
	Type:    "[]string"
	Default: "nil"
	Help:    "CUE expression for single path component when placing data files"
}, {
	Name:    "verbosity"
	Long:    "verbosity"
	Short:   "v"
	Type:    "int"
	Default: ""
	Help:    "set the verbosity of output"
}, {
	Name:    "IncludeData"
	Long:    "include-data"
	Short:   "D"
	Type:    "bool"
	Default: ""
	Help:    "auto include all data files found with cue files"
}, {
	Name:    "InjectEnv"
	Long:    "inject-env"
	Short:   "V"
	Type:    "bool"
	Default: ""
	Help:    "inject all ENV VARs as default tag vars"
}, {
	Name:    "quiet"
	Long:    "quiet"
	Short:   "q"
	Type:    "bool"
	Default: ""
	Help:    "turn off output and assume defaults at prompts"
}, {
	Name:    "AllErrors"
	Long:    "all-errors"
	Short:   "E"
	Type:    "bool"
	Default: ""
	Help:    "print all available errors"
}, {
	Name:    "IngoreErrors"
	Long:    "ignore-errors"
	Short:   "i"
	Type:    "bool"
	Default: ""
	Help:    "turn off output and assume defaults at prompts"
}]
