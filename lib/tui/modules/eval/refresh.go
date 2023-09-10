package eval

import (
	"fmt"

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
		} else {
			return nil
		}	
	}

	// handle dashboard (panel collection), CRUD actions
	handled, err := M.handleDashActions(action, args, context)
	if err != nil {
		tui.Tell("error", err)
		tui.Log("error", err)
		return err
	}
	if handled {
		tui.Log("warn", fmt.Sprintf("Eval.handleDashActions: %v %v", action, args))
		return nil
	}

	// get the current focused panel
	p := M.GetMostFocusedPanel()
	if p == nil {
		p = M.Panel
	}

	handled, err = M.handlePanelActions(p, action, args, context)
	if handled {
		tui.Log("warn", fmt.Sprintf("Eval.handlePanelActions: %v %v", action, args))
		p.Refresh(context)
		return nil
	}

	handled, err = M.handleItemActions(p, action, args, context)
	if handled {
		tui.Log("warn", fmt.Sprintf("Eval.handleItemActions: %v %v", action, args))
		p.Refresh(context)
		return nil
	}

	p.Refresh(context)

	err = fmt.Errorf("unhandled inputs: %v %v %v", action, args, context)
	tui.Log("warn", err)
	return nil

	return err
}

func (M *Eval) handleDashActions(action string, args []string, context map[string]any) (bool, error) {
	var err error
	handled := true

	switch action {
	case "save":
		if len(args) < 1 {
			err = fmt.Errorf("missing filename")
		} else {
			err = M.Save(args[0])
		}

	case "load":
		if len(args) < 1 {
			err = fmt.Errorf("missing filename")
		} else {
			_, err = M.LoadEval(args[0])
		}

	case "show":
		if len(args) < 1 {
			err = fmt.Errorf("missing filename")
		} else {
			_, err = M.ShowEval(args[0])
		}

	case "list":
		err = M.ListEval()
	default:
	  handled = false
	}

	return handled, err
}

func (M *Eval) handlePanelActions(p *panel.Panel, action string, args []string, context map[string]any) (bool, error) {
	// panel naviation
	// this should move down to the panel as well
	switch action {
	case 
		"nav.left",	
		"nav.right",
		"nav.up",
		"nav.down":

		M.doNav(p, action)
		return true, nil
	}

	return false, nil
}

func (M *Eval) handleItemActions(p *panel.Panel, action string, args []string, context map[string]any) (bool, error) {
	var err error
	var handled bool

	// if not panel action, then item?
	// get the item, save if playground
	cfi := p.ChildFocus()
	if cfi < 0 {
		cfi = 0
		// return false, nil
	}
	itm := p.GetItem(cfi).(*panel.BaseItem)
	w := itm.Widget()
	switch t := w.(type) {
	case *playground.Playground:
		handled, err = t.HandleAction(action, args, context)
	}

	return handled, err
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

