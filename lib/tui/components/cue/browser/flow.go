package browser

import (
	"bytes"
	"fmt"
	"strings"

	"cuelang.org/go/cue"
	cueflow "cuelang.org/go/tools/flow"
	"github.com/gdamore/tcell/v2"

	"github.com/hofstadter-io/hof/lib/cuetils"
	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/tview"
)

func (VB *Browser) onFlowSelect(node *tview.TreeNode) {
	reference := node.GetReference()
	if reference == nil {
		return // Selecting the root node does nothing.
	}

	children := node.GetChildren()
	if len(children) == 0 {
		// Load and show files in this directory.
		path := reference.(string)
		VB.FlowAddAt(node, path)
	} else {
		// Collapse if visible, expand if collapsed.
		node.SetExpanded(!node.IsExpanded())
	}
}


func (VB *Browser) FlowAddAt(target *tview.TreeNode, path string) {
	tui.Log("info", fmt.Sprintf("VB.FlowAddAt: %s", path))
	if VB.flow == nil {
		return
	}

	tasks := VB.flow.Ctrl.Tasks()
	if strings.HasPrefix(path, "<root>") {
		for _, task := range tasks {
			path := task.Path().String()
			// TODO, count fields here (even different types)
			l := fmt.Sprintf("%s", path)
			node := tview.NewTreeNode(l)
			node.
				SetColor(tcell.ColorCornflowerBlue).
				SetSelectable(true)
			node.SetReference(path)
			target.AddChild(node)
		}
	}

	if strings.HasPrefix(path, ".") {
		path = path[1:]
	}

	var T *cueflow.Task

	for _, task := range tasks {
		tpath := task.Path().String()
		if path == tpath {
			// exact path to task
			T = task

			//
			// add a bunch of fields by hand
			//

			node := tview.NewTreeNode(fmt.Sprintf("%s", T.State()))
			node.SetColor(tcell.ColorWhite).SetSelectable(false)
			node.SetReference(path + ".State")
			target.AddChild(node)

			node = tview.NewTreeNode("Value")
			node.SetColor(tcell.ColorCornflowerBlue).SetSelectable(true)
			node.SetReference(path + ".Value")
			target.AddChild(node)

			if len(T.Dependencies()) > 0 {
				node = tview.NewTreeNode("Deps")
				node.SetColor(tcell.ColorCornflowerBlue).SetSelectable(true)
			} else {
				node = tview.NewTreeNode("No Deps")
				node.SetColor(tcell.ColorWhite).SetSelectable(false)
			}
			node.SetReference(path + ".Deps")
			target.AddChild(node)

			node = tview.NewTreeNode("Stats")
			node.SetColor(tcell.ColorCornflowerBlue).SetSelectable(true)
			node.SetReference(path + ".Stats")
			target.AddChild(node)

		} else if strings.HasPrefix(path, tpath) {
			// path to task sub-field
			T = task
			spath := strings.TrimPrefix(path, tpath + ".")

			parts := strings.Split(spath, ".")

			first := parts[0]
			parts = parts[1:]

			switch first {
			case "Value":
				p := strings.Join(parts, ".")
				VB.FlowValueAddAt(target, T.Value(), tpath + ".Value", p)
			case "Deps":
				if len(parts) == 0 {
					for _, D := range T.Dependencies() {
						p := fmt.Sprint(D.Path())
						node := tview.NewTreeNode(p)
						node.SetColor(tcell.ColorWhite).SetSelectable(false)
						node.SetReference(path + ".Deps." + p)
						target.AddChild(node)
					}
				}
			case "Stats":
				stats := T.Stats()

				node := tview.NewTreeNode(fmt.Sprintf("Unifications: %d", stats.Unifications))
				node.SetColor(tcell.ColorWhite).SetSelectable(false)
				node.SetReference(path + ".Stats.Unifications")
				target.AddChild(node)

				node = tview.NewTreeNode(fmt.Sprintf("Disjuncts: %d", stats.Disjuncts))
				node.SetColor(tcell.ColorWhite).SetSelectable(false)
				node.SetReference(path + ".Stats.Disjuncts")
				target.AddChild(node)

				node = tview.NewTreeNode(fmt.Sprintf("Conjuncts: %d", stats.Conjuncts))
				node.SetColor(tcell.ColorWhite).SetSelectable(false)
				node.SetReference(path + ".Stats.Conjuncts")
				target.AddChild(node)

				node = tview.NewTreeNode(fmt.Sprintf("Freed: %d", stats.Freed))
				node.SetColor(tcell.ColorWhite).SetSelectable(false)
				node.SetReference(path + ".Stats.Freed")
				target.AddChild(node)

				node = tview.NewTreeNode(fmt.Sprintf("Reused: %d", stats.Reused))
				node.SetColor(tcell.ColorWhite).SetSelectable(false)
				node.SetReference(path + ".Stats.Reused")
				target.AddChild(node)

				node = tview.NewTreeNode(fmt.Sprintf("Allocs: %d", stats.Allocs))
				node.SetColor(tcell.ColorWhite).SetSelectable(false)
				node.SetReference(path + ".Stats.Allocs")
				target.AddChild(node)

				node = tview.NewTreeNode(fmt.Sprintf("Retained: %d", stats.Retained))
				node.SetColor(tcell.ColorWhite).SetSelectable(false)
				node.SetReference(path + ".Stats.Retained")
				target.AddChild(node)

				node = tview.NewTreeNode(fmt.Sprintf("Leaks: %d", stats.Leaks()))
				node.SetColor(tcell.ColorWhite).SetSelectable(false)
				node.SetReference(path + ".Stats.Leaks")
				target.AddChild(node)

			}
		}

	}

}

