package datamodel

import (
	"fmt"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
)

func RunLogFromArgs(args []string, flgs flags.DatamodelPflagpole) error {
	// fmt.Println("lib/datamodel.Info", args, flgs)

	dms, err := PrepDatamodels(args, flgs)
	if err != nil {
		return err
	}

	for _, dm := range dms {
		err = CalcDatamodelStepwiseDiff(dm)
		if err != nil {
			return err
		}
	}

	for _, dm := range dms {
		// shortcut if no history
		if dm.History == nil || len(dm.History.Past) == 0 {
			fmt.Printf("%s: no history")
			fmt.Println(dm.Value)
			continue
		}

		fmt.Println("// current")
		// print current subsume
		if dm.Subsume != nil {
			fmt.Println("// subsume:", dm.Subsume.Error())
		} else {
			fmt.Println("// subsume: yes")
		}

		// print current diff
		if dm.Diff.Exists() {
			fmt.Printf("%s: HEAD: %v\n\n", dm.Name, dm.Diff)
		} else {
			fmt.Printf("%s: HEAD: %v\n\n", dm.Name, "\"no diff\"")
		}

		// print history diff
		past := dm.History.Past
		for i := 0; i < len(past)-1; i++ {
			curr := past[i]
			fmt.Println("//", curr.Version)
			if curr.Subsume != nil {
				fmt.Println("// subsume:", curr.Subsume.Error())
			} else {
				fmt.Println("// subsume: yes")
			}
			fmt.Printf("%s: \"HEAD~%d\": %v\n\n", dm.Name, i+1, curr.Diff)
		}

		// print original (last) value
		last := past[len(past)-1]
		fmt.Println("//", last.Version, "(original)")
		fmt.Println("// subsume: n/a")
		fmt.Printf("%s: \"HEAD~%d\": %v\n", dm.Name, len(past), last.Value)
	}

	return nil
}
