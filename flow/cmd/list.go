package cmd

import (
	"github.com/codemodus/kace"
	"github.com/olekukonko/tablewriter"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/yagu"
)

func List(args []string, rflags flags.RootPflagpole, cflags flags.FlowPflagpole) error {
	R, err := prepRuntime(args, rflags, cflags)
	if err != nil {
		return err
	}

	return yagu.PrintAsTable(
		[]string{"Name", "Path", "ID"},
		func(table *tablewriter.Table) ([][]string, error) {
			var rows = make([][]string, 0, len(R.Workflows))
			// fill with data
			for _, wf := range R.Workflows {
				id := wf.Hof.Metadata.ID
				if id == "" {
					id = kace.Snake(wf.Hof.Metadata.Name) + " (auto)"
				}

				name := wf.Hof.Flow.Name
				if name == "" {
					name = "(anon)"
				}
				path := wf.Hof.Path

				row := []string{name, path, id}
				rows = append(rows, row)
			}
			return rows, nil
		},
	)
}
