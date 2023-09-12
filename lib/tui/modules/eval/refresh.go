package eval

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/components/cue/browser"
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
			err := fmt.Errorf("unknown command: %q {%s}", action, strings.Join(args, ","))
			tui.Tell("error", err)
			tui.Log("error", err)
			return nil
		}	
	}

	// get the current focused panel
	p := M.GetMostFocusedPanel()
	if p == nil {
		p = M.Panel
	}

	// tui.Log("warn", fmt.Sprintf("Eval.handleDashActions.BEFORE: %v %v", action, args))
	// handle dashboard (panel collection), CRUD actions
	handled, err := M.handleDashActions(p, action, args, context)
	if err != nil {
		tui.Tell("error", err)
		tui.Log("error", err)
		return err
	}
	if handled {
		tui.Log("warn", fmt.Sprintf("Eval.handleDashActions: %v %v %v", handled, action, args))
		return nil
	}

	// panel actions (that need to happen at the eval level)
	handled, err = M.handlePanelActions(p, action, args, context)
	if handled {
		tui.Log("warn", fmt.Sprintf("Eval.handlePanelActions: %v %v", action, args))
		p.Refresh(context)
		return nil
	}

	// item actions (passed down to current item)
	handled, err = M.handleItemActions(p, action, args, context)
	if handled {
		tui.Log("warn", fmt.Sprintf("Eval.handleItemActions: %v %v", action, args))
		p.Refresh(context)
		return nil
	}

	if !handled {
		err = fmt.Errorf("unhandled inputs: %v %v", action, args)
		tui.Tell("crit", err)
		tui.Log("error", err)
	}

	p.Refresh(context)

	return nil
}

func (M *Eval) handleDashActions(p *panel.Panel, action string, args []string, context map[string]any) (bool, error) {
	var err error
	handled := true

	switch action {
	case "preview":
		if len(args) < 1 {
			err = fmt.Errorf("missing argument to preview")
		} else {
			err = M.Save(args[0], true)
		}

	case "save":
		if len(args) < 1 {
			err = fmt.Errorf("missing argument to save")
		} else {
			err = M.Save(args[0], false)
		}

	case "load":
		if len(args) < 1 {
			err = fmt.Errorf("missing argument to load")
		} else {
			_, err = M.LoadEval(args[0])
		}

	case "show":
		if len(args) < 1 {
			err = fmt.Errorf("missing argument to show")
		} else {
			_, err = M.ShowEval(args[0])
		}

	case "list":
		err = M.ListEval()

	case "connect":
		return M.handleConnectAction(p, args, context)

	default:
	  handled = false
	}

	return handled, err
}

func (M *Eval) handleConnectAction(p *panel.Panel, args []string, context map[string]any) (bool, error) {
	var (
		err error
		handled bool
	)

	if len(args) < 1 {
		return false, fmt.Errorf("missing connect path")
	}

	var dstPlay *playground.Playground
	dstPrim := p.GetChildFocusItem().(panel.PanelItem).Widget()
	switch t := dstPrim.(type) {
	case *playground.Playground:
		dstPlay = t
	default:
		return true, fmt.Errorf("connection destination must be a playground, got %v", reflect.TypeOf(dstPrim))
	}

	path := args[0]
	srcItem, err := M.getItemByPath(path)
	if err != nil {
		return true, err
	}

	var expr string
	if len(args) == 2 {
		expr = args[1]
	}

	tui.Log("trace", fmt.Sprintf("setting connection for %s -> %s", srcItem.Id(), dstPlay.Id()))
	handled = true
	switch t := srcItem.Widget().(type) {
	case *playground.Playground:
		if expr == "" {
			dstPlay.SetConnection(args, t.GetConnValue)
		} else {
			dstPlay.SetConnection(args, t.GetConnValueExpr(expr))
		}
	case *browser.Browser:
		if expr == "" {
			dstPlay.SetConnection(args, t.GetConnValue)
		} else {
			dstPlay.SetConnection(args, t.GetConnValueExpr(expr))
		}
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

