package cmd

import (
	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/yagu"
	"github.com/olekukonko/tablewriter"
)

func List(args []string, rflags flags.RootPflagpole) error {
	R, err := prepRuntime(args, rflags)
	if err != nil {
		return err
	}

	return yagu.PrintAsTable(
		[]string{"Name", "Kind", "Path"},
		func(table *tablewriter.Table) ([][]string, error) {
			var rows = make([][]string, 0, len(R.Envs))
			// fill with data
			for _, e := range R.Envs {
				name := e.Hof.Env.Name
				if name == "" {
					name = "(anon)"
				}
				kind := e.Hof.Env.Kind
				path := e.Hof.Path

				row := []string{name, kind, path}
				rows = append(rows, row)
			}
			return rows, nil
		},
	)
}
