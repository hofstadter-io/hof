package flags

import (
	"github.com/spf13/pflag"
)

var _ *pflag.FlagSet

var Datamodel__LogFlagSet *pflag.FlagSet

type Datamodel__LogFlagpole struct {
	ByValue bool
	Details bool
}

var Datamodel__LogFlags Datamodel__LogFlagpole

func SetupDatamodel__LogFlags(fset *pflag.FlagSet, fpole *Datamodel__LogFlagpole) {
	// flags

	fset.BoolVarP(&(fpole.ByValue), "by-value", "", false, "display snapshot log by value")
	fset.BoolVarP(&(fpole.Details), "details", "", false, "print more when displaying the log")
}

func init() {
	Datamodel__LogFlagSet = pflag.NewFlagSet("Datamodel__Log", pflag.ContinueOnError)

	SetupDatamodel__LogFlags(Datamodel__LogFlagSet, &Datamodel__LogFlags)

}
