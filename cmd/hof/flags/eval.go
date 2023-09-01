package flags

import (
	"github.com/spf13/pflag"
)

var _ *pflag.FlagSet

var EvalFlagSet *pflag.FlagSet

type EvalFlagpole struct {
	Expression    []string
	List          bool
	Simplify      bool
	Out           string
	Outfile       string
	InlineImports bool
	Comments      bool
	All           bool
	Concrete      bool
	Escape        bool
	Attributes    bool
	Definitions   bool
	Hidden        bool
	Optional      bool
	Resolve       bool
	Defaults      bool
	Final         bool
	Tui           bool
}

var EvalFlags EvalFlagpole

func SetupEvalFlags(fset *pflag.FlagSet, fpole *EvalFlagpole) {
	// flags

	fset.StringArrayVarP(&(fpole.Expression), "expression", "e", nil, "evaluate these expressions only")
	fset.BoolVarP(&(fpole.List), "list", "", false, "concatenate multiple objects into a list")
	fset.BoolVarP(&(fpole.Simplify), "simplify", "", false, "simplify CUE statements where possible")
	fset.StringVarP(&(fpole.Out), "out", "", "", "output data format, when detection does not work")
	fset.StringVarP(&(fpole.Outfile), "outfile", "o", "", "filename or - for stdout with optional file prefix")
	fset.BoolVarP(&(fpole.InlineImports), "inline-imports", "", false, "expand references to non-core imports")
	fset.BoolVarP(&(fpole.Comments), "comments", "C", false, "include comments in output")
	fset.BoolVarP(&(fpole.All), "all", "a", false, "show optional and hidden fields")
	fset.BoolVarP(&(fpole.Concrete), "concrete", "c", false, "require the evaluation to be concrete")
	fset.BoolVarP(&(fpole.Escape), "escape", "", false, "use HTLM escaping")
	fset.BoolVarP(&(fpole.Attributes), "attributes", "A", false, "display field attributes")
	fset.BoolVarP(&(fpole.Definitions), "definitions", "S", true, "display defintions")
	fset.BoolVarP(&(fpole.Hidden), "hidden", "H", false, "display hidden fields")
	fset.BoolVarP(&(fpole.Optional), "optional", "O", false, "display optional fields")
	fset.BoolVarP(&(fpole.Resolve), "resolve", "", false, "resolve references in value")
	fset.BoolVarP(&(fpole.Defaults), "defaults", "", false, "use default values if not set")
	fset.BoolVarP(&(fpole.Final), "final", "", true, "finalize the value")
	fset.BoolVarP(&(fpole.Tui), "tui", "", false, "open hof's TUI and browse your CUE")
}

func init() {
	EvalFlagSet = pflag.NewFlagSet("Eval", pflag.ContinueOnError)

	SetupEvalFlags(EvalFlagSet, &EvalFlags)

}
