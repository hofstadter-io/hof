package components

import (
	"bytes"
	"fmt"
	"strings"

	"cuelang.org/go/cue"
	"github.com/rivo/tview"
	"github.com/gdamore/tcell/v2"

	"github.com/hofstadter-io/hof/lib/tui/app"
)

type ValueBrowser struct {
	*tview.TreeView

	OnFieldSelect func(string)

	Root *tview.TreeNode

	App *app.App

	Value cue.Value

	docs,
	attrs,
	defs,
	optional,
	ignore,
	inline,
	resolve,
	concrete,
	hidden,
	final bool
}

func NewValueBrowser(app *app.App, val cue.Value, OnFieldSelect func(path string)) *ValueBrowser {
	FB := &ValueBrowser {
		App: app,
		Value: val,
		OnFieldSelect: OnFieldSelect,
	}

	// tree view
	FB.TreeView = tview.NewTreeView()
	FB.Root = tview.NewTreeNode("no results yet")
	FB.Root.SetColor(tcell.ColorSilver)

	FB.
		SetRoot(FB.Root).
		SetCurrentNode(FB.Root)

	// set our selected handler
	FB.SetSelectedFunc(FB.OnSelect)
	FB.SetBorder(true)

	FB.SetupKeybinds()

	return FB
}

func (FB *ValueBrowser) Rebuild(path string) {
	if path == "" {
		path = "<root>"
	}

	FB.Root = tview.NewTreeNode(path)
	FB.Root.SetColor(tcell.ColorSilver)
	FB.AddAt(FB.Root, path)

	FB.
		SetRoot(FB.Root).
		SetCurrentNode(FB.Root)

}

func (FB *ValueBrowser) OnSelect(node *tview.TreeNode) {
	reference := node.GetReference()
	if reference == nil {
		return // Selecting the root node does nothing.
	}

	children := node.GetChildren()
	if len(children) == 0 {
		// Load and show files in this directory.
		path := reference.(string)
		FB.AddAt(node, path)
	} else {
		// Collapse if visible, expand if collapsed.
		node.SetExpanded(!node.IsExpanded())
	}
}


func (FB *ValueBrowser) AddAt(target *tview.TreeNode, path string) {
	// FB.App.Logger(fmt.Sprintf("FB.AddAt: %s\n", path))

	if strings.HasPrefix(path, "<root>") {
		path = ""
	}
	if strings.HasPrefix(path, ".") {
		path = path[1:]
	}
	val := FB.Value.LookupPath(cue.ParsePath(path))
	// FB.App.Logger(fmt.Sprintf("#v\n", val))

	if val.Err() != nil {
		FB.App.Logger(fmt.Sprintf("Error: %s\n", val.Err()))
		return
	}

	// get fields at path, need to know what format options are at play here
	var iter *cue.Iterator
	switch val.IncompleteKind() {
	case cue.StructKind:
		opts := []cue.Option{
			cue.Docs(FB.docs),
			cue.Attributes(FB.attrs),
			cue.Definitions(FB.defs),
			cue.Optional(FB.optional),
			cue.InlineImports(FB.inline),
			cue.ErrorsAsValues(FB.ignore),
			cue.ResolveReferences(FB.resolve),
		}
		if FB.concrete {
			opts = append(opts, cue.Concrete(true))
		}
		if FB.hidden {
			opts = append(opts, cue.Hidden(true))
		}

		if FB.final {
			// prepend final, so others still apply
			opts = append([]cue.Option{cue.Final()}, opts...)
		}

		iter, _ = val.Fields(opts...)
	case cue.ListKind:
		i, _ := val.List()
		iter = &i
	}
	if iter == nil {
		FB.App.Logger(fmt.Sprintf("nil iter for: %s\n", path))
		return
	}

	// sort dirs first, then by name
	//sort.Slice(files, func(x, y int) bool {
	//  X, Y := files[x], files[y]
	//  // deal with file vs dir
	//  if X.IsDir() && !Y.IsDir() {
	//    return true
	//  } else if !X.IsDir() && Y.IsDir() {
	//    return false
	//  } else {
	//    return X.Name() < Y.Name()
	//  }
	//})

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
			l := fmt.Sprintf("{ %s }:", sel)
			line := fmt.Sprintf("%-42s [goldenrod]%s", l, attr)
			node = tview.NewTreeNode(line)
			node.
				SetColor(tcell.ColorCornflowerBlue).
				SetSelectable(true)

		case cue.ListKind:
			l := fmt.Sprintf("[ %s ]:", sel)
			node = tview.NewTreeNode(l)
			node.
				SetColor(tcell.ColorLime).
				SetSelectable(true)
			

		default:
			l := fmt.Sprintf("%s: %v %s", sel, value, attr)
			node = tview.NewTreeNode(l)

		}


		node.SetReference(fullpath)
		target.AddChild(node)
	}
}

func (VB *ValueBrowser) SetupKeybinds() {

	VB.SetInputCapture(func(evt *tcell.EventKey) *tcell.EventKey {

		switch evt.Key() {

		case tcell.KeyCtrlE:
			VB.App.Logger("eval!\n")

		case tcell.KeyCtrlD:
			VB.docs = !VB.docs

		case tcell.KeyCtrlA:
			VB.attrs = !VB.attrs

		case tcell.KeyCtrlO:
			VB.optional = !VB.optional

		case tcell.KeyCtrlJ:
			VB.inline = !VB.inline

		case tcell.KeyCtrlI:
			VB.ignore = !VB.ignore

		case tcell.KeyCtrlR:
			VB.resolve = !VB.resolve

		case tcell.KeyCtrlT:
			VB.concrete = !VB.concrete

		case tcell.KeyCtrlH:
			VB.hidden = !VB.hidden

		case tcell.KeyCtrlF:
			VB.final = !VB.final

		case tcell.KeyRune:
			switch evt.Rune() {
				case '#':
					VB.defs = !VB.defs
			default:
				return evt
			}

		default:
			return evt
		}

		VB.Rebuild("")

		return nil
	})

}
