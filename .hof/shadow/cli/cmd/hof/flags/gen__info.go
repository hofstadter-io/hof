package flags

import (
	"github.com/spf13/pflag"
)

var _ *pflag.FlagSet

var Gen__InfoFlagSet *pflag.FlagSet

type Gen__InfoFlagpole struct {
	Expression []string
}

var Gen__InfoFlags Gen__InfoFlagpole

func SetupGen__InfoFlags(fset *pflag.FlagSet, fpole *Gen__InfoFlagpole) {
	// flags

	fset.StringArrayVarP(&(fpole.Expression), "expr", "e", nil, "CUE paths to select outputs, depending on the command")
}

func init() {
	Gen__InfoFlagSet = pflag.NewFlagSet("Gen__Info", pflag.ContinueOnError)

	SetupGen__InfoFlags(Gen__InfoFlagSet, &Gen__InfoFlags)

}
