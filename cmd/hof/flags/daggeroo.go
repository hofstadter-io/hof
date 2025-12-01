package flags

import (
	"github.com/spf13/pflag"
)

var _ *pflag.FlagSet

var DaggerooFlagSet *pflag.FlagSet

type DaggerooFlagpole struct {
	Workdir string
	Image   string
}

var DaggerooFlags DaggerooFlagpole

func SetupDaggerooFlags(fset *pflag.FlagSet, fpole *DaggerooFlagpole) {
	// flags

	fset.StringVarP(&(fpole.Workdir), "workdir", "", "/work", "work directory to start in")
	fset.StringVarP(&(fpole.Image), "image", "", "qmcgaw/godevcontainer:debian", "image uri to launch into")
}

func init() {
	DaggerooFlagSet = pflag.NewFlagSet("Daggeroo", pflag.ContinueOnError)

	SetupDaggerooFlags(DaggerooFlagSet, &DaggerooFlags)

}
