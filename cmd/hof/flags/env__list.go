package flags

import (
	"github.com/spf13/pflag"
)

var _ *pflag.FlagSet

var Env__ListFlagSet *pflag.FlagSet

type Env__ListFlagpole struct {
	Kind []string
	Sort []string
}

var Env__ListFlags Env__ListFlagpole

func SetupEnv__ListFlags(fset *pflag.FlagSet, fpole *Env__ListFlagpole) {
	// flags

	fset.StringArrayVarP(&(fpole.Kind), "kind", "k", nil, "kinds to include, defaults to all")
	fset.StringArrayVarP(&(fpole.Sort), "sort", "s", nil, "sort columns, default is the order CUE defines")
}

func init() {
	Env__ListFlagSet = pflag.NewFlagSet("Env__List", pflag.ContinueOnError)

	SetupEnv__ListFlags(Env__ListFlagSet, &Env__ListFlags)

}
