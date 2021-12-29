package datamodel

import (
	"fmt"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
)

func RunHistoryFromArgs(args []string, flgs flags.DatamodelPflagpole) error {
	fmt.Println("lib/datamodel.History", args)

	return nil
}
