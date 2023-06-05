package hof

import (
	"fmt"
	"strings"

	"cuelang.org/go/cue"
	"github.com/codemodus/kace"
)

func FindHofs(value cue.Value) (roots []*Node[any], err error) {
	// fmt.Println("FindHofs!")
	var stack *Node[any] // cue stack
	var nodes *Node[any] // hof nodes

	before := func (val cue.Value) bool {
		// get some info
		path := val.Path()
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

		// do not decend into $hof value itself
		// or any definition
		if label == "#hof" {
			return false
		}

		// update cue stack
		curr := New[any](label, val, nil, stack)
		stack = curr


		// did we find something of interest?
		found := false

		// look for #hof: _
		hv := val.LookupPath(cue.ParsePath("#hof"))
		if hv.Exists() {
			found = true
			err = hv.Decode(&(stack.Hof))
			if err != nil {
				fmt.Println(err)
				stack = stack.Parent
				return false
			}

			// fmt.Println("id:", stack.Hof.Label, stack.Hof.Path, stack.Hof.Metadata)
		}

		//
		// Do not modify stack.Hof before we Decode it in the condition just above
		//

		// need to re-add label, path, name? here, after decode
		stack.Hof.Label = label
		stack.Hof.Path = val.Path().String()

		// look for attributes
		attrs := val.Attributes(cue.ValueAttr)
		for _, A := range attrs {
			an, ac := A.Name(), A.Contents()
			found = true
			switch an {
				case "hof":
					switch ac {
						case "datamodel":
						 stack.Hof.Datamodel.Root = true
					}
				case "id":
					stack.Hof.Metadata.ID = ac

			case "datamodel":
				 stack.Hof.Datamodel.Root = true
			case "history":
				 stack.Hof.Datamodel.History = true
			case "ordered":
				 stack.Hof.Datamodel.Ordered = true
			case "cue":
				 stack.Hof.Datamodel.Cue = true

			// doesn't handle empty case, do we support that
			// we probably should
			case "gen":
				stack.Hof.Gen.Root = true
				stack.Hof.Gen.Name = ac

			// this doesnt handle empty @flow()
			case "flow":
				stack.Hof.Flow.Root = true
				stack.Hof.Flow.Name = ac
			// this doesn't handle task names
			// maybe we split into parts
			case "task":
				stack.Hof.Flow.Task = ac
				stack.Hof.Flow.Name = label

			case "chat":
				stack.Hof.Chat.Root = true
				stack.Hof.Chat.Name = label
				stack.Hof.Chat.Extra = ac

			default:
				found = false
			}
		}

		// filters to end recursion
		// check datamodel root because of nested history and roots snafu
		if stack.Hof.Datamodel.Root {
			// backtrack, walking parents		
			for bt := nodes; bt != nil; bt = bt.Parent {
				// we found a nested root datamodel
				if bt.Hof.Datamodel.Root {
					// stop recursion
					fmt.Println("hof.DM: want to stop recursion here", bt.Hof.Path, stack.Hof.Path)
					// return false
				}	
			}

			// fmt.Println("found datamodel:", stack.Hof.Path)
		}
		if stack.Hof.Label == "History" {
			stack = stack.Parent
			return false
		}

		// we should update the nodes
		if found {
			// update hof node
			node := New[any](label, val, nil, nodes)
			nodes = node
			nodes.Hof = stack.Hof

			// more enrichment
			if nodes.Hof.Metadata.Name == "" {
				nodes.Hof.Metadata.Name = nodes.Hof.Label
				nodes.Value = nodes.Value.FillPath(cue.ParsePath("#hof.metadata.name"), nodes.Hof.Metadata.Name)
			}
			if nodes.Hof.Metadata.ID == "" {
				nodes.Hof.Metadata.ID = nodes.Hof.Metadata.Name
				nodes.Value = nodes.Value.FillPath(cue.ParsePath("#hof.metadata.id"), kace.Kebab(nodes.Hof.Metadata.ID))
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

	after := func (val cue.Value) {
		// paths for matching trees
		np, sp := "", ""
		if nodes != nil {
			np = nodes.Value.Path().String()
		}
		if stack != nil {
			sp = stack.Value.Path().String()
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
	cue.Definitions(true),
	cue.Hidden(false),
	cue.Optional(true),
	cue.Docs(true),
}

// Walk is an alternative to cue.Value.Walk which handles more field types
// You can customize this with your own options
// returning false will stop recursion for that node
func Walk(v cue.Value, before func(cue.Value) bool, after func(cue.Value), options ...cue.Option) {

	// call before and possibly stop recursion
	if before != nil && !before(v) {
		return
	}

	// possibly recurse
	switch v.IncompleteKind() {
	case cue.StructKind:
		if options == nil {
			options = defaultWalkOptions
		}
		s, _ := v.Fields(options...)

		for s.Next() {
			Walk(s.Value(), before, after, options...)
		}

	case cue.ListKind:
		l, _ := v.List()
		for l.Next() {
			Walk(l.Value(), before, after, options...)
		}

		// no default (basic lit types)

	}

	if after != nil {
		after(v)
	}

}
