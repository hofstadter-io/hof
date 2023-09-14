package browser

import (
	"fmt"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/format"
	"github.com/alecthomas/chroma/quick"
	"github.com/gdamore/tcell/v2"
	"github.com/kr/pretty"

	"github.com/hofstadter-io/hof/lib/cuetils"
	"github.com/hofstadter-io/hof/lib/gen"
	"github.com/hofstadter-io/hof/lib/singletons"
	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/components/cue/helpers"
	"github.com/hofstadter-io/hof/lib/tui/tview"
)

func (C *Browser) RebuildValue() {
	C.setThinking(true)
	defer C.setThinking(false)

	if len(C.sources) > 0 {
		// start with an empty value
		val := singletons.EmptyValue()

		// fill all the sources into one
		for _, S := range C.sources {
			v, _ := S.GetValue()
			val = val.FillPath(cue.ParsePath(S.Path), v)
		}

		// cache in the browser...?
		C.value = val
	}
}

func (C *Browser) setThinking(thinking bool) {
	c := tcell.ColorWhite
	if thinking {
		c = tcell.ColorViolet
	}

	C.SetBorderColor(c)
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

	if C.nextMode == "settings" {
		C.code.Clear()
		C.SetPrimitive(C.code)

		for _, s := range C.sources {
			fmt.Fprintf(C.codeW, "%# v\n\n", pretty.Formatter(*s))
		}
	} else if C.nextMode == "tree" {
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

		// possible load error
		wrote := false
		if err != nil {
			if err != nil {
				writeErr(err)
				wrote = true
			}
		}

		if !wrote {
			// first, possibly validate and possibly write an error
			if C.validate || C.nextMode == "json" || C.nextMode == "yaml" {
				err := C.value.Validate(C.Options()...)
				if err != nil {
					writeErr(err)
					wrote = true
				}
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
					// return
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

	// todo, show sources with a Frame and a hotkey?
	if len(VB.sources) > 0 {
		for _, src := range VB.sources {
			if src.Source != helpers.EvalNone && len(src.Args) > 0 {
				s += "[violet](" + strings.Join(src.Args, " ") + ")[-] "
			}
			if src.Source == helpers.EvalNone && VB.usingScope {
				s += "[blue]<S>[-] "
			}
		}
	} else {
		if VB.usingScope {
			s += "[blue]<S>[-] "
		} else {
			s += "[darkgray]<empty>[-] "
		}
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

