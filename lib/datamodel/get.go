package datamodel

import (
	"fmt"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/cuetils"
)

func RunGetFromArgs(args []string, cmdpflags flags.DatamodelPflagpole) error {
	fmt.Println("lib/datamodel.Get", args, cmdpflags)

	cueFiles := args

	// Loadup our Cue files
	crt, err := cuetils.CueRuntimeFromEntrypointsAndFlags(cueFiles)
	if err != nil {
		return err
	}

	// TODO: find values from flags / attributes
	val := crt.CueValue

	syn, err := cuetils.PrintCueValue(val)
	if err != nil {
		return err
	}

	fmt.Println(syn)

	return nil
}
