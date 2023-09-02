package components

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/format"
	"github.com/alecthomas/chroma/quick"
	"github.com/gdamore/tcell/v2"

	"github.com/hofstadter-io/hof/lib/cuetils"
	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/tview"
)

type ValueBrowser struct {
	*tview.Frame
	Tree *tview.TreeView
	Code *tview.TextView
	CodeW io.Writer

	OnFieldSelect func(string)

	Root *tview.TreeNode

	Value cue.Value

	// if root is expanded or not
	expanded bool

	// mode mode [tree,cue,yaml,json]
	mode string
	nextMode string

	// eval settings
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

func NewValueBrowser(val cue.Value, mode string, OnFieldSelect func(path string)) *ValueBrowser {
	VB := &ValueBrowser {
		Value: val,
		mode: mode,
	}

	// code view
	VB.Code = tview.NewTextView()
	VB.CodeW = tview.ANSIWriter(VB.Code)
	VB.Code.SetWordWrap(true).
		SetDynamicColors(true)

	// tree view
	VB.Root = tview.NewTreeNode("no results yet")
	VB.Root.SetColor(tcell.ColorSilver)

	VB.Tree = tview.NewTreeView()
	VB.Tree. SetRoot(VB.Root).SetCurrentNode(VB.Root)

	// set our selected handler for tree
	VB.Tree.SetSelectedFunc(VB.OnSelect)


	if VB.mode == "tree" {
		VB.Frame = tview.NewFrame(VB.Tree)
	} else {
		VB.Frame = tview.NewFrame(VB.Code)
	}

	VB.SetBorder(true)
	VB.SetupKeybinds()

	return VB
}

func (VB *ValueBrowser) Rebuild(path string) {
	if path == "" {
		path = "<root>"
	}

	if VB.nextMode == "" {
		VB.nextMode = VB.mode
	}

	if VB.nextMode == "tree" {

		root := tview.NewTreeNode(path)
		root.SetColor(tcell.ColorSilver)
		tree := tview.NewTreeView()

		VB.AddAt(root, path)
		tree.SetRoot(root).SetCurrentNode(root)
		tree.SetSelectedFunc(VB.OnSelect)
		tree.SetDoubleClickedFunc(func(node *tview.TreeNode) {
			// double clicking a node impacts all children
			if node.GetLevel() == 0 {
				VB.expanded = !VB.expanded
				if VB.expanded {
					node.ExpandAll()
				} else {
					node.CollapseAll()
					node.Expand()
				}
			}  else {
				if node.IsExpanded() {
					node.CollapseAll()
				} else {
					node.ExpandAll()
				}
			}

		})

		VB.SetPrimitive(tree)

		// TODO, dual-walk old-new tree's too keep things open
		VB.Tree = tree
		VB.Root = root

	} else {
		VB.Code.Clear()
		syn := VB.Value.Syntax(VB.Options()...)

		b, err := format.Node(syn)
		if !VB.ignore {
			if err != nil {
				s := cuetils.CueErrorToString(err)
				fmt.Fprintln(VB.CodeW, s)
			}
		}

		err = quick.Highlight(VB.CodeW, string(b), "Go", "terminal256", "solarized-dark")
		if err != nil {
			go tui.SendCustomEvent("/console/error", fmt.Sprintf("error highlighing %v", err))
			return
		}

		VB.SetPrimitive(VB.Code)

	}

	if VB.nextMode != VB.mode {
		VB.mode = VB.nextMode
		VB.Focus(func(p tview.Primitive){
			p.Focus(nil)
		})
	}
	VB.nextMode = ""

	VB.Frame.SetTitle(VB.buildStatusString())

}
func (VB *ValueBrowser) buildStatusString() string {

	var s string

	add := func(on bool, char string) {
		if on {
			s += "[lime]" + char + "[-]"
		} else {
			s += char
		}
	}

	s += VB.mode + " ["

	add(VB.mode == "tree", "T")
	add(VB.mode == "cue", "C")
	add(VB.mode == "yaml", "Y")
	add(VB.mode == "json", "J")

	s += "] "

	add(VB.concrete, "c")
	add(VB.final, "f")
	add(VB.resolve, "r")

	s += " "

	add(VB.ignore, "e")
	add(VB.inline, "i")

	s += " "

	add(VB.defs, "d")
	add(VB.optional, "o")
	add(VB.hidden, "h")

	s += " "

	add(VB.docs, "D")
	add(VB.attrs, "A")


	// add some buffer
	s = "  " + s + "  "
	return s
}

func (VB *ValueBrowser) OnSelect(node *tview.TreeNode) {
	reference := node.GetReference()
	if reference == nil {
		return // Selecting the root node does nothing.
	}

	children := node.GetChildren()
	if len(children) == 0 {
		// Load and show files in this directory.
		path := reference.(string)
		VB.AddAt(node, path)
	} else {
		// Collapse if visible, expand if collapsed.
		node.SetExpanded(!node.IsExpanded())
	}
}


func (VB *ValueBrowser) AddAt(target *tview.TreeNode, path string) {
	// VB.App.Logger(fmt.Sprintf("VB.AddAt: %s\n", path))

	if strings.HasPrefix(path, "<root>") {
		path = ""
	}
	if strings.HasPrefix(path, ".") {
		path = path[1:]
	}
	val := VB.Value.LookupPath(cue.ParsePath(path))
	// VB.App.Logger(fmt.Sprintf("#v\n", val))

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
		i, _ := val.List()
		iter = &i
	}
	if iter == nil {
		// VB.App.Logger(fmt.Sprintf("nil iter for: %s\n", path))
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

		}


		node.SetReference(fullpath)
		target.AddChild(node)
	}
}

