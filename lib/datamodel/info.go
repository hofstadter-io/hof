package datamodel

import (
	"fmt"
	"io"
	"strings"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
)

func (dm *Datamodel) PrintInfo(out io.Writer, max int, dflags flags.DatamodelPflagpole) error {
	return dm.T.printInfoR(out, "", "  ", max, dflags)
}

func (V *Value) printInfoR(out io.Writer, indent, spaces string, max int, dflags flags.DatamodelPflagpole) error {
	// print current info
	if err := V.printInfo(out, indent, spaces, max, dflags); err != nil {
		return err
	}

	// recurse into any child nodes
	for _, c := range V.Children {
		if err := c.T.printInfoR(out, indent + spaces, spaces, max, dflags); err != nil {
			return err
		}
	}

	return nil
}

func (V *Value) printInfo(out io.Writer, indent, spaces string, max int, dflags flags.DatamodelPflagpole) error {
	if len(dflags.Expression) > 0 {
		path := V.Hof.Path
		found := false
		for _, ex := range dflags.Expression {
			if strings.HasPrefix(path, ex) {
				found = true
				break
			}
		}
		if !found {
			return nil
		}
	}
	name := V.Hof.Label
	dm := V.Hof.Datamodel

	s := ""
	if dm.Root {
		s += " @datamodel()"
	}
	if dm.Node {
		s += " @node()"
	}
	if dm.History {
		s += " @history()"
	}
	if dm.Ordered {
		s += " @ordered()"
	}
	if dm.Cue {
		s += " @cue()"
	}

	fstr := fmt.Sprintf("%%s%%-%ds %%s\n", max - len(indent))
	fmt.Fprintf(out, fstr, indent, name, s)

	return nil
}
