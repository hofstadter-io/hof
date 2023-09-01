package flags

import (
	"github.com/spf13/pflag"
)

var _ *pflag.FlagSet

var Mod__GetFlagSet *pflag.FlagSet

type Mod__GetFlagpole struct {
	Prerelease bool
}

var Mod__GetFlags Mod__GetFlagpole

func SetupMod__GetFlags(fset *pflag.FlagSet, fpole *Mod__GetFlagpole) {
	// flags

	fset.BoolVarP(&(fpole.Prerelease), "prerelease", "P", false, "include prerelease version when using @latest")
}

func init() {
	Mod__GetFlagSet = pflag.NewFlagSet("Mod__Get", pflag.ContinueOnError)

	SetupMod__GetFlags(Mod__GetFlagSet, &Mod__GetFlags)

}
