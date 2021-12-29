package datamodel

import (
	"fmt"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
)

func RunInfoFromArgs(args []string, flgs flags.DatamodelPflagpole) error {
	fmt.Println("lib/datamodel.Status", args)

	return nil
}
