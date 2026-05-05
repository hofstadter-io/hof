package cmd

import (
	"fmt"
	"strings"

	"cuelang.org/go/cue"
	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/env"
	"github.com/hofstadter-io/hof/lib/yagu"
	"github.com/olekukonko/tablewriter"
)

func List(args []string, rflags flags.RootPflagpole, eflags flags.EnvPflagpole) error {
	// some quick setup and early filtering
	R, matches, err := commonStart(args, rflags, eflags)
	if err != nil {
		return err
	}

	// gather rows
	var rows = make([][]string, 0, len(R.Envs))
	// fill with data
	for _, e := range matches {
		_, kind, mname := extractMeta(e)

		path := e.Hof.Path
		extra := genExtra(e)

		row := []string{mname, kind, path, extra}
		rows = append(rows, row)
	}

	return yagu.PrintAsTable(
		[]string{"Name", "Kind", "Path", "Extra"},
		func(table *tablewriter.Table) ([][]string, error) {
			return rows, nil
		},
	)
}

func genExtra(e *env.Env) string {
	name, kind, mname := extractMeta(e)
	if mname != "" {
		name = mname
	}

	extra := ""
	switch kind {

	case "container":
		b := new(strings.Builder)
		containerExtra(b, e.Value)
		extra = b.String()

	// "name" (HostImage uses this directly)
	case "hostImage":
		extra = "<- " + name

	case "hostExec":
		sv := e.Value.LookupPath(cue.ParsePath("args"))
		if sv.Exists() {
			var args []string
			sv.Decode(&args)
			extra = strings.Join(args, " ")
		}

	case "hostService":
		sv := e.Value.LookupPath(cue.ParsePath("host"))
		if sv.Exists() {
			b := new(strings.Builder)
			s, _ := sv.String()
			fmt.Fprintf(b, "<- %s", s)
			addPorts(b, e.Value)
			extra = b.String()
		}

	case "service":
		b := new(strings.Builder)
		fmt.Fprintf(b, "%s", name)
		addPorts(b, e.Value)
		extra = b.String()

	case "hostTunnel":
		b := new(strings.Builder)
		fmt.Fprintf(b, "-> %s", name)
		addPorts(b, e.Value)
		extra = b.String()

	// "path"
	case "dir", "file", "cuefigSBOM":
		sv := e.Value.LookupPath(cue.ParsePath("path"))
		if sv.Exists() {
			b := new(strings.Builder)
			s, _ := sv.String()
			fmt.Fprintf(b, "%s", s)
			extra = b.String()
		}

	case "hostDir", "hostFile", "hostSocket":
		sv := e.Value.LookupPath(cue.ParsePath("path"))
		if sv.Exists() {
			b := new(strings.Builder)
			s, _ := sv.String()
			fmt.Fprintf(b, "<- %s", s)
			extra = b.String()
		}

	case "exportFile", "exportDir", "exportImageFile":
		sv := e.Value.LookupPath(cue.ParsePath("path"))
		if sv.Exists() {
			b := new(strings.Builder)
			s, _ := sv.String()
			fmt.Fprintf(b, "-> %s", s)
			extra = b.String()
		}

	// "url"
	case "gitRepo":
		sv := e.Value.LookupPath(cue.ParsePath("url"))
		if sv.Exists() {
			s, _ := sv.String()
			extra = s
		}
		rv := e.Value.LookupPath(cue.ParsePath("ref"))
		if sv.Exists() {
			s, _ := rv.String()
			if s != "" {
				extra += "@" + s
			}
		}

	case "exportImage", "publishImage":
		b := new(strings.Builder)
		fmt.Fprintf(b, "-> ")
		rv := e.Value.LookupPath(cue.ParsePath("reg"))
		if rv.Exists() {
			s, _ := rv.String()
			if s != "" {
				fmt.Fprintf(b, "%s/", s)
			}
		}
		nv := e.Value.LookupPath(cue.ParsePath("name"))
		if nv.Exists() {
			s, _ := nv.String()
			fmt.Fprintf(b, "%s", s)
		} else {
			nv := e.Value.LookupPath(cue.ParsePath("image.name"))
			if nv.Exists() {
				s, _ := nv.String()
				fmt.Fprintf(b, "%s", s)
			}
		}
		tv := e.Value.LookupPath(cue.ParsePath("tag"))
		if tv.Exists() {
			s, _ := tv.String()
			if s != "" {
				fmt.Fprintf(b, ":%s", s)
			}
		}
		extra = b.String()

	}

	return extra
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
	// fmt.Println("containerExtra", val)

	// name := val.LookupPath(cue.ParsePath("name"))
	from := val.LookupPath(cue.ParsePath("from"))
	if from.Exists() && from.IncompleteKind() == cue.StringKind {
		// fmt.Println("from: string")
		s, _ := from.String()
		fmt.Fprintf(b, "<- %v", s)
		return
	}
	fromname := val.LookupPath(cue.ParsePath("from.name"))
	if fromname.Exists() && fromname.IncompleteKind() == cue.StringKind {
		// fmt.Println("from: name")
		s, _ := fromname.String()
		fmt.Fprintf(b, "%s", s)
		return
	}

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
