package datamodel

import (
	"fmt"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
)

func RunDiffFromArgs(args []string, flgs flags.DatamodelPflagpole) error {
	// fmt.Println("lib/datamodel.Diff", args)

	dms, err := LoadDatamodels(args, flgs)
	if err != nil {
		return err
	}

	for _, dm := range dms {
		fmt.Println("---", dm.Name, "---")
		if len(dm.History.Past) == 0 {
			fmt.Println("no history to diff against")
		} else {
			fmt.Println("compare curr to", dm.History.Past[0])
		}
	}

	return nil
}
