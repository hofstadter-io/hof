package hof

import (
	"fmt"
)

func (hn *Node[T]) Print() {
	if hn.Hof.Datamodel.Root {
		hn.printDatamodel("  ", "")
	}

	if hn.Hof.Gen.Root {
		hn.printGen("  ", "")
	}

	if hn.Hof.Flow.Root {
		hn.printFlow("  ", "")
	}
}

func (hn *Node[T]) printDatamodel(spacer, indent string) {
	s := fmt.Sprintf("%s%s:", indent, hn.Hof.Path)
	dm := hn.Hof.Datamodel
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
	
	// print line
	fmt.Println(s)
	// pretty.Println(hn.Hof)

	// recurse
	for _, c := range hn.Children {
		c.printDatamodel(spacer, indent+spacer)
	}
}

func (hn *Node[T]) printGen(spacer, indent string) {
	fmt.Printf("%s%s: @gen(%s)\n", indent, hn.Hof.Path, hn.Hof.Gen.Name)
}

func (hn *Node[T]) printFlow(spacer, indent string) {
	s := fmt.Sprintf("%s%s:", indent, hn.Hof.Path)
	f := hn.Hof.Flow

	if f.Root {
		s += fmt.Sprintf(" @flow(%s)", f.Name)
	}
	if f.Task != "" {
		s += fmt.Sprintf(" @task(%s)", f.Task)
	}
	// print line
	fmt.Println(s)
	// pretty.Println(hn.Hof)

	// recurse
	for _, c := range hn.Children {
		c.printFlow(spacer, indent+spacer)
	}
}
