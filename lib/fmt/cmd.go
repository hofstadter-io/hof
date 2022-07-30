package fmt

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

func Run(args []string) error {
	fmt.Println("Run:", args)

	return nil
}

func Info(which string) error {
	err := updateFormatterStatus()
	if err != nil {
		return err
	}

	return printAsTable(
		[]string{"Name", "Port", "Image", "Status"},
		func(table *tablewriter.Table) ([][]string, error) {
			var rows = make([][]string, 0, len(fmtrNames))
			// fill with data
			for _,f := range fmtrNames {
				fmtr := formatters[f]

				if which == "" && fmtr.Running == false {
					continue
				}
				
				if fmtr.Container != nil {
					rows = append(rows, []string{
						fmtr.Name,
						fmtr.Port,
						fmtr.Container.Image,
						fmtr.Container.Status,
					})
				} else {
					rows = append(rows, []string{ fmtr.Name, "", "", "", })
				}
			}
			return rows, nil
		},
	)

	return nil
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

type dataPrinter func(table *tablewriter.Table) ([][]string, error)

func printAsTable(headers []string, printer dataPrinter) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)

	defaultTableFormat(table)

	rows, err := printer(table)
	if err != nil {
		return err
	}

	table.AppendBulk(rows)

	// render
	table.Render()

	return nil
}

