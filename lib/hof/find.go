package hof

import (
	"fmt"
	"strings"

	"cuelang.org/go/cue"
	"github.com/codemodus/kace"
)

func FindHofs(value *Value) (roots []*Node[any], err error) {
	fmt.Println("FindHofs:", value.Path())
	var stack *Node[any] // cue stack
	var nodes *Node[any] // hof nodes

	before := func (val *Value) bool {
		// fmt.Println("find:", val.Path())
		// get some info
		path := val.CueValue().Path()
		sels := path.Selectors()
		last := cue.Selector{}
		label := ""
		if len(sels) > 0 {
			last = sels[len(sels)-1]
			label = last.String()
		}

		// return early and recurse for root value
		if label == "" {
			return true
		}

		curr := New[any](label, path.String(), value, nil, stack)
		stack = curr

		// do not decend into $hof value itself
		// or any definition
		if label == "$hof" {
			stack = stack.Parent
			return false
		}

		var tmp Node[any]

		// did we find something of interest?
		found := false

		// look for $hof: _
		hv := val.LookupPath("$hof")
		if hv.Exists() {
			found = true
			// fmt.Println("id:", stack.Hof.Label, stack.Hof.Path, stack.Hof.Metadata)
			err = hv.Decode(&tmp.Hof)
			if err != nil {
				fmt.Println(err)
				return false
			}

		}

		tmp.Hof.Label = label
		tmp.Hof.Path = val.Path()

		// look for attributes
		attrs := val.CueValue().Attributes(cue.ValueAttr)
		for _, A := range attrs {
			an, ac := A.Name(), A.Contents()
			found = true
			switch an {
				case "hof":
					switch ac {
						case "datamodel":
						 tmp.Hof.Datamodel.Root = true
					}
				case "id":
					tmp.Hof.Metadata.ID = ac

			case "datamodel":
				 tmp.Hof.Datamodel.Root = true
			case "history":
				 tmp.Hof.Datamodel.History = true
			case "ordered":
				 tmp.Hof.Datamodel.Ordered = true
			case "cue":
				 tmp.Hof.Datamodel.Cue = true

			// doesn't handle empty case, do we support that
			// we probably should
			case "gen":
				tmp.Hof.Gen.Root = true
				tmp.Hof.Gen.Name = ac

			// this doesnt handle empty @flow()
			case "flow":
				tmp.Hof.Flow.Root = true
				tmp.Hof.Flow.Name = ac
			// this doesn't handle task names
			// maybe we split into parts
			case "task":
				tmp.Hof.Flow.Task = ac
				tmp.Hof.Flow.Name = label

			case "chat":
				tmp.Hof.Chat.Root = true
				tmp.Hof.Chat.Name = label
				tmp.Hof.Chat.Extra = ac

			default:
				found = false
			}
		}

		// filters to end recursion
		// check datamodel root because of nested history and roots snafu
		if tmp.Hof.Datamodel.Root {
			// backtrack, walking parents		
			for bt := nodes; bt != nil; bt = bt.Parent {
				// we found a nested root datamodel
				if bt.Hof.Datamodel.Root {
					// stop recursion
					return false
				}	
			}
		}

		if tmp.Hof.Label == "History" {
			return false
		}

		// we should update the nodes
		if found {
			fmt.Println("found:", tmp.Hof.Path)



			// update hof node
			node := New[any](label, stack.Hof.Path, value, nil, nodes)
			nodes = node
			nodes.Hof = stack.Hof

			// more enrichment
			if nodes.Hof.Metadata.Name == "" {
				nodes.Hof.Metadata.Name = nodes.Hof.Label
				path := nodes.Hof.Path + ".$hof.metadata.name"
				value.FillPath(path, nodes.Hof.Metadata.Name)
			}
			if nodes.Hof.Metadata.ID == "" {
				nodes.Hof.Metadata.ID = nodes.Hof.Metadata.Name
				path := nodes.Hof.Path + ".$hof.metadata.id"
				value.FillPath(path, kace.Kebab(nodes.Hof.Metadata.ID))
			}

			if nodes.Parent == nil {
				// add to root if no parent
				roots = append(roots, nodes)
			} else {
				// add to parent's Children
				nodes.Parent.Children = append(nodes.Parent.Children, nodes)
			}
		}

		return true
	}

	after := func (val *Value) {
		// paths for matching trees
		np, sp := "", ""
		if nodes != nil {
			np = nodes.Hof.Path
		}
		if stack != nil {
			sp = stack.Hof.Path
		}

		if nodes != nil && np == sp {
			// unwind hof nodes
			nodes = nodes.Parent
		}
		// unwind node stack
		if stack != nil {
			stack = stack.Parent
		}
	}

	// this is a depth first walk
	Walk(value, before, after)

	return roots, nil
}

func (n *Node[T]) indent() string {
	if n == nil {
		return ""
	}
	d := 1
	for n.Parent != nil {
		d++
		n = n.Parent
	}
	return strings.Repeat("  ", d)
}

var defaultWalkOptions = []cue.Option{
	cue.Attributes(true),
	cue.Concrete(false),
	cue.Definitions(false),
	cue.Hidden(false),
	cue.Optional(true),
	cue.Docs(false),
}

// Walk is an alternative to cue.Value.Walk which handles more field types
// You can customize this with your own options
// returning false will stop recursion for that node
func Walk(v *Value, before func(*Value) bool, after func(*Value), options ...cue.Option) {

	// call before and possibly stop recursion
	if before != nil && !before(v) {
		return
	}

	// possibly recurse
	switch v.CueValue().IncompleteKind() {
	case cue.StructKind:
		if options == nil {
			options = defaultWalkOptions
		}
		s, _ := v.Fields(options...)

		for s.Next() {
			f := v.LookupPath(s.Selector().String())
			Walk(f, before, after, options...)
		}

	case cue.ListKind:
		l, _ := v.List()
		for l.Next() {
			e := v.LookupPath(l.Selector().String())
			Walk(e, before, after, options...)
		}

		// no default (basic lit types)

	}

	if after != nil {
		after(v)
	}

}
