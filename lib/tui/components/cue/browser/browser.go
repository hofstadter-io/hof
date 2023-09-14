package browser

import (
	"fmt"
	"io"

	"cuelang.org/go/cue"
	"github.com/gdamore/tcell/v2"

	"github.com/hofstadter-io/hof/lib/singletons"
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

	// source configs
	sources []*helpers.SourceConfig

	// value from filled sources
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

func New() *Browser {
	C := &Browser {
		Frame: tview.NewFrame(),
		sources: make([]*helpers.SourceConfig,0),

		value: singletons.EmptyValue(),

		// some sane defaults
		mode: "cue",
		ignore: true,
		resolve: true,
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
		C.Frame.SetPrimitive(C.tree)
	} else {
		C.Frame.SetPrimitive(C.code)
	}

	C.SetBorder(true)
	C.SetupKeybinds()

	C.Frame.SetTitle(C.BuildStatusString())
	return C
}

func (VB *Browser) SetMode(mode string) {
	VB.mode = mode
}

func (VB *Browser) GetMode() string {
	return VB.mode
}

func (VB *Browser) GetText() string {
	return VB.code.GetText(true)
}

func (B *Browser) GetSourceConfigs() (sources []*helpers.SourceConfig) {
	return B.sources
}

func (B *Browser) AddSourceConfig(source *helpers.SourceConfig) {
	source.Name = fmt.Sprintf("%s-src.%d", B.Name(), len(B.sources))
	B.sources = append(B.sources, source)
}

func (B *Browser) SetSourceConfig(index int, source *helpers.SourceConfig) {
	source.Name = fmt.Sprintf("%s-src.%d", B.Name(), index)
	B.sources[index] = source
}

func (B *Browser) RemoveSourceConfig(index int) {
	B.sources = append(B.sources[:index], B.sources[index+1:]...)
}

func (B *Browser) ClearSourceConfigs() {
	B.sources = make([]*helpers.SourceConfig,0)
}

func (VB *Browser) GetUsingScope() bool {
	return VB.usingScope
}

func (VB *Browser) SetUsingScope(usingScope bool) {
	VB.usingScope = usingScope
}

func (C *Browser) SetWatchCallback(callback func()) {
	for _, s := range C.sources {
		s.WatchFunc = callback
	}
}

func (C *Browser) GetValue() cue.Value {
	return C.value
}

func (VB *Browser) GetValueExpr(expr string) func () cue.Value {
	// tui.Log("trace", fmt.Sprintf("View.GetConnValueExpr from: %s/%s %s", VB.Id(), VB.Name(), expr))
	p := cue.ParsePath(expr)

	return func() cue.Value {
		return VB.GetValue().LookupPath(p)
	}

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

			// info about CUE value? (stats)
			//case 'I':
			//  VB.nextMode = "info"

			// show settings, hidden?

				VB.nextMode = "settings"

			// todo, dive values, and walk back up
			//case 'I': // in
			//case 'U': // up

			case 'R':
				VB.RebuildValue()

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
