package datamodel

import (
	"fmt"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
)

func RunHistoryFromArgs(args []string, flgs flags.DatamodelPflagpole) error {
	// fmt.Println("lib/datamodel.History", args)

	dms, err := LoadDatamodels(args, flgs)
	if err != nil {
		return err
	}

	for _, dm := range dms {
		fmt.Println("---", dm.Name, "---")
		for _, ver := range dm.History.Past {
			fmt.Println(ver.version)
		}
		fmt.Println()
	}
	return nil
}
