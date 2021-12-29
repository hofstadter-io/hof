package datamodel

import (
	"fmt"
	"regexp"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/olekukonko/tablewriter"
)

func RunInfoFromArgs(args []string, flgs flags.DatamodelPflagpole) error {
	// fmt.Println("lib/datamodel.Info", args, flgs)

	dms, err := LoadDatamodels(args, flgs)
	if err != nil {
		return err
	}

	dms, err = filterDatamodelsByVersion(dms, flgs)
	if err != nil {
		return err
	}

	return infoDatamodels(dms, flgs)
}

func infoDatamodels(dms []*Datamodel, flgs flags.DatamodelPflagpole) error {
	switch flgs.Output {
	case "cue":
		return infoDatamodelsCue(dms, flgs)
	case "table":
		return infoDatamodelsTable(dms, flgs)
	default:
		return fmt.Errorf("Unknown format %q", flgs.Output)
	}

	return nil
}

func infoDatamodelsTable(dms []*Datamodel, flgs flags.DatamodelPflagpole) error {
	return printAsTable(
		[]string{"DM Name", "Models", "Fields", "Type", "Status"},
		func(table *tablewriter.Table) ([][]string, error) {
			var rows = make([][]string, 0, len(dms))
			// fill with data
			for _, dm := range dms {
				dmn := dm.Name
				for _, m := range dm.Models {
					nf := fmt.Sprint(len(m.Fields))
					rows = append(rows, []string{dmn, m.Name, nf, "model", m.status})
					if len(flgs.Models) > 0 {
						for _, f := range m.Fields {
							rows = append(rows, []string{"", "", f.Name, f.Type, ""})
						}
					}
					// only print once
					if dmn != "" {
						dmn = ""
					}
				}
			}
			return rows, nil
		},
	)
}

func infoDatamodelsCue(dms []*Datamodel, flgs flags.DatamodelPflagpole) error {

	for _, dm := range dms {
		// print whole datamodels
		if len(flgs.Models) == 0 {
			fmt.Printf("%s: %v\n", dm.Name, dm.value)
			continue
		}

		// print whole models
		if len(flgs.Models) > 0 {
			for _, m := range dm.Models {
				match := false
				for _, regx := range flgs.Models {
					if match, _ = regexp.MatchString(regx, m.Name); match {
						fmt.Printf("%s: %s: %v\n", dm.Name, m.Name, m.value)
					}
				}
			}
			continue
		}

	}

	return nil
}
