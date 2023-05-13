package datamodel

import (
	"fmt"
	"io"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/format"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
)

func (dm *Datamodel) PrintLog(out io.Writer, max int, ts string, dflags flags.DatamodelPflagpole, cflags flags.Datamodel__LogFlagpole) error {
	// we add an extra 2 here because we want to indenty everything under a timestamp
	return dm.T.printLogR(out, ts, "  ", "  ", max+2, dflags, cflags)
}

func (V *Value) printLogR(out io.Writer, ts, indent, spaces string, max int, dflags flags.DatamodelPflagpole, cflags flags.Datamodel__LogFlagpole) error {
	if err := V.printLog(out, ts, indent, spaces, max, dflags, cflags); err != nil {
		return err
	}

	// recurse if children to load any nested histories
	for _, c := range V.Children {
		if err := c.T.printLogR(out, ts, indent + spaces, spaces, max, dflags, cflags); err != nil {
			return err
		}
	}

	return nil
}

func (V *Value) printLog(out io.Writer, ts, indent, spaces string, max int, dflags flags.DatamodelPflagpole, cflags flags.Datamodel__LogFlagpole) error {
	if len(dflags.Expression) > 0 {
		path := V.Hof.Path
		found := false
		for _, ex := range dflags.Expression {
			if ex == path {
				found = true
				break
			}
		}
		if !found {
			return nil
		}
	}

	name := V.Hof.Label
	extra := ""

	hasNoHist := len(V.history) == 0
	// fmt.Println("hasHist:", V.Hof.Path, hasHist)
	if hasNoHist {
		hasMoreHist := V.hasHistBelow()
		if hasMoreHist {
			fmt.Fprintf(out, "%s%s\n", indent, name)
		}
		return nil
	}

	if V.hasSnapshotAt(ts) {
		pos := V.getSnapshotPos(ts)
		newVal := false
		if pos == len(V.history)-1 {
			extra = "+ new value"
			newVal = true
		} else {
			extra = "~ has changes"
		}

		fstr := fmt.Sprintf("%%s%%-%ds %%s\n", max - len(indent))
		fmt.Fprintf(out, fstr, indent, name, extra)

		if cflags.Details {
			S := V.history[pos]
			val := S.Lense.CurrDiff

			if !newVal {
				node := val.Syntax(
					cue.Final(),
					cue.Docs(true),
					cue.Attributes(true),
					cue.Definitions(true),
					cue.Optional(true),
					cue.Hidden(true),
					cue.Concrete(true),
					cue.ResolveReferences(true),
				)
				bytes, err := format.Node(
					node,
					format.Simplify(),
				)
				if err != nil {
					return err
				}
				str := string(bytes)

				lines := strings.Split(str, "\n")
				str = ""
				for _, line := range lines {
					str += indent + line + "\n"
				}
				
				fmt.Fprintln(out, str)
			}

		}
	}

	return nil
}

func (dm *Datamodel) PrintLogByValue(out io.Writer, max int, dflags flags.DatamodelPflagpole, cflags flags.Datamodel__LogFlagpole) error {
		return dm.T.printLogByValueR(out, "", "  ", max, dflags)
}

func (V *Value) printLogByValueR(out io.Writer, indent, spaces string, max int, dflags flags.DatamodelPflagpole) error {
	// load own history
	if err := V.printLogByValue(out, indent, spaces, max, dflags); err != nil {
		return err
	}
	if V.Hof.Datamodel.History {
		V.printLogEntriesByValue(out, indent + spaces, dflags)
	}

	// recurse if children to load any nested histories
	for _, c := range V.Children {
		if err := c.T.printLogByValueR(out, indent + spaces, spaces, max, dflags); err != nil {
			return err
		}
	}

	return nil
}

func (V *Value) printLogByValue(out io.Writer, indent, spaces string, max int, dflags flags.DatamodelPflagpole) error {
	if len(dflags.Expression) > 0 {
		path := V.Hof.Path
		found := false
		for _, ex := range dflags.Expression {
			if ex == path {
				found = true
				break
			}
		}
		if !found {
			return nil
		}
	}
	name := V.Hof.Label
	extra := ""
	if V.hasDiff() {
		extra = "~ has changes"
	}
	if V.Hof.Datamodel.History && len(V.history) == 0 {
		extra = "+ new value"
	}

	fstr := fmt.Sprintf("%%s%%-%ds %%s\n", max - len(indent))
	fmt.Fprintf(out, fstr, indent, name, extra)

	return nil
}

func (V *Value) printLogEntriesByValue(out io.Writer, indent string, dflags flags.DatamodelPflagpole) error {
	if len(dflags.Expression) > 0 {
		path := V.Hof.Path
		found := false
		for _, ex := range dflags.Expression {
			if ex == path {
				found = true
				break
			}
		}
		if !found {
			return nil
		}
	}
	// fmt.Println("got here", len(V.history))
	for _, S := range V.history {
		fmt.Fprintf(out, "%s  @%s  %q\n", indent, S.Timestamp, S.Message)
	}

	return nil
}
