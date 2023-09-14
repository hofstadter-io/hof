package playground

import (
	"fmt"
	"strings"
	"time"

	"github.com/atotto/clipboard"

	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/components/cue/helpers"
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

	case "update":
		err = C.updateFromArgsAndContext(args, context)

	case "set.value":
		C.setThinking(true, "final")
		defer C.setThinking(false, "final")
		context["target"] = "value"
		err = C.updateFromArgsAndContext(args, context)
	case "set.scope":
		C.setThinking(true, "scope")
		defer C.setThinking(false, "scope")
		context["target"] = "scope"
		err = C.updateFromArgsAndContext(args, context)

	case "get.refresh":
	  tui.StatusMessage(fmt.Sprintf("refresh debounce time is %s", C.debounceTime))
	case "refresh", "set.refresh":
		if len(args) != 1 {
			err = fmt.Errorf("refresh requires a duration like 1s or 300ms")
			tui.Tell("error", err)
		} else {
			d, derr := time.ParseDuration(args[0])
			if derr != nil {
				err = derr
				tui.Tell("error", err)
			} else {
				C.debounceTime = d
				if d.Nanoseconds() <= 0 {
					// disabled
					C.edit.SetChangedFunc(nil)
				} else {
					C.debouncer = watch.NewDebouncer(C.debounceTime)
					// enabled
					C.edit.SetChangedFunc(func() {
						C.debouncer(func(){
							C.Rebuild(false)
						})
					})
				}
			}
		}

	case "get.scope.watch":
	  tui.StatusMessage(fmt.Sprintf("watch debounce time is %s", C.scope.config.WatchTime))
	case "get.value.watch":
	  tui.StatusMessage(fmt.Sprintf("watch debounce time is %s", C.editCfg.WatchTime))
	case "set.scope.watchGlobs", "watchGlobs":
		C.scope.config.WatchGlobs = args
	case "set.value.watchGlobs":
		C.editCfg.WatchGlobs = args
	case "watch", "set.scope.watch", "set.value.watch":
		d := time.Duration(42*time.Millisecond)
		if len(args) < 1 {
			//aerr := fmt.Errorf("watch requires a duration like 1s or 300ms")
			//tui.Log("warn", aerr)
			tui.Log("warn", fmt.Sprintf("no watch duration given, setting to %s", d))
			tui.StatusMessage(fmt.Sprintf("no watch duration given, setting to %s", d))
		} else {
			d, err = time.ParseDuration(args[0])
		}

		if err != nil {
			tui.Tell("error", err)
		} else {

			do := func(cfg *helpers.SourceConfig, callback func()) {
				// tui.Log("trace", fmt.Sprintf("DOOOOOO: %v", cfg))
				cfg.WatchTime = d
				// tui.Log("trace", fmt.Sprintf("watch config: %v", cfg))

				if d.Nanoseconds() > 0 {
					// startup new watch
					tui.StatusMessage(fmt.Sprintf("start %sing... %v", strings.TrimPrefix(action, "set."), cfg))
					err = cfg.Watch(C.Name(), callback, d)
					// tui.Log("crit", fmt.Sprintf("watch err: %v", err))
				} else {
					// or stop any watches
					cfg.StopWatch()
				}
			}

			// set scope watch
			if action == "set.scope.watch" || action == "watch" {
				// tui.Log("crit", fmt.Sprintf("got over here.0"))
				cfg := C.scope.config
				C.scope.viewer.SetSourceConfig(cfg)
				// tui.Log("crit", fmt.Sprintf("got over here: %v", cfg))
				go do(cfg, func() {
					// tui.Log("crit", fmt.Sprintf("got over here again: %v", cfg))
					C.Rebuild(true)
				})
			}

			// set value watch
			if action == "set.value.watch" || action == "watch" {
				cfg := C.editCfg
				go do(cfg, func() {
					txt, err := cfg.GetText()
					if err != nil {
						tui.Log("error", err)
					}
					C.SetText(txt)
					// when action, we only want the scope full rebuild to happen
					if action != "watch" && C.scope.config.WatchQuit != nil {
						// tui.Log("crit", fmt.Sprintf("got here again.2: %v", C.scope.config))
						C.Rebuild(false)
					}
				})
			} 

		}

	default:
		handled = false
		// err = fmt.Errorf("unknown command %q", action)
	}

	return handled, err
}

func (C *Playground) updateFromArgsAndContext(args[] string, context map[string]any) error {
	// tui.Log("warn", fmt.Sprintf("Playground.updateHandler.1: %v %v", args, context))
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

	rebuildScope := false
	switch target {
	case "value":
		// local source default, assume it was a filename
		if source == "" {
			source = "file"
		} else if source == "<new-play>" {
			source = ""
		}
		C.editCfg.Source = helpers.EvalSource(source)
		C.editCfg.Args = args

		// tui.Log("warn", fmt.Sprintf("Playground.updateHandler.3.cfg: %v", C.editCfg))

		txt, err := C.editCfg.GetText()
		if err != nil {
			tui.Log("error", err)
			return err
		}
		// tui.Log("extra", fmt.Sprintf("Playground.updateHandler.4.text: %v", txt ))
		C.SetText(txt)

	case "scope":
		if source == "" {
			source = "runtime"
		}
		srcCfg := C.scope.config
		srcCfg.Source = helpers.EvalSource(source)
		srcCfg.Args = args

		// tui.Log("warn", fmt.Sprintf("Playground.updateHandler.3.S: %v", srcCfg))
		C.SetScopeConfig(srcCfg)

		rebuildScope = true
		C.seeScope = true
		C.UseScope(true)
	}

	return C.Rebuild(rebuildScope)
}

