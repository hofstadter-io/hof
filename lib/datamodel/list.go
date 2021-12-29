package datamodel

import (
	"fmt"
	"os"

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

	return printDatamodelList(dms, "table")
}

func printDatamodelList(dms []*Datamodel, format string) error {
	switch format {
	case "table":
		printDatamodelListTable(dms)

	default:
		return fmt.Errorf("Unknown format %q", format)
	}

	return nil
}

func printDatamodelListTable(dms []*Datamodel) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Models", "Versions", "Status"})
	defaultTableFormat(table)

	// fill with data
	for _, dm := range dms {
		nm := fmt.Sprint(len(dm.Models))
		nv := fmt.Sprint(len(dm.History.Past))
		table.Append([]string{dm.Name, nm, nv, dm.status})
	}

	// render
	table.Render()
}

func defaultTableFormat(table *tablewriter.Table) {
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("  ") // pad with tabs
	table.SetNoWhiteSpace(true)
}
