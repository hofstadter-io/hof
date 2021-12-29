package datamodel

import (
	"fmt"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
)

func RunCheckpointFromArgs(args []string, flgs flags.DatamodelPflagpole) error {
	// fmt.Println("lib/datamodel.Checkpoint", args)

	dms, err := LoadDatamodels(args, flgs)
	if err != nil {
		return err
	}

	for _, dm := range dms {
		if dm.status == "dirty" {
			fmt.Println("checkpoint:", dm.Name)
		}
	}

	return nil
}
