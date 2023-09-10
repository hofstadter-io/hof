package playground

import (
	"fmt"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"

	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/components/cue/helpers"
)

func (C *Playground) Rebuild(rebuildScope bool) error {
	// tui.Log("info", fmt.Sprintf("Play.rebuildScope %v %v %v", rebuildScope, C.useScope, C.scope.config))
	var (
		v cue.Value
		err error
	)

	ctx := cuecontext.New()
	src := C.edit.GetText()

	// compile a value
	if !C.useScope {
		// just compile the text
		v = ctx.CompileString(src, cue.InferBuiltins(true))
	} else {
		// compile the text with a scope

		// tui.Log("warn", fmt.Sprintf("%#v", s))
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
				// C.scope.config.Rebuild()
				cfg := helpers.SourceConfig{Value: sv}
				C.scope.viewer.SetSourceConfig(cfg)
				C.scope.viewer.Rebuild()
			}

			// tui.Log("warn", fmt.Sprintf("recompile with scope: %v", rebuildScope))
			ctx := sv.Context()
			v = ctx.CompileString(src, cue.InferBuiltins(true), cue.Scope(sv))
		}
	}

	cfg := helpers.SourceConfig{Value: v}
	if err != nil {
		tui.Log("error", err)
		cfg = helpers.SourceConfig{Text: err.Error()}
	}
	// only update view value, that way, if we erase everything, we still see the value
	C.final.viewer.SetUsingScope(C.useScope)
	C.final.viewer.SetSourceConfig(cfg)
	C.final.viewer.Rebuild()

	// show/hide scope as needed
	if C.useScope {
		C.SetItem(0, C.scope.viewer, 0, 1, true)
	} else {
		C.SetItem(0, nil, 0, 0, false)
	}


	// tui.Draw()
	return nil
}

