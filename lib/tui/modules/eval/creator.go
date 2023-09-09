package eval

import (
	"fmt"
	"os"
	"strconv"

	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/components/panel"
	"github.com/hofstadter-io/hof/lib/tui/components/widget"
	"github.com/hofstadter-io/hof/lib/tui/components/cue/browser"
	"github.com/hofstadter-io/hof/lib/tui/components/cue/helpers"
	"github.com/hofstadter-io/hof/lib/tui/components/cue/playground"
)

// used for debugging panel CRUD & KEYS
var panel_debug = false

func init() {
	if v := os.Getenv("HOF_TUI_PANEL_DEBUG"); v != "" {
		vb, _ := strconv.ParseBool(v)
		if vb {
			panel_debug = true
		}
	}

	if !panel_debug {
		setupCreator()
	}
}

var itemCreator *panel.Factory

func setupCreator() {
	f := panel.NewFactory()

	f.Register("default", helpItem)
	f.Register("help", helpItem)
	f.Register("play", playItem)
	f.Register("tree", treeItem)

	itemCreator = f
}


// this function is responsable for creating the components that fill slots in the panel
// these are the widgets that make up the application and should have their own operation
func (E *Eval) creator(context panel.ItemContext, parent *panel.Panel) (panel.PanelItem, error) {
	tui.Log("extra", fmt.Sprintf("Eval.creator: %v", context ))

	// short-circuit for developer mode (first, before user custom)
	if panel_debug {
		t := widget.NewTextView()
		i := panel.NewBaseItem(context, parent)
		i.SetWidget(t)
		return i, nil
	}

	// set default item
	if _, ok := context["item"]; !ok {
		context["item"] = "help"
	}

	return itemCreator.Creator(context, parent)
}

func helpItem(context panel.ItemContext, parent *panel.Panel) (panel.PanelItem, error) {
	// tui.Log("extra", fmt.Sprintf("new helpItem %v", context ))
	I := panel.NewBaseItem(context, parent)

	txt := widget.NewTextView()
	fmt.Fprint(txt, EvalHelpText)

	I.SetWidget(txt)

	return I, nil
}

func playItem(context panel.ItemContext, parent *panel.Panel) (panel.PanelItem, error) {
	// tui.Log("extra", fmt.Sprintf("Eval.playItem.context: %v", context ))

	args := []string{}
	if _args, ok := context["args"]; ok {
		args = _args.([]string)
	}

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

	// special case
	if len(args) == 0 || (len(args) == 1 && args[0] == "new") {
		source = "<none>" // very temporary setting
		target = "value"
	}

	// tui.Log("extra", fmt.Sprintf("Eval.playItem.config: %v", srcCfg ))

	play := playground.New("", srcCfg)

	rebuildScope := false
	switch target {
	case "value":
		// local source default, assume it was a filename
		if source == "" {
			source = "file"
		} else if source == "<none>" {
			source = ""
		}
		srcCfg.Source = helpers.EvalSource(source)

		// tui.Log("extra", fmt.Sprintf("Eval.playItem.value: %v", srcCfg ))
		play.UseScope(false)
		// need to get the text once at startup
		txt, err := srcCfg.GetText()
		if err != nil {
			return nil, err
		}
		// tui.Log("extra", fmt.Sprintf("Eval.playItem.value.text: %v", txt ))
		play.SetText(txt)

	case "scope":
		if source == "" {
			source = "runtime"
		}
		srcCfg.Source = helpers.EvalSource(source)

		// tui.Log("extra", fmt.Sprintf("Eval.playItem.scope: %v", srcCfg ))

		play.SetScopeConfig(srcCfg)

		rebuildScope = true
		play.UseScope(true)
	}



	play.Rebuild(rebuildScope)

	I := panel.NewBaseItem(context, parent)
	I.SetWidget(play)

	return I, nil
}

func treeItem(context panel.ItemContext, parent *panel.Panel) (panel.PanelItem, error) {
	// tui.Log("extra", fmt.Sprintf("new treeItem %v", context ))

	args := []string{}
	if _args, ok := context["args"]; ok {
		args = _args.([]string)
	}

	// get source, defaults to runtime
	source := "runtime"
	if _source, ok := context["source"]; ok {
		source = _source.(string)
	}

	srcCfg := helpers.SourceConfig{
		Source: helpers.EvalSource(source),
		Args: args,
	}
	if len(args) == 1 && args[0] == "new" {
		srcCfg.Source = helpers.EvalNone
	}

	b := browser.New(srcCfg, "cue")
	b.SetTitle(fmt.Sprintf("  %v  ", args)).SetBorder(true)
	b.Rebuild()

	I := panel.NewBaseItem(context, parent)
	I.SetWidget(b)

	return I, nil
}
