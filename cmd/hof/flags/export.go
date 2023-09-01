package flags

import (
	"github.com/spf13/pflag"
)

var _ *pflag.FlagSet

var ExportFlagSet *pflag.FlagSet

type ExportFlagpole struct {
	Expression []string
	List       bool
	Simplify   bool
	Out        string
	Outfile    string
	Escape     bool
	Comments   bool
}

var ExportFlags ExportFlagpole

func SetupExportFlags(fset *pflag.FlagSet, fpole *ExportFlagpole) {
	// flags

	fset.StringArrayVarP(&(fpole.Expression), "expression", "e", nil, "evaluate these expressions only")
	fset.BoolVarP(&(fpole.List), "list", "", false, "concatenate multiple objects into a list")
	fset.BoolVarP(&(fpole.Simplify), "simplify", "", false, "simplify CUE statements where possible")
	fset.StringVarP(&(fpole.Out), "out", "", "", "output data format, when detection does not work")
	fset.StringVarP(&(fpole.Outfile), "outfile", "o", "", "filename or - for stdout with optional file prefix")
	fset.BoolVarP(&(fpole.Escape), "escape", "", false, "use HTLM escaping")
	fset.BoolVarP(&(fpole.Comments), "comments", "C", false, "include comments in output")
}

func init() {
	ExportFlagSet = pflag.NewFlagSet("Export", pflag.ContinueOnError)

	SetupExportFlags(ExportFlagSet, &ExportFlags)

}
