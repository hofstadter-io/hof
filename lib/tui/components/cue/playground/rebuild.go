package playground

import (
	"fmt"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"github.com/gdamore/tcell/v2"

	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/components/cue/helpers"
)

func (C *Playground) setThinking(thinking bool, which string) {
	c := tcell.ColorWhite
	if thinking {
		c = tcell.ColorViolet
	}

	switch which {
	case "scope":
		C.scope.viewer.SetBorderColor(c)

	case "final":
		C.final.viewer.SetBorderColor(c)

	default:
		C.scope.viewer.SetBorderColor(c)
		C.edit.SetBorderColor(c)
		C.final.viewer.SetBorderColor(c)
	}
	go tui.Draw()
}


func (C *Playground) Rebuild(rebuildScope bool) error {
	tui.Log("info", fmt.Sprintf("Play.rebuildScope %v %v %v", rebuildScope, C.useScope, C.scope.config))
	// fmt.Println("got here")
	var (
		v cue.Value
		err error
	)

	// just to be sure any children get updated
	C.UseScope(C.useScope)

	ctx := cuecontext.New()
	src := C.edit.GetText()

	C.setThinking(true, "final")
	defer C.setThinking(false, "final")

	// compile a value
	if !C.useScope {
		// just compile the text
		v = ctx.CompileString(src, cue.InferBuiltins(true))
	} else {
		// compile the text with a scope

		sv, serr := C.scope.config.GetValue()
		err = serr

		if err != nil {
			tui.Log("error", err)
		}
		// we shouldn't have to worry about this, but we aren't catching all the ways
		// that we get into this code, in particular, hotkey can set scope to true when none exists
		if !sv.Exists() {
			tui.Log("error", "scope value does not exist")
			err = fmt.Errorf("scope value does not exist")
		}

		if err == nil && sv.Exists() {
			if rebuildScope {
				C.scope.viewer.Rebuild()
			}

			// tui.Log("warn", fmt.Sprintf("recompile with scope: %v", rebuildScope))
			ctx := sv.Context()
			v = ctx.CompileString(src, cue.InferBuiltins(true), cue.Scope(sv))
		}
	}


	cfg := &helpers.SourceConfig{Value: v}
	if err != nil {
		tui.Log("error", err)
		cfg = &helpers.SourceConfig{Text: err.Error()}
	}

	// only update view value, that way, if we erase everything, we still see the value
	C.final.config = cfg
	C.final.viewer.SetSourceConfig(cfg)
	C.final.viewer.SetUsingScope(C.useScope)
	C.final.viewer.Rebuild()

	// show/hide scope as needed
	if C.seeScope {
		C.SetItem(0, C.scope.viewer, 0, 1, true)
	} else {
		C.SetItem(0, nil, 0, 0, false)
	}

	// tui.Draw()
	return nil
}

