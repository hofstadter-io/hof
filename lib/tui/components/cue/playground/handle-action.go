package playground

import (
	"fmt"
	"time"

	"github.com/atotto/clipboard"

	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/watch"
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

	case "create":
		err = C.updateFromArgsAndContext(args, context)

	case "update":
		err = C.updateFromArgsAndContext(args, context)

	case "set.value":
		C.setThinking(true, "final")
		defer C.setThinking(false, "final")
		context["target"] = "value"
		err = C.updateFromArgsAndContext(args, context)

	case "refresh":
		if len(args) < 1 {
			tui.StatusMessage(fmt.Sprintf("refresh debounce time is %s", C.debounceTime))
		} else {
			C.SetDebounceTime(args, context)
		}
	case "watch":
		C.startScopeWatch(args, context)

	case "rebuild.scope":
		C.scope.RebuildValue()
		C.scope.Rebuild()
		C.Rebuild()

	case "flow.clear":
		C.mode = ModeEval
		C.Rebuild()
	case "flow.run":
		C.mode = ModeFlow
		C.Rebuild()

	default:
		handled = false
		// err = fmt.Errorf("unknown command %q", action)
	}

	if !handled {
		// try handling in the scope browser
		handled, err = C.scope.HandleAction(action, args, context)

		// hmm, there are cases where it would be nice to set or see scope
		// let's just do it by default for now?
		C.UseScope(true)
		C.Rebuild()
	}

	return handled, err
}

func (C *Playground) startScopeWatch(args []string, context map[string]any) {
	C.scope.SetWatchCallback(func() {
		C.scope.RebuildValue()
		C.scope.Rebuild()
		C.Rebuild()
	})
	C.scope.HandleAction("watch", args, context)
}

func (C *Playground) startValueWatch(args []string, context map[string]any) {
	cfg := C.editCfg
	cfg.WatchFunc = func() {
		txt, err := cfg.GetText()
		if err != nil {
			tui.Log("error", err)
			txt = err.Error()
		}
		C.edit.SetText(txt, false)
		C.Rebuild()
	}
	go cfg.Watch()
}

func (C *Playground) updateFromArgsAndContext(args[] string, context map[string]any) error {
	tui.Log("warn", fmt.Sprintf("Playground.updateHandler.1: %v %v", args, context))
	// get source, defaults to empty, new runtime?
	source := ""
	if _source, ok := context["source"]; ok {
		source = _source.(string)
	}

	target := "value"
	if _target, ok := context["target"]; ok {
		target = _target.(string)
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

	// tui.Log("warn", fmt.Sprintf("Playground.updateHandler.2: %v %v %v", source, target, srcCfg))

	switch target {
	case "value":
		// local source default, assume it was a filename
		if source == "" {
			context["source"] = "file"
		} else if source == "<new-play>" {
			context["source"] = ""
		}
		C.editCfg.UpdateFrom("set", args, context)

		// tui.Log("warn", fmt.Sprintf("Playground.updateHandler.3.cfg: %v", C.editCfg))
		{
			C.setThinking(true, "edit")
			defer C.setThinking(false, "edit")

			txt, err := C.editCfg.GetText()
			if err != nil {
				tui.Log("error", err)
				return err
			}

			// tui.Log("extra", fmt.Sprintf("Playground.updateHandler.4.text: %v", txt ))
			C.edit.SetText(txt, false)
		}

	case "scope":
		if source == "" {
			context["source"] = "runtime"
		}
		C.scope.HandleAction("create", args, context)
		C.seeScope = true
		C.UseScope(true)

		C.scope.RebuildValue()
		C.scope.Rebuild()
	}

	return C.Rebuild()
}


func (C *Playground) SetDebounceTime(args []string, context map[string]any) error {
	d, err := time.ParseDuration(args[0])
	if err != nil {
		return err
	}
	C.debounceTime = d

	if d.Nanoseconds() <= 0 {
		// disabled
		C.edit.SetChangedFunc(nil)
	} else {
		C.debouncer = watch.NewDebouncer(C.debounceTime)
		// enabled
		C.edit.SetChangedFunc(func() {
			C.debouncer(func(){
				C.Rebuild()
			})
		})
	}

	return nil
}
