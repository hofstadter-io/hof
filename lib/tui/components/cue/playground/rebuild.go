package playground

import (
	"fmt"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"github.com/atotto/clipboard"

	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/components/cue/helpers"
)

func (C *Playground) HandleAction(action string, args []string, context map[string]any) (bool, error) {
	tui.Log("warn", fmt.Sprintf("Playground.HandleAction: %v %v %v", action, args, context))
	var err error
	handled := true

	// item actions
	switch action {
	case "push":
		id, err := C.PushToPlayground()
		// if ok...
		if err == nil {
			msg := fmt.Sprintf("snippet id: %s  (link copied!)", id)

			url := fmt.Sprintf("https://cuelang.org/play?id=%s", id)
			clipboard.WriteAll(url)

			tui.Tell("error", msg)
			tui.Log("trace", msg)
		}


	case "write":
		if len(args) != 1 {
			err = fmt.Errorf("write requires a filename")
		} else {
			filename := args[0]
			err = C.WriteEditToFile(filename)
			// if ok...
			if err == nil {
				msg := fmt.Sprintf("editor text saved to %s", filename)
				tui.Tell("error", msg)
				tui.Log("trace", msg)
			}
		}

	case "export":
		if len(args) != 1 {
			err = fmt.Errorf("export requires a filename")
		} else {
			filename := args[0]
			err := C.ExportFinalToFile(filename)
			// if ok...
			if err == nil {
				msg := fmt.Sprintf("value exported to %s", filename)
				tui.Tell("error", msg)
				tui.Log("trace", msg)
			}
		}

	case "update":
		err = C.updateFromArgsAndContext(args, context)

	default:
		err = fmt.Errorf("unknown command %q", action)
	}

	return handled, err
}

func (C *Playground) updateFromArgsAndContext(args[] string, context map[string]any) error {
	// get source, defaults to empty, new runtime?
	source := ""
	if _source, ok := context["source"]; ok {
		source = _source.(string)
	}

	target := "value"
	if _target, ok := context["target"]; ok {
		target = _target.(string)
	}

	// setup our source config
	srcCfg := helpers.SourceConfig{
		Args: args,
	}

	// special case, source will be empty when the args are all cue entrypoints
	// we want to...
	//   (1) catch special empty case for new play
	//   (2) we want different defaults for empty when there are args, based on the target
	//   for (1), we need temporary <new-play> to know we are in new play mode
	if len(args) == 0 || (len(args) == 1 && args[0] == "new") {
		source = "<new-play>" // very temporary setting
		target = "value"
	}

	rebuildScope := false
	switch target {
	case "value":
		// local source default, assume it was a filename
		if source == "" {
			source = "file"
		} else if source == "<new-play>" {
			source = ""
		}
		srcCfg.Source = helpers.EvalSource(source)

		// tui.Log("extra", fmt.Sprintf("Eval.playItem.value: %v", srcCfg ))
		C.UseScope(false)
		// need to get the text once at startup
		txt, err := srcCfg.GetText()
		if err != nil {
			tui.Log("error", err)
			return err
		}
		// tui.Log("extra", fmt.Sprintf("Eval.playItem.value.text: %v", txt ))
		C.SetText(txt)

	case "scope":
		if source == "" {
			source = "runtime"
		}
		srcCfg.Source = helpers.EvalSource(source)

		C.SetScopeConfig(srcCfg)

		rebuildScope = true
		C.UseScope(true)
	}

	return C.Rebuild(rebuildScope)
}


func (C *Playground) Rebuild(rebuildScope bool) error {
	tui.Log("info", fmt.Sprintf("Play.rebuildScope %v %v %v", rebuildScope, C.useScope, C.scope.config))
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

