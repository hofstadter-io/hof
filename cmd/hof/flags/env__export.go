package flags

import (
	"github.com/spf13/pflag"
)

var _ *pflag.FlagSet

var Env__ExportFlagSet *pflag.FlagSet

type Env__ExportFlagpole struct {
	Tag []string
}

var Env__ExportFlags Env__ExportFlagpole

func SetupEnv__ExportFlags(fset *pflag.FlagSet, fpole *Env__ExportFlagpole) {
	// flags

	fset.StringArrayVarP(&(fpole.Tag), "tag", "T", nil, "tags to give to the environment, can be set multiple times")
}

func init() {
	Env__ExportFlagSet = pflag.NewFlagSet("Env__Export", pflag.ContinueOnError)

	SetupEnv__ExportFlags(Env__ExportFlagSet, &Env__ExportFlags)

}
