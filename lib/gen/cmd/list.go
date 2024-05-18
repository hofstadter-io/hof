package cmd

import (
	"fmt"

	"github.com/codemodus/kace"
	"github.com/olekukonko/tablewriter"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/yagu"
)

func List(args []string, rflags flags.RootPflagpole, gflags flags.GenFlagpole) error {
	R, err := prepRuntime(args, rflags, gflags)
	if err != nil {
		return err
	}

	return yagu.PrintAsTable(
		[]string{"Name", "Path", "ID", "Creator"},
		func(table *tablewriter.Table) ([][]string, error) {
			var rows = make([][]string, 0, len(R.Generators))
			// fill with data
			for _, gen := range R.Generators {
				id := gen.Hof.Metadata.ID
				if id == "" {
					id = kace.Snake(gen.Hof.Metadata.Name) + " (auto)"
				}

				name := gen.Hof.Gen.Name
				if name == "" {
					name = "(anon)"
				}
				path := gen.Hof.Path

				row := []string{name, path, id, fmt.Sprint(gen.Hof.Gen.Creator)}
				rows = append(rows, row)
			}
			return rows, nil
		},
	)
}
