package browser

import (
	"io"

	"cuelang.org/go/cue"
	"github.com/gdamore/tcell/v2"

	"github.com/hofstadter-io/hof/lib/tui/components/cue/helpers"
	"github.com/hofstadter-io/hof/lib/tui/tview"
)

type Browser struct {
	*tview.Frame

	// mode mode [tree,cue,yaml,json]
	mode string
	nextMode string
	refocus bool  // possibly refocus, if we rebuild the tree or switch views
	usingScope bool // this is just for display in the status, scope is not used here, but impacts the results from the playground

	// tree view
	tree *tview.TreeView
	root *tview.TreeNode
	expanded bool // if root is expanded or not

	// code view
	code *tview.TextView
	codeW io.Writer

	// source config & value
	source helpers.SourceConfig
	value cue.Value

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

func (*Browser) TypeName() string {
	return "cue/browser"
}

func New(source helpers.SourceConfig, mode string) *Browser {
	C := &Browser {
		source: source,
		mode: mode,
	}

	// code view
	C.code = tview.NewTextView()
	C.codeW = tview.ANSIWriter(C.code)
	C.code.SetWordWrap(true).
		SetDynamicColors(true)

	// tree view
	C.root = tview.NewTreeNode("no results yet")
	C.root.SetColor(tcell.ColorSilver)

	C.tree = tview.NewTreeView()
	C.tree. SetRoot(C.root).SetCurrentNode(C.root)

	// set our selected handler for tree
	C.tree.SetSelectedFunc(C.onSelect)


	if C.mode == "tree" {
		C.Frame = tview.NewFrame(C.tree)
	} else {
		C.Frame = tview.NewFrame(C.code)
	}

	C.SetBorder(true)
	C.SetupKeybinds()

	return C
}

func (VB *Browser) SetMode(mode string) {
	VB.mode = mode
}

func (VB *Browser) GetMode() string {
	return VB.mode
}

func (VB *Browser) SetUsingScope(usingScope bool) {
	VB.usingScope = usingScope
}

func (B *Browser) SetSourceConfig(source helpers.SourceConfig) {
	B.source = source
}

func (VB *Browser) GetUsingScope() bool {
	return VB.usingScope
}

func (VB *Browser) Options() []cue.Option {
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
		// PRE-PEND Final, so others still apply:
		opts = append([]cue.Option{cue.Final()}, opts...)
	}

	return opts
}

func (VB *Browser) SetupKeybinds() {

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
			case 't':
				VB.nextMode = "text"

			// todo, dive values, and walk back up
			//case 'I': // in
			//case 'U': // up

			case 'R':
				// this should just trigger a refresh by nature of being here

			default:
				return evt
			}

			VB.refocus = true
			VB.Rebuild()

			return nil
		}

		return evt
	})	

}

func (VB *Browser) Focus(delegate func(p tview.Primitive)) {
	switch VB.mode {
	case "tree":
		delegate(VB.tree)
	default:
		delegate(VB.code)
	}
}
