package cmd

import (
	// "fmt"

	"github.com/codemodus/kace"
	"github.com/olekukonko/tablewriter"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/datamodel"
	"github.com/hofstadter-io/hof/lib/runtime"
)

func list(R *runtime.Runtime, dflags flags.DatamodelPflagpole) error {
	return printAsTable(
		[]string{"Name", "Type", "Version", "Status", "ID"},
		func(table *tablewriter.Table) ([][]string, error) {
			var rows = make([][]string, 0, len(R.Datamodels))
			// fill with data
			for _, dm := range R.Datamodels {
				id := dm.Hof.Metadata.ID
				if id == "" {
					id = kace.Snake(dm.Hof.Metadata.Name) + " (auto)"
				}

				name := dm.Hof.Metadata.Name
				typ  := datamodel.DatamodelType(dm)
				ver := dm.Hof.Datamodel.Version
				if ver == "" {
					ver = "-"
				}
				status := dm.Status()
				if status == "" {
					status = "-"
				}

				row := []string{name, typ, ver, status, id}
				rows = append(rows, row)
			}
			return rows, nil
		},
	)
}
