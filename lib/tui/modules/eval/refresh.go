package eval

import (
	"fmt"
	"reflect"

	"github.com/atotto/clipboard"

	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/components/cue/playground"
	"github.com/hofstadter-io/hof/lib/tui/components/panel"
	"github.com/hofstadter-io/hof/lib/tui/tview"
)

// this is basically the entryp point
func (M *Eval) Refresh(context map[string]any) error {

	// reprocess args, all commands should enter the Eval page first
	// needed for when we come in from the command line first time, or the command box later
	context = enrichContext(context)
	args := []string{}
	if _args, ok := context["args"]; ok {
		args = _args.([]string)
	}
	tui.Log("debug", fmt.Sprintf("Eval.Refresh: %v %# v", args, context))

	// handle any top-leval eval commands
	action := ""
	if _action, ok := context["action"]; ok {
		action = _action.(string)
	}
	// default action to update if item is set
	if action == "" {
		// default action to update if item is set
		if _, ok := context["item"]; ok {
			context["action"] = "update"
		}
	}

	// intercept our top-level commands first
	switch action {
	case "save":
		if len(args) < 1 {
			err := fmt.Errorf("missing filename")
			tui.Tell("error", err)
			tui.Log("error", err)
			return nil
		}
		return M.Save(args[0])

	case "load":
		if len(args) < 1 {
			err := fmt.Errorf("missing filename")
			tui.Tell("error", err)
			tui.Log("error", err)
			return err
		}
		_, err := M.LoadEval(args[0])
		if err != nil {
			tui.Tell("error", err)
			tui.Log("error", err)
			return err
		}
		return nil

	case "show":
		if len(args) < 1 {
			err := fmt.Errorf("missing filename")
			tui.Tell("error", err)
			tui.Log("error", err)
			return err
		}
		_, err := M.ShowEval(args[0])
		if err != nil {
			tui.Tell("error", err)
			tui.Log("error", err)
			return err
		}
		return nil

	case "list":
		err := M.ListEval()
		if err != nil {
			tui.Tell("error", err)
			tui.Log("error", err)
			return err
		}
		return nil
	}

	p := M.GetMostFocusedPanel()
	if p == nil {
		p = M.Panel
	}

	// item actions
	switch action {
	case "push":

		tui.Log("debug", fmt.Sprintf("push cmd: %# v", context))
		cfi := p.ChildFocus()

		itm := p.GetItem(cfi).(*panel.BaseItem)
		w := itm.Widget()
		switch play := w.(type) {
		case *playground.Playground:
			id, err := play.PushToPlayground()
			if err != nil {
				tui.Tell("error", err)
				tui.Log("error", err)
				return err
			}

			msg := fmt.Sprintf("snippet id: %s  (link copied!)", id)

			url := fmt.Sprintf("https://cuelang.org/play?id=%s", id)
			clipboard.WriteAll(url)

			tui.Tell("error", msg)
			tui.Log("trace", msg)
			return nil

		default:
			err := fmt.Errorf("unable to push this item %v", reflect.TypeOf(w))
			tui.Tell("error", err)
			tui.Log("error", err)
			return err
		}

	case "write":
		tui.Log("debug", fmt.Sprintf("write cmd: %# v", context))
		cfi := p.ChildFocus()

		if len(args) != 1 {
			err := fmt.Errorf("write requires a filename")
			tui.Tell("error", err)
			tui.Log("error", err)
			return err
		}

		filename := args[0]

		itm := p.GetItem(cfi).(*panel.BaseItem)
		w := itm.Widget()
		switch play := w.(type) {
		case *playground.Playground:
			err := play.WriteEditToFile(filename)
			if err != nil {
				tui.Tell("error", err)
				tui.Log("error", err)
				return err
			}

			msg := fmt.Sprintf("editor text saved to %s", filename)

			tui.Tell("error", msg)
			tui.Log("trace", msg)
			return nil

		default:
			err := fmt.Errorf("unable to write this item %v", reflect.TypeOf(w))
			tui.Tell("error", err)
			tui.Log("error", err)
			return err
		}

	case "export":
		tui.Log("debug", fmt.Sprintf("export cmd: %# v", context))
		cfi := p.ChildFocus()

		if len(args) != 1 {
			err := fmt.Errorf("export requires a filename")
			tui.Tell("error", err)
			tui.Log("error", err)
			return err
		}

		filename := args[0]

		itm := p.GetItem(cfi).(*panel.BaseItem)
		w := itm.Widget()
		switch play := w.(type) {
		case *playground.Playground:
			err := play.ExportFinalToFile(filename)
			if err != nil {
				tui.Tell("error", err)
				tui.Log("error", err)
				return err
			}

			msg := fmt.Sprintf("value exported to %s", filename)

			tui.Tell("error", msg)
			tui.Log("trace", msg)
			return nil

		default:
			err := fmt.Errorf("unable to export this item %v", reflect.TypeOf(w))
			tui.Tell("error", err)
			tui.Log("error", err)
			return err
		}

	case 
		"nav.left",	
		"nav.right",
		"nav.up",
		"nav.down":

		M.doNav(p, action)

	case "":
		// the empty string should only happen on startup

	default:
		err := fmt.Errorf("unknown command %q", action)
		tui.Tell("error", err)
		tui.Log("error", err)
		return err
	}

	err := p.Refresh(context)
	if err != nil {
		return M.showError(err)	
	}

	return nil
}

func (M *Eval) doNav(panel *panel.Panel, action string) {

	dir := panel.GetDirection()

	// calculated movement values
	local := false // local to panel, or between panels
	inc := false   // increase index, or decrease index
	switch action {
	case "nav.left", "nav.right":
		local = dir == tview.FlexColumn
		inc = action == "nav.right"
	case "nav.up", "nav.down":
		local = dir == tview.FlexRow
		inc = action == "nav.down"
	}

	// let's move to the parent
	if !local {
		p := panel.GetParent()
		// nothing to do, we should be at the root
		if p == nil {
			return
		}
		// update panel, so now we work on parent
		panel = p
	}

	cfi := panel.ChildFocus()
	cnt := panel.GetItemCount()
	// nothing to do, shortcut return

	if cnt < 2 {
		return
	}

	// new index to focus
	if inc {
		cfi += 1
	} else {
		cfi -= 1
	}
	if cfi < 0 {
		cfi = 0
	}
	if cfi >= cnt {
		cfi = cnt - 1
	}

	tui.SetFocus(panel.GetItem(cfi))
}
