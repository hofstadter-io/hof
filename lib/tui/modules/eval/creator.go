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
	tui.Log("extra", fmt.Sprintf("new helpItem %v", context ))
	I := panel.NewBaseItem(context, parent)

	txt := widget.NewTextView()
	fmt.Fprint(txt, EvalHelpText)

	I.SetWidget(txt)

	return I, nil
}

func playItem(context panel.ItemContext, parent *panel.Panel) (panel.PanelItem, error) {
	tui.Log("extra", fmt.Sprintf("new playItem %v", context ))

	args := []string{}
	if _args, ok := context["args"]; ok {
		args = _args.([]string)
	}

	// get source, defaults to runtime
	source := "runtime"
	if _source, ok := context["source"]; ok {
		source = _source.(string)
	}

	// setup our source config
	srcCfg := helpers.SourceConfig{
		Source: helpers.EvalSource(source),
		Args: args,
	}

	if len(args) == 1 && args[0] == "new" {
		srcCfg.Source = helpers.EvalNone
	}

	e := playground.New("", srcCfg)
	e.Rebuild(true)

	I := panel.NewBaseItem(context, parent)
	I.SetWidget(e)

	return I, nil
}

func treeItem(context panel.ItemContext, parent *panel.Panel) (panel.PanelItem, error) {
	tui.Log("extra", fmt.Sprintf("new treeItem %v", context ))

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
