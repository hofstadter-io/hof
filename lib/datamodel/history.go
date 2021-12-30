package datamodel

import (
	"fmt"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
)

func RunHistoryFromArgs(args []string, flgs flags.DatamodelPflagpole) error {
	// fmt.Println("lib/datamodel.History", args)

	dms, err := PrepDatamodels(args, flgs)
	if err != nil {
		return err
	}

	for _, dm := range dms {
		fmt.Println("---", dm.Name, "---")
		for _, ver := range dm.History.Past {
			fmt.Println(ver.Version)
		}

		if len(dms) > 1 {
			fmt.Println()
		}
	}
	return nil
}
