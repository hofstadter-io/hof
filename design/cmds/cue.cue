package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

CueFlags: [...schema.Flag] & [{
	Name:    "all"
	Long:    "all"
	Short:   "a"
	Type:    "bool"
	Default: "false"
	Help:    "show optional and hidden fields"
}, {
	Name:    "path"
	Long:    "path"
	Short:   "l"
	Type:    "[]string"
	Default: "\"feedback\""
	Help:    "labels,comma,separated"
}]

DefCommand: schema.Command & {
	Name:  "def"
	Usage: "def"
	Short: "print consolidated CUE definitions"
	Long:  Short

	Flags: [{
		Name:    "issue"
		Long:    "issue"
		Short:   "i"
		Type:    "bool"
		Default: "false"
		Help:    "create an issue (discussion is default)"
	}]
}

EvalCommand: schema.Command & {
	Name:  "eval"
	Usage: "eval"
	Short: "evaluate and print CUE configuration"
	Long:  Short

	Flags: [{
		Name:    "issue"
		Long:    "issue"
		Short:   "i"
		Type:    "bool"
		Default: "false"
		Help:    "create an issue (discussion is default)"
	}]
}

ExportCommand: schema.Command & {
	Name:  "export"
	Usage: "export"
	Short: "output data in a standard format"
	Long:  Short

	Flags: [{
		Name:    "issue"
		Long:    "issue"
		Short:   "i"
		Type:    "bool"
		Default: "false"
		Help:    "create an issue (discussion is default)"
	}]
}

VetCommand: schema.Command & {
	Name:  "vet"
	Usage: "vet"
	Short: "validate data with CUE"
	Long:  Short

	Flags: [{
		Name:    "issue"
		Long:    "issue"
		Short:   "i"
		Type:    "bool"
		Default: "false"
		Help:    "create an issue (discussion is default)"
	}]
}

// Eval, Export, Vet
