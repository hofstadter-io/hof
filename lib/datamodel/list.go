package datamodel

import (
	"fmt"

	"github.com/olekukonko/tablewriter"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
)

func RunListFromArgs(args []string, flgs flags.DatamodelPflagpole) error {
	// fmt.Println("lib/datamodel.Get", args, flgs)

	// Loadup our Cue files
	dms, err := LoadDatamodels(args, flgs)
	if err != nil {
		return err
	}

	dms, err = filterDatamodelsByVersion(dms, flgs)
	if err != nil {
		return err
	}

	return listDatamodels(dms, flgs)
}

func listDatamodels(dms []*Datamodel, flgs flags.DatamodelPflagpole) error {
	switch flgs.Output {
	case "table":
		return printAsTable(
			[]string{"Name", "Models", "Versions", "Status"},
			func(table *tablewriter.Table) ([][]string, error) {
				var rows = make([][]string, 0, len(dms))
				// fill with data
				for _, dm := range dms {
					nm := fmt.Sprint(len(dm.Models))
					nv := fmt.Sprint(len(dm.History.Past))
					rows = append(rows, []string{dm.Name, nm, nv, dm.status})
				}
				return rows, nil
			},
		)

	default:
		return fmt.Errorf("Unknown output format %q", flgs.Output)
	}

	return nil
}
