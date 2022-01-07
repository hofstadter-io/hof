package datamodel

import (
	"fmt"

	"github.com/olekukonko/tablewriter"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
)

func RunListFromArgs(args []string, flgs flags.DatamodelPflagpole) error {
	// fmt.Println("lib/datamodel.Get", args, flgs)

	// Loadup our Cue files
	dms, err := PrepDatamodels(args, flgs)
	if err != nil {
		return err
	}

	return listDatamodels(dms, flgs)
}

func listDatamodels(dms []*Datamodel, flgs flags.DatamodelPflagpole) error {
	switch flgs.Output {
	case "table":
		return printAsTable(
			[]string{"Name", "Models", "Versions", "Status", "Subsume"},
			func(table *tablewriter.Table) ([][]string, error) {
				var rows = make([][]string, 0, len(dms))
				// fill with data
				for _, dm := range dms {
					lh := len(dm.History.Past)
					nm := fmt.Sprint(len(dm.Models))
					nv := fmt.Sprint(lh)

					sub := "n/a"
					if lh > 0 {
						if dm.Subsume != nil {
							sub = dm.Subsume.Error()
						} else {
							sub = "yes"
						}
					}

					rows = append(rows, []string{dm.Name, nm, nv, dm.Status, sub})
				}
				return rows, nil
			},
		)

	default:
		return fmt.Errorf("Unknown output format %q", flgs.Output)
	}

	return nil
}
