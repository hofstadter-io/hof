package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

// ideally this could be a separate flagpole,
// and then embedded into other flagpoles
SharedCueFlags: [...schema.Flag] & [{
	Name:    "expression"
	Long:    "expression"
	Short:   "e"
	Type:    "[]string"
	Default: "nil"
	Help:    "evaluate these expressions only"
}, {
	//  Name:    "extensions"
	//  Long:    "extensions"
	//  Short:   "x"
	//  Type:    "bool"
	//  Default: "false"
	//  Help:    "include hof extensions when evaluating CUE code"
	//}, {
	Name:    "list"
	Long:    "list"
	Type:    "bool"
	Default: "false"
	Help:    "concatenate multiple objects into a list"
}, {
	Name:    "out"
	Long:    "out"
	Type:    "string"
	Default: "\"\""
	Help:    "output data format, when detection does not work"
}, {
	// TODO, consider adding the -T flag from gen here, but only data outputs
	Name:    "outfile"
	Long:    "outfile"
	Short:   "o"
	Type:    "string"
	Default: "\"\""
	Help:    "filename or - for stdout with optional file prefix"
}]

DefCommand: schema.Command & {
	Name:  "def"
	Usage: "def"
	Short: "print consolidated CUE definitions"
	Long:  Short

	Flags: SharedCueFlags + [{
		Name:    "InlineImports"
		Long:    "inline-imports"
		Type:    "bool"
		Default: "false"
		Help:    "expand references to non-core imports"
	}, {
		Name:    "comments"
		Long:    "comments"
		Short:   "C"
		Type:    "bool"
		Default: "false"
		Help:    "include comments in output"
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

	Flags: SharedCueFlags + [{
		Name:    "InlineImports"
		Long:    "inline-imports"
		Type:    "bool"
		Default: "false"
		Help:    "expand references to non-core imports"
	}, {
		Name:    "comments"
		Long:    "comments"
		Short:   "C"
		Type:    "bool"
		Default: "false"
		Help:    "include comments in output"
	}, {
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

	Flags: SharedCueFlags + [{
		Name:    "escape"
		Long:    "escape"
		Type:    "bool"
		Default: "false"
		Help:    "use HTLM escaping"
	}, {
		Name:    "comments"
		Long:    "comments"
		Short:   "C"
		Type:    "bool"
		Default: "false"
		Help:    "include comments in output"
	}]
}

VetCommand: schema.Command & {
	Name:  "vet"
	Usage: "vet"
	Short: "validate data with CUE"
	Long:  Short

	Flags: SharedCueFlags + [{
		Name:    "concrete"
		Long:    "concrete"
		Short:   "c"
		Type:    "bool"
		Default: "false"
		Help:    "require the evaluation to be concrete"
	}]
}
