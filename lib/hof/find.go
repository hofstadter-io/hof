package hof

import (
	"fmt"
	"strconv"
	"strings"

	"cuelang.org/go/cue"
	"github.com/codemodus/kace"
)

// this is where we upgrade our @sugar() to equivalent #hof: { ... } in the Golang Hof type
func upgradeAttrs[T any](node *Node[T], label string) bool {
	val := node.Value
	attrs := val.Attributes(cue.ValueAttr)

	// did we find any attribute
	found := false
	for _, A := range attrs {
		an, ac := A.Name(), A.Contents()
		lfound := true
		switch an {
			case "hof":
				switch ac {
					case "datamodel":
					 node.Hof.Datamodel.Root = true
				}
			case "id":
				node.Hof.Metadata.ID = ac

		case "datamodel":
			 node.Hof.Datamodel.Root = true
		case "history":
			 node.Hof.Datamodel.History = true
		case "ordered":
			 node.Hof.Datamodel.Ordered = true
		case "cue":
			 node.Hof.Datamodel.Cue = true

		// doesn't handle empty case, do we support that
		// we probably should
		case "gen":
			node.Hof.Gen.Root = true
			node.Hof.Gen.Name = ac

		// this doesnt handle empty @flow()
		case "flow":
			node.Hof.Flow.Root = true
			node.Hof.Flow.Name = ac
			// extra, a bit hacky
			node.Hof.Metadata.Name = ac

		// this doesn't handle task names
		// maybe we split into parts
		case "task":
			node.Hof.Flow.Task = ac
			node.Hof.Flow.Name = label

		case "pool":

			parts := strings.Split(ac, ",")

			if len(parts) > 1 {
				node.Hof.Flow.Pool.Name = parts[0]
				node.Hof.Flow.Pool.Make = true
				c, err := strconv.Atoi(parts[1])
				if err != nil {
					fmt.Println("warning: unable to parse %q to int", parts[1])
					
				} else {
					node.Hof.Flow.Pool.Number = c
				}

			} else {
				node.Hof.Flow.Pool.Name = ac
				node.Hof.Flow.Pool.Take = true
			}

		case "chat":
			node.Hof.Chat.Root = true
			node.Hof.Chat.Name = label
			node.Hof.Chat.Extra = ac

		case "print":
			// TODO, better parsing of AC to get parts
			node.Hof.Flow.Print.Level = 1
			node.Hof.Flow.Print.Path  = ac

		default:
		  lfound = false
		}
		// write to outer found
		if lfound {
			found = true
		}
	}

	return found
}

// parse out a #hof for a single value
func ParseHof[T any](val cue.Value) (*Node[T], error) {

	// get some info
	path := val.Path()
	sels := path.Selectors()
	last := cue.Selector{}
	label := ""
	if len(sels) > 0 {
		last = sels[len(sels)-1]
		label = last.String()
	}

	// do not decend into #hof value itself
	// or any definition
	if label == "#hof" {
		return nil, fmt.Errorf("you are trying to parse the #hof value, rather than the containing node at %v", path)
	}

	// return early and recurse for root value?
	if label == "" {
		label = "<<root>>"
		// return nil, nil
	}

	// create new node
	node := New[T](label, val, nil, nil)
	found := false
	
	// look for #hof: _
	hv := val.LookupPath(cue.ParsePath("#hof"))
	if hv.Exists() {
		err := hv.Decode(&(node.Hof))
		if err != nil {
			return node, err
		}
		found = true
	}

	// need to re-add label, path, name? here, after decode
	node.Hof.Label = label
	node.Hof.Path = path.String()

	//
	// look for attributes
	ufound := upgradeAttrs(node, label)
	if ufound {
		found = ufound
	}
	

	// filters to end recursion
	// check datamodel root because of nested history and roots snafu
	// This was here from before we made this singular function, in the context of nodes and the stack
	//if node.Hof.Datamodel.Root {
	//  // backtrack, walking parents		
	//  for bt := nodes; bt != nil; bt = bt.Parent {
	//    // we found a nested root datamodel
	//    if bt.Hof.Datamodel.Root {
	//      // stop recursion
	//      fmt.Println("hof.DM: want to stop recursion here", bt.Hof.Path, node.Hof.Path)
	//      // return false
	//    }	
	//  }
	//  // fmt.Println("found datamodel:", stack.Hof.Path)
	//}

	//// skip any injected History labels?
	//if node.Hof.Label == "History" {
	//  return nil, nil
	//}

	// HACK for nested task, so that it does not get picked up twice
	// (1) inside of a flow  (2) as the root of a flow
	// when we get to rethinking hof/flow with really using the #hof
	// we shouldn't even have something like pool or nest tasks at all
	// (also, CUE might support nested tasks more naturally)
	//if node.Hof.Flow.Pool.Take && node.Hof.Flow.Task == "nest" {
	//  return nil, nil
	//}

	if !found {
		return nil, nil
	}

	// These are EXPENSIVE, only do if we care about and will definitely be returning the node
	// more enrichment
	if node.Hof.Metadata.Name == "" {
		node.Hof.Metadata.Name = node.Hof.Label
		node.Value = node.Value.FillPath(cue.ParsePath("#hof.metadata.name"), node.Hof.Metadata.Name)
	}
	if node.Hof.Metadata.ID == "" {
		node.Hof.Metadata.ID = node.Hof.Metadata.Name
		node.Value = node.Value.FillPath(cue.ParsePath("#hof.metadata.id"), kace.Kebab(node.Hof.Metadata.ID))
	}

	return node, nil
}

func FindHofs(value cue.Value) (roots []*Node[any], err error) {
	// fmt.Println("FindHofs!")
	var stack *Node[any] // cue stack

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

		if label == "#hof" {
			return false
		}

		node, err := ParseHof[any](val)
		if err != nil {
			fmt.Println("while parsing for #hof:", err)
			return false
		}

		// did we find something?
		if node != nil {

			// skip any injected History labels?
			if node.Hof.Label == "History" {
				return false
			}

			// is this the root of that interesting thing?
			// otherwise, push onto the stack and update parent/child pointers
			if stack == nil {
				stack = node	
				roots = append(roots, node)
			} else {
				// two-way relation setting
				node.Parent = stack
				stack.Children = append(stack.Children, node)

				// push stack
				stack = node
			}

			// hmm, should this go before or after the stack update?
			// it doesn't work if we do, so... trying to reconcile this with history
			// this is anything that can have a nested flow, ideally this could be inferred from the node / enrichment
			if node.Hof.Flow.Task == "nest" ||
				(node.Hof.Flow.Name != "" && node.Hof.Flow.Task == ""){
				// fmt.Println("ending hof recursion in nest", node.Hof.Path, node.Hof.Label)
				return false
			}

		}

		return true
	}

	after := func (val cue.Value) {
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
