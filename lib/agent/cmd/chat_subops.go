package cmd

import (
	"fmt"
	"maps"
	"sort"
	"strings"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/agent/runtime/handlers/common"
	"github.com/hofstadter-io/hof/lib/consts"
	"github.com/hofstadter-io/hof/lib/yagu"
	"github.com/olekukonko/tablewriter"
)

func ChatInfo(args []string, rflags flags.RootPflagpole, aflags flags.AgentPflagpole, cflags flags.Agent__ChatPflagpole) error {
	// hard code the kind here, as we only support one at this point for chat
	aflags.Kind = []string{"agent"}
	R, _, matches, err := commonStart(args, rflags, aflags)
	if err != nil {
		return err
	}

	fmt.Println("agentics found:", len(R.Agentics), len(matches))

	return nil
}

func ChatList(args []string, rflags flags.RootPflagpole, aflags flags.AgentPflagpole, cflags flags.Agent__ChatPflagpole) error {
	R, AR, _, err := commonStart(args, rflags, aflags)
	if err != nil {
		return err
	}

	// these loops need to be switched
	// REALLY, we need a third list to calc intermediate rows with comparable data
	sessions, err := common.SessionList(R.Ctx, AR, consts.VEG_DEFAULT_USER)
	rows := make([][]string, 0, len(sessions))
	for _, s := range sessions {
		id := s.ID()

		numState := fmt.Sprintf("%2d", len(maps.Collect(s.State().All())))
		lastTime := s.LastUpdateTime().Local().Format("Mon, Jan 2, 2006 15:04")
		tVal, err := s.State().Get("title")
		title := ""
		if err == nil {
			title = tVal.(string)
		}

		row := []string{title, id, numState, lastTime}
		rows = append(rows, row)
	}

	if len(aflags.Sort) == 0 {
		aflags.Sort = []string{"time"}
	}
	sort.Slice(rows, func(i, j int) bool {
		lhs, rhs := rows[i], rows[j]
		for _, s := range aflags.Sort {
			s = strings.ToLower(s)
			switch s {
			case "name", "title":
				return lhs[0] < rhs[0]
			case "id":
				return lhs[1] < rhs[1]
			case "state":
				return lhs[2] < rhs[2]
			case "time", "last", "update", "lastupdate":
				return lhs[3] < rhs[3]
			default:
				fmt.Printf("WARN: unknown sort field %q ignored", s)
				continue
			}
		}
		return false
	})

	return yagu.PrintAsTable(
		[]string{"Title", "ID", "|State|", "Last Update"},
		func(table *tablewriter.Table) ([][]string, error) {
			return rows, nil
		},
	)
}

func ChatDelete(args []string, rflags flags.RootPflagpole, aflags flags.AgentPflagpole, cflags flags.Agent__ChatPflagpole) error {
	// hard code the kind here, as we only support one at this point for chat
	aflags.Kind = []string{"agent"}
	R, _, matches, err := commonStart(args, rflags, aflags)
	if err != nil {
		return err
	}

	fmt.Println("agentics found:", len(R.Agentics), len(matches))

	return nil
}
