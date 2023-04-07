package cmd

import (
	"os"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/runtime"
)

func diff(R *runtime.Runtime, dflags flags.DatamodelPflagpole) error {

	for _, dm := range R.Datamodels {
		if !dm.HasDiff() {
			continue
		}
		if err := dm.PrintDiff(os.Stdout, dflags); err != nil {
			return err
		}
	}

	return nil
}