func (VB *Browser) FlowValueAddAt(target *tview.TreeNode, value cue.Value, fpath, path string) {
	// tui.Log("debug", fmt.Sprintf("VB.FlowValueAddAt: %q %q\n", fpath, path))

	if strings.HasPrefix(path, "<root>") {
		path = ""
	}
	if strings.HasPrefix(path, ".") {
		path = path[1:]
	}
	val := value.LookupPath(cue.ParsePath(path))

	if val.Err() != nil {
		tui.SendCustomEvent("/console/err", cuetils.CueErrorToString(val.Err()))
		return
	}

	// get fields at path, need to know what format options are at play here
	var iter *cue.Iterator
	switch val.IncompleteKind() {
	case cue.StructKind:
		iter, _ = val.Fields(VB.Options()...)
	case cue.ListKind:
		// because List does not return a pointer like Fields
		i, _ := val.List()
		iter = &i
	}

	// val is not an iterable CUE type
	if iter == nil {
		// VB.App.Logger(fmt.Sprintf("nil iter for: %s\n", path))
		return
	}

	// build tree nodes
	for iter.Next() {
		sel := iter.Selector()
		value := iter.Value()
		attrs := value.Attributes(cue.ValueAttr)

		fullpath := path
		// input value that we are iterating over
		switch val.IncompleteKind() {
		case cue.ListKind:
			fullpath += fmt.Sprintf("[%s]", sel)
		default:
			fullpath += fmt.Sprintf(".%s", sel)
		}


		var node *tview.TreeNode

		var buf bytes.Buffer
		for _, a := range attrs {
			fmt.Fprintf(&buf, "%v", a)
		}
		attr := buf.String()

		switch value.IncompleteKind() {
		case cue.StructKind:
			// TODO, count fields here (even different types)
			l := fmt.Sprintf("{ %s }", sel)
			line := fmt.Sprintf("%-42s [goldenrod]%s", l, attr)
			node = tview.NewTreeNode(line)
			node.
				SetColor(tcell.ColorCornflowerBlue).
				SetSelectable(true)

		case cue.ListKind:
			l := fmt.Sprintf("[ %s (%d) ]", sel, value.Len())
			node = tview.NewTreeNode(l)
			node.
				SetColor(tcell.ColorLime).
				SetSelectable(true)
			

		default:
			l := fmt.Sprintf("%s: %v %s", sel, value, attr)
			node = tview.NewTreeNode(l)
			node.SetSelectable(false)

		}

		if strings.HasPrefix(fullpath, ".") {
			fullpath = fullpath[1:]
		}
		node.SetReference(fpath + "." + fullpath)
		target.AddChild(node)
	}
}



