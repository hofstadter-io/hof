package cmd

import (
	"sort"
	"strings"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/yagu"
	"github.com/olekukonko/tablewriter"
)

func List(args []string, rflags flags.RootPflagpole, cflags flags.Env__ListFlagpole) error {
	R, err := prepRuntime(args, rflags)
	if err != nil {
		return err
	}

	// gather rows
	var rows = make([][]string, 0, len(R.Envs))
	// fill with data
	for _, e := range R.Envs {
		name := e.Hof.Env.Name
		if name == "" {
			name = "(anon)"
		}
		kind := e.Hof.Env.Kind

		if len(cflags.Kind) > 0 {
			match := false
			for _, k := range cflags.Kind {
				if kind == k {
					match = true
					break
				}
			}
			if !match {
				continue
			}
		}

		path := e.Hof.Path

		row := []string{name, kind, path}
		rows = append(rows, row)
	}

	// multi-column sort based on cflags.Sort ([]string)
	if len(cflags.Sort) > 0 {
		sort.Slice(rows, func(i, j int) bool {
			for _, s := range cflags.Sort {
				s = strings.ToLower(s)
				idx := -1
				switch s {
				case "name":
					idx = 0
				case "kind":
					idx = 1
				case "path":
					idx = 2
				}

				if idx == -1 {
					continue
				}

				if rows[i][idx] != rows[j][idx] {
					return rows[i][idx] < rows[j][idx]
				}
			}
			return false
		})
	}

	return yagu.PrintAsTable(
		[]string{"Name", "Kind", "Path"},
		func(table *tablewriter.Table) ([][]string, error) {
			return rows, nil
		},
	)
}
