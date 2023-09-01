package flags

import (
	"github.com/spf13/pflag"
)

var _ *pflag.FlagSet

var DefFlagSet *pflag.FlagSet

type DefFlagpole struct {
	Expression    []string
	List          bool
	Simplify      bool
	Out           string
	Outfile       string
	InlineImports bool
	Comments      bool
	Attributes    bool
}

var DefFlags DefFlagpole

func SetupDefFlags(fset *pflag.FlagSet, fpole *DefFlagpole) {
	// flags

	fset.StringArrayVarP(&(fpole.Expression), "expression", "e", nil, "evaluate these expressions only")
	fset.BoolVarP(&(fpole.List), "list", "", false, "concatenate multiple objects into a list")
	fset.BoolVarP(&(fpole.Simplify), "simplify", "", false, "simplify CUE statements where possible")
	fset.StringVarP(&(fpole.Out), "out", "", "", "output data format, when detection does not work")
	fset.StringVarP(&(fpole.Outfile), "outfile", "o", "", "filename or - for stdout with optional file prefix")
	fset.BoolVarP(&(fpole.InlineImports), "inline-imports", "", false, "expand references to non-core imports")
	fset.BoolVarP(&(fpole.Comments), "comments", "C", false, "include comments in output")
	fset.BoolVarP(&(fpole.Attributes), "attributes", "a", false, "diplay field attributes")
}

func init() {
	DefFlagSet = pflag.NewFlagSet("Def", pflag.ContinueOnError)

	SetupDefFlags(DefFlagSet, &DefFlags)

}
