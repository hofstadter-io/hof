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

	Root *tview.TreeNode

	Value cue.Value

	// if set, gets the value on refresh
	ValueGetter func() cue.Value

	// if root is expanded or not
	expanded bool

	// mode mode [tree,cue,yaml,json]
	mode string
	nextMode string
	refocus bool  // possibly refocus, if we rebuild the tree or switch views

	UsingScope bool

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
	final,
	validate bool
}

func (*ValueBrowser) TypeName() string {
	return "ValueBrowser"
}
func (V *ValueBrowser) EncodeMap() (map[string]any, error) {
	return map[string]any{
		"type": V.TypeName(),
		"mode": V.mode,
		"useScope": V.UsingScope,
		"docs": V.docs,
		"attrs": V.attrs,
		"defs": V.defs,
		"optional": V.optional,
		"ignore": V.ignore,
		"inline": V.inline,
		"resolve": V.resolve,
		"concrete": V.concrete,
		"hidden": V.hidden,
		"final": V.final,
		"validate": V.validate,
	}, nil
}

func NewValueBrowser(val cue.Value, mode string) *ValueBrowser {
	C := &ValueBrowser {
		Value: val,
		mode: mode,
	}

	// code view
	C.Code = tview.NewTextView()
	C.CodeW = tview.ANSIWriter(C.Code)
	C.Code.SetWordWrap(true).
		SetDynamicColors(true)

	// tree view
	C.Root = tview.NewTreeNode("no results yet")
	C.Root.SetColor(tcell.ColorSilver)

	C.Tree = tview.NewTreeView()
	C.Tree. SetRoot(C.Root).SetCurrentNode(C.Root)

	// set our selected handler for tree
	C.Tree.SetSelectedFunc(C.OnSelect)


	if C.mode == "tree" {
		C.Frame = tview.NewFrame(C.Tree)
	} else {
		C.Frame = tview.NewFrame(C.Code)
	}

	C.SetBorder(true)
	C.SetupKeybinds()

	return C
}

func (C *ValueBrowser) Rebuild(path string) {
	if C.ValueGetter != nil {
		C.Value = C.ValueGetter()
	}

	if path == "" {
		path = "<root>"
	}

	if C.nextMode == "" {
		C.nextMode = C.mode
	}

	if C.nextMode == "tree" {

		root := tview.NewTreeNode(path)
		root.SetColor(tcell.ColorSilver)
		tree := tview.NewTreeView()

		C.AddAt(root, path)
		tree.SetRoot(root).SetCurrentNode(root)
		tree.SetSelectedFunc(C.OnSelect)
		//tree.SetDoubleClickedFunc(func(node *tview.TreeNode) {
		//  // double clicking a node impacts all children
		//  if node.GetLevel() == 0 {
		//    C.expanded = !C.expanded
		//    if C.expanded {
		//      node.ExpandAll()
		//    } else {
		//      node.CollapseAll()
		//      node.Expand()
		//    }
		//  }  else {
		//    if node.IsExpanded() {
		//      node.CollapseAll()
		//    } else {
		//      node.ExpandAll()
		//    }
		//  }

		//})

		C.SetPrimitive(tree)

		// TODO, dual-walk old-new tree's too keep things open
		C.Tree = tree
		C.Root = root

	} else {
		C.Code.Clear()
		wrote := false
		if C.validate {
			err := C.Value.Validate(C.Options()...)
			if err != nil {
				fmt.Fprint(C.CodeW, cuetils.CueErrorToString(err))
				wrote = true
			}
		}


		if !wrote {
			syn := C.Value.Syntax(C.Options()...)

			b, err := format.Node(syn)
			if !C.ignore {
				if err != nil {
					s := cuetils.CueErrorToString(err)
					fmt.Fprintln(C.CodeW, s)
				}
			}

			err = quick.Highlight(C.CodeW, string(b), "Go", "terminal256", "solarized-dark")
			if err != nil {
				go tui.SendCustomEvent("/console/error", fmt.Sprintf("error highlighing %v", err))
				return
			}
		}

		C.SetPrimitive(C.Code)

	}

	if C.refocus {
		C.refocus = false
		if C.nextMode != C.mode {
			C.mode = C.nextMode
		}
		C.Focus(func(p tview.Primitive){
			p.Focus(nil)
		})
	}

	C.nextMode = ""
	C.Frame.SetTitle(C.buildStatusString())

}
func (VB *ValueBrowser) buildStatusString() string {

	var s string

	if n := VB.Name(); len(n) > 0 {
		s += n + " -  "
	}

	add := func(on bool, char string) {
		if on {
			s += "[lime]" + char + "[-]"
		} else {
			s += char
		}
	}

	s += VB.mode + " ["
	add(VB.mode == "tree", "T")
	add(VB.mode == "cue",  "C")
	add(VB.mode == "json", "J")
	add(VB.mode == "yaml", "Y")
	add(VB.UsingScope, " S")
	s += "] "

	add(VB.validate, "v")
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

	// add some space around the final result
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
		//path := reference.(string)
		//VB.AddAt(node, path)
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

			case 'v':
				VB.validate = !VB.validate
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

			VB.refocus = true
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

