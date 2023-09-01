package flags

import (
	"github.com/spf13/pflag"
)

var _ *pflag.FlagSet

var Datamodel__CheckpointFlagSet *pflag.FlagSet

type Datamodel__CheckpointFlagpole struct {
	Message string
}

var Datamodel__CheckpointFlags Datamodel__CheckpointFlagpole

func SetupDatamodel__CheckpointFlags(fset *pflag.FlagSet, fpole *Datamodel__CheckpointFlagpole) {
	// flags

	fset.StringVarP(&(fpole.Message), "message", "m", "", "message describing the checkpoint")
}

func init() {
	Datamodel__CheckpointFlagSet = pflag.NewFlagSet("Datamodel__Checkpoint", pflag.ContinueOnError)

	SetupDatamodel__CheckpointFlags(Datamodel__CheckpointFlagSet, &Datamodel__CheckpointFlags)

}
