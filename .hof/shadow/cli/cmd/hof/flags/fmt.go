package flags

import (
	"github.com/spf13/pflag"
)

var _ *pflag.FlagSet

var FmtFlagSet *pflag.FlagSet

type FmtFlagpole struct {
	Data bool
}

var FmtFlags FmtFlagpole

func SetupFmtFlags(fset *pflag.FlagSet, fpole *FmtFlagpole) {
	// flags

	fset.BoolVarP(&(fpole.Data), "fmt-data", "", true, "include cue,yaml,json,toml,xml files, set to false to disable")
}

func init() {
	FmtFlagSet = pflag.NewFlagSet("Fmt", pflag.ContinueOnError)

	SetupFmtFlags(FmtFlagSet, &FmtFlags)

}
