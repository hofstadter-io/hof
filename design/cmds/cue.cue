package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

SharedCueFlags: [...schema.Flag] & [{
	Name:    "expression"
	Long:    "expression"
	Short:   "e"
	Type:    "[]string"
	Default: "nil"
	Help:    "evaluate these expressions only"
}, {
	Name:    "path"
	Long:    "path"
	Short:   "l"
	Type:    "[]string"
	Default: "nil"
	Help:    "CUE expression for single path componenti when placing data files"
}]

DefCommand: schema.Command & {
	Name:  "def"
	Usage: "def"
	Short: "print consolidated CUE definitions"
	Long:  Short

	Flags: [{
		Name:    "InlineImports"
		Long:    "inline-imports"
		Type:    "bool"
		Default: "false"
		Help:    "expand references to non-core imports"
	}, {
		Name:    "attributes"
		Long:    "attributes"
		Short:   "a"
		Type:    "bool"
		Default: "false"
		Help:    "diplay field attributes"
	}]
}

EvalCommand: schema.Command & {
	Name:  "eval"
	Usage: "eval"
	Short: "evaluate and print CUE configuration"
	Long:  Short

	Flags: [{
		Name:    "all"
		Long:    "all"
		Short:   "a"
		Type:    "bool"
		Default: "false"
		Help:    "show optional and hidden fields"
	}, {
		Name:    "concrete"
		Long:    "concrete"
		Short:   "c"
		Type:    "bool"
		Default: "false"
		Help:    "require the evaluation to be concrete"
	}, {
		Name:    "attributes"
		Long:    "attributes"
		Short:   "A"
		Type:    "bool"
		Default: "false"
		Help:    "diplay field attributes"
	}, {
		Name:    "hidden"
		Long:    "hidden"
		Short:   "H"
		Type:    "bool"
		Default: "false"
		Help:    "diplay hidden fields"
	}, {
		Name:    "optional"
		Long:    "optional"
		Short:   "O"
		Type:    "bool"
		Default: "false"
		Help:    "diplay optional fields"
	}]
}

ExportCommand: schema.Command & {
	Name:  "export"
	Usage: "export"
	Short: "output data in a standard format"
	Long:  Short

	Flags: [{
		Name:    "escape"
		Long:    "espace"
		Type:    "bool"
		Default: "false"
		Help:    "use HTLM escaping"
	}, {
		Name:    "out"
		Long:    "out"
		Type:    "string"
		Default: "\"\""
		Help:    "output data format, when detection does not work"
	}]
}

VetCommand: schema.Command & {
	Name:  "vet"
	Usage: "vet"
	Short: "validate data with CUE"
	Long:  Short

	Flags: [{
		Name:    "concrete"
		Long:    "concrete"
		Short:   "c"
		Type:    "bool"
		Default: "false"
		Help:    "require the evaluation to be concrete"
	}]
}

// Eval, Export, Vet
