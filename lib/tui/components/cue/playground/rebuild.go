package playground

import (
	"fmt"

	"cuelang.org/go/cue"
	"github.com/gdamore/tcell/v2"

	"github.com/hofstadter-io/hof/lib/singletons"
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
		C.scope.SetBorderColor(c)

	case "edit":
		C.edit.SetBorderColor(c)

	case "final":
		C.final.SetBorderColor(c)

	default:
		C.scope.SetBorderColor(c)
		C.edit.SetBorderColor(c)
		C.final.SetBorderColor(c)
	}
	go tui.Draw()
}


func (C *Playground) Rebuild() error {
	tui.Log("info", fmt.Sprintf("Play.Rebuild %v %v", C.useScope, C.scope.GetSourceConfigs()))

	var (
		v cue.Value
		err error
	)

	// just to be sure any children get updated
	C.UseScope(C.useScope)
	// show/hide scope as needed
	if C.seeScope {
		C.SetItem(0, C.scope, 0, 1, true)
	} else {
		C.SetItem(0, nil, 0, 0, false)
	}


	// user code that will be evaluated
	src := C.edit.GetText()
	C.setThinking(true, "final")

	go func() {

		defer C.setThinking(false, "final")

		// compile a value
		sv := C.scope.GetValue()
		if C.useScope && sv.Exists() {
			ctx := sv.Context()
			v = ctx.CompileString(src, cue.InferBuiltins(true), cue.Scope(sv))
		} else {
			// just compile the text
			ctx := singletons.CueContext()
			v = ctx.CompileString(src, cue.InferBuiltins(true))
		}

		var cfg *helpers.SourceConfig
		// make a new config with the latest value for the output
		cfg = &helpers.SourceConfig{Value: v}
		// only update view value, that way, if we erase everything, we still see the value


		if C.mode == ModeFlow {
			// first has to pass basic CUE checks so that errors look the same
			err = v.Validate()
			if err == nil {
				tui.Log("trace", "got here")
				// then we try to run the flow
				// we need a special way to deal with errors here
				v, err = C.runFlow(v)
				if err != nil {
					tui.Log("error", err)
					cfg = &helpers.SourceConfig{Text: err.Error()}
					C.final.SetMode("text")
				} else {
					cfg.Value = v
				}
			}

		}

		if err == nil && C.final.GetMode() == "text" {
			C.final.SetMode("cue")
		}

		C.final.ClearSourceConfigs()
		C.final.AddSourceConfig(cfg)
		C.final.RebuildValue()
		C.final.Rebuild()

		C.RebuildEditTitle()

		tui.Draw()
	}()
		
	tui.Draw()

	return nil
}

