package cmd

import (
	"strings"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/agent"
	"github.com/hofstadter-io/hof/lib/yagu"
	"github.com/olekukonko/tablewriter"
)

func List(args []string, rflags flags.RootPflagpole, aflags flags.AgentPflagpole) error {
	// some quick setup and early filtering
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

func genExtra(a *agent.Agentic) string {
	name, kind, mname := extractMeta(a)
	if mname != "" {
		name = mname
	}

	extra := ""
	switch kind {

	case "agent":
		b := new(strings.Builder)
		// containerExtra(b, a.Value)
		extra = b.String()

	// "name" (HostImage uses this directly)
	case "hostImage":
		extra = "<- " + name
	}

	return extra
}
