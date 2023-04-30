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
	found := false
	if dm.Root {
		s += " @datamodel()"
		found = true
	}
	if dm.Node {
		s += " @node()"
		found = true
	}
	if dm.History {
		s += " @history()"
		found = true
	}
	if dm.Ordered {
		s += " @ordered()"
		found = true
	}
	if dm.Cue {
		s += " @cue()"
		found = true
	}
	
	if !found {
		return
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
	for _, c := range hn.Children {
		c.printDatamodel(spacer, indent+spacer)
	}
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
