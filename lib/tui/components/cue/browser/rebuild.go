package browser

import (
	"fmt"
	"strings"

	"cuelang.org/go/cue/format"
	"github.com/alecthomas/chroma/quick"
	"github.com/gdamore/tcell/v2"

	"github.com/hofstadter-io/hof/lib/cuetils"
	"github.com/hofstadter-io/hof/lib/gen"
	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/tview"
)

func (C *Browser) setThinking(thinking bool) {
	c := tcell.ColorWhite
	if thinking {
		c = tcell.ColorViolet
	}

	C.tree.SetBorderColor(c)
	C.code.SetBorderColor(c)
	go tui.Draw()
}

func (C *Browser) Rebuild() {
	var err error

	C.setThinking(true)
	defer C.setThinking(false)

	path := "<root>"

	if C.nextMode == "" {
		C.nextMode = C.mode
	}

	writeErr := func(err error) {
		C.code.Clear()
		fmt.Fprint(C.codeW, cuetils.CueErrorToString(err))
		C.SetPrimitive(C.code)
	}

	C.value, err = C.source.GetValue()
	if err != nil {
		writeErr(err)
		return
	}

	// tui.Log("info", fmt.Sprintf("View.Rebuild %v %v", C.usingScope, C.nextMode))
	// tui.Log("info", fmt.Sprintf("View.Rebuild %v %v %v", C.usingScope, C.nextMode, C.value))
	// todo, rebuild value from source config
	// for now, we are still doing this outside

	// special case tree first
	if C.nextMode == "tree" {
		root := tview.NewTreeNode(path)
		root.SetColor(tcell.ColorSilver)
		tree := tview.NewTreeView()

		C.AddAt(root, path)
		tree.SetRoot(root).SetCurrentNode(root)
		tree.SetSelectedFunc(C.onSelect)

		C.SetPrimitive(tree)

		// TODO, dual-walk old-new tree's too keep things open
		C.tree = tree
		C.root = root

	} else {

		// otherwise, we need to turn the value into a string for code browsing
		C.code.Clear()
		C.SetPrimitive(C.code)

		// first, possibly validate and possibly write an error
		wrote := false
		if C.validate || C.nextMode == "json" || C.nextMode == "yaml" {
			err := C.value.Validate(C.Options()...)
			if err != nil {
				writeErr(err)
				wrote = true
			}
		}

		if !wrote {
			var (
				b []byte
				err error
			)
			switch C.nextMode {
			case "cue":
				syn := C.value.Syntax(C.Options()...)

				b, err = format.Node(syn)
				if !C.ignore {
					if err != nil {
						writeErr(err)
					}
				}

			case "json":
				f := &gen.File{}
				b, err = f.FormatData(C.value, "json")
				if err != nil {
					writeErr(err)
				}

			case "yaml":
				f := &gen.File{}
				b, err = f.FormatData(C.value, "yaml")
				if err != nil {
					writeErr(err)
				}

			}

			if err == nil {
				err = quick.Highlight(C.codeW, string(b), "Go", "terminal256", "solarized-dark")
				// tui.Log("info", fmt.Sprintf("View.Rebuild writing..."))
				if err != nil {
					writeErr(err)
					tui.Log("crit", fmt.Sprintf("error highlighing %v", err))
					return
				}
			}
		}


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
	C.Frame.SetTitle(C.BuildStatusString())
	// tui.Draw()
}

func (VB *Browser) BuildStatusString() string {

	var s string

	if n := VB.Name(); len(n) > 0 {
		s += n + ": "
	}
	if len(VB.source.Args) > 0 {
		s += "[violet](" + strings.Join(VB.source.Args, " ") + ")[-] "
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
	add(VB.usingScope, " S")
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

