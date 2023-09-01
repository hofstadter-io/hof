package flags

import (
	"github.com/spf13/pflag"
)

var _ *pflag.FlagSet

var DatamodelFlagSet *pflag.FlagSet

type DatamodelPflagpole struct {
	Datamodels []string
	Expression []string
}

func SetupDatamodelPflags(fset *pflag.FlagSet, fpole *DatamodelPflagpole) {
	// pflags

	fset.StringArrayVarP(&(fpole.Datamodels), "model", "M", nil, "specify one or more data models to operate on")
	fset.StringArrayVarP(&(fpole.Expression), "expr", "e", nil, "CUE paths to select outputs, depending on the command")
}

var DatamodelPflags DatamodelPflagpole

func init() {
	DatamodelFlagSet = pflag.NewFlagSet("Datamodel", pflag.ContinueOnError)

	SetupDatamodelPflags(DatamodelFlagSet, &DatamodelPflags)

}
