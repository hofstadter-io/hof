package cmd

import (
	"fmt"

	"github.com/codemodus/kace"
	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/yagu"
	"github.com/olekukonko/tablewriter"
)

func Pull(args []string, rflags flags.RootPflagpole) error {
	R, err := prepRuntime(args, rflags)
	if err != nil {
		return err
	}

	return yagu.PrintAsTable(
		[]string{"Name", "Path", "ID", "Extra"},
		func(table *tablewriter.Table) ([][]string, error) {
			var rows = make([][]string, 0, len(R.Envs))
			// fill with data
			for _, e := range R.Envs {
				id := e.Hof.Metadata.ID
				if id == "" {
					id = kace.Snake(e.Hof.Metadata.Name) + " (auto)"
				}

				name := e.Hof.Env.Name
				if name == "" {
					name = "(anon)"
				}
				path := e.Hof.Path

				row := []string{name, path, id, fmt.Sprint(e.Hof.Env.Extra)}
				rows = append(rows, row)
			}
			return rows, nil
		},
	)
}
