package datamodel

import (
	"fmt"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/olekukonko/tablewriter"
)

func RunHistoryFromArgs(args []string, flgs flags.DatamodelPflagpole) error {
	// fmt.Println("lib/datamodel.History", args)

	dms, err := PrepDatamodels(args, flgs)
	if err != nil {
		return err
	}

	return histDatamodels(dms, flgs)
}

func histDatamodels(dms []*Datamodel, flgs flags.DatamodelPflagpole) error {
	switch flgs.Output {
	case "table":
		return printAsTable(
			[]string{"Name", "Version", "Timestamp", "Subsume"},
			func(table *tablewriter.Table) ([][]string, error) {
				var rows = make([][]string, 0, len(dms))
				// fill with data
				for _, dm := range dms {
					for _, ver := range dm.History.Past {
						sub := "yes"
						if ver.Subsume != nil {
							sub = dm.Subsume.Error()
						}
						rows = append(rows, []string{dm.Name, ver.Version, ver.Timestamp, sub})
					}
				}
				return rows, nil
			},
		)

	default:
		return fmt.Errorf("Unknown output format %q", flgs.Output)
	}

	return nil
}
