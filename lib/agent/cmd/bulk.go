package cmd

import (
	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/yagu"
	"github.com/olekukonko/tablewriter"
)

func Bulk(args []string, rflags flags.RootPflagpole, aflags flags.AgentPflagpole) error {
	R, _, matches, err := commonStart(args, rflags, aflags)
	if err != nil {
		return err
	}

	// gather rows
	var rows = make([][]string, 0, len(R.Agentics))
	// fill with data
	for _, e := range matches {
		_, kind, mname := extractMeta(e)

		path := e.Hof.Path
		// extra := genExtra(e)

		row := []string{mname, kind, path}
		rows = append(rows, row)
	}

	return yagu.PrintAsTable(
		[]string{"Name", "Kind", "Path"},
		func(table *tablewriter.Table) ([][]string, error) {
			return rows, nil
		},
	)
}
