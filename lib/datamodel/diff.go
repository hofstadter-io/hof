package datamodel

import (
	"fmt"

	"github.com/hofstadter-io/cuetils/structural"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
)

func RunDiffFromArgs(args []string, flgs flags.DatamodelPflagpole) error {

	dms, err := LoadDatamodels(args, flgs)
	if err != nil {
		return err
	}

	dms, err = filterDatamodelsByVersion(dms, flgs)
	if err != nil {
		return err
	}

	for _, dm := range dms {
		if len(dm.History.Past) == 0 {
			fmt.Printf("%s: no history\n", dm.Name)
		} else {
			past := dm.History.Past[0]
			if flgs.Since != "" {
				past = dm.History.Past[len(dm.History.Past)-1]
			}

			fmt.Printf("// %s -> %s\n%s: ", dm.History.Past[0].Version, dm.Version, dm.Name)
			diff, err := structural.DiffValue(past.Value, dm.Value, nil)
			if err != nil {
				return err
			}
			if !diff.Exists() {
				fmt.Println("_|_")
			} else {
				fmt.Println(diff)
			}
		}
	}

	return nil
}

func CalcDatamodelStepwiseDiff(dm *Datamodel) error {
	if dm.History == nil || len(dm.History.Past) == 0 {
		return nil
	}
	past := dm.History.Past

	// loop back through time (checkpoints)
	curr := dm
	for i := 0; i < len(past); i++ {
		// get prev to compare against
		prev := past[i]

		// calculate what needs to be done to prev to get to curr
		diff, err := structural.DiffValue(prev.Value, curr.Value, nil)
		if err != nil {
			return err
		}

		// set changes need to arrive at curr
		curr.Diff = diff
		// update before relooping
		curr = prev
	}

	return nil
}
