package cmd

import (
	"regexp"
	"sort"
	"strings"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/yagu"
	"github.com/olekukonko/tablewriter"
)

func List(args []string, rflags flags.RootPflagpole, cflags flags.Env__ListFlagpole) error {
	// fmt.Println("args:", args)
	cueargs := args
	args = []string{}
	for i, a := range cueargs {
		// fmt.Println("-", i, a, cueargs[:i], cueargs[i:])
		if a == "%" {
			if i+1 < len(cueargs) {
				args = cueargs[i+1:]
			}
			cueargs = cueargs[:i]
			break
		}
	}
	// fmt.Println(cueargs, args)

	R, err := prepRuntime(cueargs, rflags)
	if err != nil {
		return err
	}

	// gather rows
	var rows = make([][]string, 0, len(R.Envs))
	// fill with data
	for _, e := range R.Envs {
		name := e.Hof.Env.Name
		if name == "" {
			continue
		}
		kind := e.Hof.Env.Kind

		do := true
		if len(args) > 0 {
			do = false
			for _, a := range args {
				re, err := regexp.Compile(a)
				if err == nil && re.MatchString(e.Hof.Env.Name) {
					do = true
					break
				}
			}
		}
		if !do {
			continue
		}

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