func (VB *ValueBrowser) SetMode(mode string) {
	VB.mode = mode
}

func (VB *ValueBrowser) GetMode() string {
	return VB.mode
}

func (VB *ValueBrowser) Options() []cue.Option {
		opts := []cue.Option{
			cue.ResolveReferences(VB.resolve),
			cue.InlineImports(VB.inline),
			cue.ErrorsAsValues(VB.ignore),
			cue.Docs(VB.docs),
			cue.Attributes(VB.attrs),
			cue.Optional(VB.optional),
			cue.Definitions(VB.defs),
		}
		if VB.concrete {
			opts = append(opts, cue.Concrete(true))
		}
		if VB.hidden {
			opts = append(opts, cue.Hidden(true))
		}

		if VB.final {
			// prepend final, so others still apply
			opts = append([]cue.Option{cue.Final()}, opts...)
		}

	return opts
}

func (VB *ValueBrowser) SetupKeybinds() {

	VB.SetInputCapture(func(evt *tcell.EventKey) *tcell.EventKey {

		if evt.Key() == tcell.KeyRune {
			switch evt.Rune() {

			case 'c':
				VB.concrete = !VB.concrete
			case 'f':
				VB.final = !VB.final
			case 'r':
				VB.resolve = !VB.resolve

			case 'i':
				VB.inline = !VB.inline
			case 'e':
				VB.ignore = !VB.ignore

			case 'd':
				VB.defs = !VB.defs
			case 'o':
				VB.optional = !VB.optional
			case 'h':
				VB.hidden = !VB.hidden

			case 'D':
				VB.docs = !VB.docs
			case 'A':
				VB.attrs = !VB.attrs

			case 'Y':
				VB.nextMode = "yaml"
			case 'J':
				VB.nextMode = "json"
			case 'C':
				VB.nextMode = "cue"
			case 'T':
				VB.nextMode = "tree"

			// todo, dive values, and walk back up
			case 'I': // in
			case 'U': // up

			default:
				return evt
			}

			VB.Rebuild("")

			return nil
		}

		return evt
	})	

}

func (VB *ValueBrowser) Focus(delegate func(p tview.Primitive)) {
	switch VB.mode {
	case "tree":
		delegate(VB.Tree)
	default:
		delegate(VB.Code)
	}
}

//func (VB *ValueBrowser) Mount(context map[string]any) error {


//  return nil
//}

