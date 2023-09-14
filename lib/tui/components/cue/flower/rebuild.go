package flower

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cuelang.org/go/cue"
	"github.com/gdamore/tcell/v2"

	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/components/cue/helpers"
	"github.com/hofstadter-io/hof/lib/yagu"
)

func (F *Flower) setThinking(thinking bool, which string) {
	c := tcell.ColorWhite
	if thinking {
		c = tcell.ColorViolet
	}

	switch which {
	case "scope":
		F.scope.viewer.SetBorderColor(c)

	case "final":
		F.final.viewer.SetBorderColor(c)

	default:
		F.scope.viewer.SetBorderColor(c)
		F.edit.SetBorderColor(c)
		F.final.viewer.SetBorderColor(c)
	}
	go tui.Draw()
}

func (F *Flower) Rebuild() error {
	tui.Log("info", fmt.Sprintf("Flow.rebuildScope %v", F.scope.config))
	var (
		v cue.Value
		err error
	)

	// get latest user text
	src := F.edit.GetText()

	var sv cue.Value

	{
		F.setThinking(true, "scope")
		defer F.setThinking(false, "scope")

		// get latest scope
		sv, err = F.scope.config.GetValue()
		if err != nil {
			tui.Log("error", err)
		}

		F.scope.viewer.Rebuild()
	}

	{
		F.setThinking(true, "final")
		defer F.setThinking(false, "final")
		// rebuild the value
		ctx := sv.Context()
		v = ctx.CompileString(src, cue.InferBuiltins(true), cue.Scope(sv))
		cfg := &helpers.SourceConfig{Value: v}

		// only update view value, that way, if we erase everything, we still see the value
		F.final.config = cfg
		F.final.viewer.SetSourceConfig(cfg)
		F.final.viewer.Rebuild()
	}

	tui.Draw()
	return nil
}

func (F *Flower) HandleAction(action string, args []string, context map[string]any) (bool, error) {
	tui.Log("warn", fmt.Sprintf("Flower.HandleAction: %v %v %v", action, args, context))
	var err error
	handled := true

	// item actions
	switch action {
	case "export":
		if len(args) != 1 {
			err = fmt.Errorf("export requires a filename")
		} else {
			filename := args[0]
			err := F.ExportFinalToFile(filename)
			// if ok...
			if err == nil {
				msg := fmt.Sprintf("value exported to %s", filename)
				tui.Tell("error", msg)
				tui.Log("trace", msg)
			}
		}

	case "update":
		err = F.updateFromArgsAndContext(args, context)

	case "set.value":
		F.setThinking(true, "final")
		defer F.setThinking(false, "final")
		context["target"] = "value"
		err = F.updateFromArgsAndContext(args, context)
	case "set.scope":
		F.setThinking(true, "scope")
		defer F.setThinking(false, "scope")
		context["target"] = "scope"
		err = F.updateFromArgsAndContext(args, context)

	default:
		handled = false
		// err = fmt.Errorf("unknown command %q", action)
	}

	return handled, err
}

func (F *Flower) updateFromArgsAndContext(args[] string, context map[string]any) error {
	// tui.Log("warn", fmt.Sprintf("Flower.updateHandler.1: %v %v", args, context))
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
	srcCfg := &helpers.SourceConfig{
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

	tui.Log("warn", fmt.Sprintf("Flower.updateHandler.2: %v %v %v", source, target, srcCfg))

	switch target {
	case "value":
		// local source default, assume it was a filename
		if source == "" {
			source = "file"
		} else if source == "<new-play>" {
			source = ""
		}
		srcCfg.Source = helpers.EvalSource(source)

		// tui.Log("warn", fmt.Sprintf("Flower.updateHandler.3.V: %v", srcCfg))
		// tui.Log("extra", fmt.Sprintf("Eval.playItem.value: %v", srcCfg ))
		// C.UseScope(false)
		// need to get the text once at startup
		txt, err := srcCfg.GetText()
		if err != nil {
			tui.Log("error", err)
			return err
		}
		// tui.Log("extra", fmt.Sprintf("Eval.playItem.value.text: %v", txt ))
		F.edit.SetText(txt, false)

	case "scope":
		if source == "" {
			source = "runtime"
		}
		srcCfg.Source = helpers.EvalSource(source)

		// tui.Log("warn", fmt.Sprintf("Flower.updateHandler.3.S: %v", srcCfg))
		F.scope.config = srcCfg
		F.scope.viewer.SetSourceConfig(srcCfg)
	}

	return F.Rebuild()
}

func (F *Flower) ExportFinalToFile(filename string) (error) {
	ext := filepath.Ext(filename)
	ext = strings.TrimPrefix(ext, ".")
	src, err := F.final.viewer.GetValueText(ext)
	if err != nil {
		return err
	}

	dir := filepath.Dir(filename)
	err = yagu.Mkdir(dir)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, []byte(src), 0644)
}
