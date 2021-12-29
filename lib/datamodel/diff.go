package datamodel

import (
	"fmt"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
)

func RunDiffFromArgs(args []string, flgs flags.DatamodelPflagpole) error {
	fmt.Println("lib/datamodel.Diff", args)

	return nil
}
