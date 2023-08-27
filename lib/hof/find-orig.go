package hof

import (
	"fmt"

	"cuelang.org/go/cue"
	"github.com/codemodus/kace"
)

func FindHofsOrig(value cue.Value) (roots []*Node[any], err error) {
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
			label = "<<root>>"
			// return true
		}

		// do not decend into $hof value itself
		// or any definition
		if label == "#hof" {
			return false
		}
		if label == "History" {
			return false
		}

		// update cue stack
		curr := New(label, val, nil, stack)
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

		ufound := upgradeAttrs(stack, label)
		if ufound {
			found = ufound
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
		//if stack.Hof.Label == "History" {
		//  stack = stack.Parent
		//  return false
		//}

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
	WalkOrig(value, before, after)

	return roots, nil
}

// Walk is an alternative to cue.Value.Walk which handles more field types
// You can customize this with your own options
// returning false will stop recursion for that node
func WalkOrig(v cue.Value, before func(cue.Value) bool, after func(cue.Value), options ...cue.Option) {

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
