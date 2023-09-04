package flags

import (
	"github.com/spf13/pflag"
)

var _ *pflag.FlagSet

var RootFlagSet *pflag.FlagSet

type RootPflagpole struct {
	Package      string
	InputData    []string
	StdinEmpty   bool
	Tags         []string
	Path         []string
	Schema       []string
	IncludeData  bool
	WithContext  bool
	InjectEnv    bool
	AllErrors    bool
	IngoreErrors bool
	Stats        bool
	Quiet        bool
	Verbosity    int
}

func SetupRootPflags(fset *pflag.FlagSet, fpole *RootPflagpole) {
	// pflags

	fset.StringVarP(&(fpole.Package), "package", "p", "", "the Cue package context to use during execution")
	fset.StringArrayVarP(&(fpole.InputData), "input", "I", nil, "extra data to unify into the root value")
	fset.BoolVarP(&(fpole.StdinEmpty), "stdin-empty", "0", false, "A flag that ensure stdin is zero and does not block")
	fset.StringArrayVarP(&(fpole.Tags), "tags", "t", nil, "@tags() to be injected into CUE code")
	fset.StringArrayVarP(&(fpole.Path), "path", "l", nil, "CUE expression for single path component when placing data files")
	fset.StringArrayVarP(&(fpole.Schema), "schema", "d", nil, "expression to select schema to apply to data files")
	fset.BoolVarP(&(fpole.IncludeData), "include-data", "D", false, "auto include all data files found with cue files")
	fset.BoolVarP(&(fpole.WithContext), "with-context", "", false, "add extra context for data files, usable in the -l/path flag")
	fset.BoolVarP(&(fpole.InjectEnv), "inject-env", "V", false, "inject all ENV VARs as default tag vars")
	fset.BoolVarP(&(fpole.AllErrors), "all-errors", "E", false, "print all available errors")
	fset.BoolVarP(&(fpole.IngoreErrors), "ignore-errors", "i", false, "turn off output and assume defaults at prompts")
	fset.BoolVarP(&(fpole.Stats), "stats", "", false, "print generator statistics")
	fset.BoolVarP(&(fpole.Quiet), "quiet", "q", false, "turn off output and assume defaults at prompts")
	fset.IntVarP(&(fpole.Verbosity), "verbosity", "v", 0, "set the verbosity of output")
}

var RootPflags RootPflagpole

func init() {
	RootFlagSet = pflag.NewFlagSet("Root", pflag.ContinueOnError)

	SetupRootPflags(RootFlagSet, &RootPflags)

}
