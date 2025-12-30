package cmd

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"cuelang.org/go/cue"
	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/yagu"
	"github.com/olekukonko/tablewriter"
)

func List(args []string, rflags flags.RootPflagpole, cflags flags.Env__ListFlagpole) error {
	args, cueargs := splitArgs(args)

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

		extra := ""
		switch kind {

		case "container":
			b := new(strings.Builder)
			containerExtra(b, e.Value)
			extra = b.String()

		// "name" (HostImage uses this directly)
		case "hostImage":
			sv := e.Value.LookupPath(cue.ParsePath("name"))
			if sv.Exists() {
				s, _ := sv.String()
				extra = s
			}
		case "hostService":
			sv := e.Value.LookupPath(cue.ParsePath("host"))
			if sv.Exists() {
				b := new(strings.Builder)
				s, _ := sv.String()
				fmt.Fprintf(b, "%s", s)
				addPorts(b, e.Value)
				extra = b.String()
			}
		case "service", "hostTunnel":
			sv := e.Value.LookupPath(cue.ParsePath("name"))
			if sv.Exists() {
				b := new(strings.Builder)
				s, _ := sv.String()
				fmt.Fprintf(b, "%s", s)
				addPorts(b, e.Value)
				extra = b.String()
			}

		// "path"
		case "dir", "file", "hostDir", "hostFile", "hostSocket", "exportFile", "exportDir", "exportImageFile":
			sv := e.Value.LookupPath(cue.ParsePath("path"))
			if sv.Exists() {
				b := new(strings.Builder)
				s, _ := sv.String()
				fmt.Fprintf(b, "%s", s)
				extra = b.String()
			}

		// "url"
		case "gitRepo", "exportImage":
			sv := e.Value.LookupPath(cue.ParsePath("url"))
			if sv.Exists() {
				s, _ := sv.String()
				extra = s
			}
		}

		row := []string{name, kind, path, extra}
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
		[]string{"Name", "Kind", "Path", "Extra"},
		func(table *tablewriter.Table) ([][]string, error) {
			return rows, nil
		},
	)
}

func addPorts(b *strings.Builder, val cue.Value) {
	ports := val.LookupPath(cue.ParsePath("ports"))
	if !ports.Exists() {
		return
	}
	iter, _ := ports.List()
	for iter.Next() {
		pv := iter.Value()
		bev := pv.LookupPath(cue.ParsePath("backend"))
		be, _ := bev.Int64()
		if be > 0 {
			fmt.Fprintf(b, ":%d", be)
		}
		fev := pv.LookupPath(cue.ParsePath("frontend"))
		fe, _ := fev.Int64()
		if fe > 0 {
			fmt.Fprintf(b, ":%d", fe)
		}
	}
}

func containerExtra(b *strings.Builder, val cue.Value) {
	// name := val.LookupPath(cue.ParsePath("name"))
	from := val.LookupPath(cue.ParsePath("from.name"))
	s, _ := from.String()
	fmt.Fprintf(b, "from: %v", s)

	// switch ik := from.IncompleteKind(); ik {
	// case cue.StringKind:
	// 	s, _ := from.String()
	// 	fmt.Fprintf(b, "from: %v", s)
	// case cue.StructKind:
	// 	name := val.LookupPath(cue.ParsePath("name"))
	// 	s, _ := name.String()
	// 	fmt.Fprintf(b, "from: %v", s)
	// }
}
