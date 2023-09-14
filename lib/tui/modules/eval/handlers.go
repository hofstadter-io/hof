package eval

import (
	"fmt"
	"reflect"

	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/components/cue/browser"
	"github.com/hofstadter-io/hof/lib/tui/components/cue/playground"
	"github.com/hofstadter-io/hof/lib/tui/components/panel"
	"github.com/hofstadter-io/hof/lib/tui/components/widget"
	"github.com/hofstadter-io/hof/lib/tui/tview"
)

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

	case "conn":
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

	dstItem := p.GetChildFocusItem().(panel.PanelItem)

	var receiver widget.ActionHandler
	dstPrim := dstItem.Widget()
	switch t := dstPrim.(type) {
	case widget.ActionHandler:
		receiver = t
	default:
		return true, fmt.Errorf("connection destination is not a supported widget, got %v", reflect.TypeOf(dstPrim))
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

	tui.Log("trace", fmt.Sprintf("setting connection for %s -> %s", srcItem.Id(), dstItem.Id()))
	switch t := srcItem.Widget().(type) {
	case widget.ValueProducer:
		fn := t.GetValue
	  if expr != "" {
	    fn = t.GetValueExpr(expr)
	  }
		context["source"] = "conn"
		context["fn"] = fn
	}

	// hmm, add vs set, how do we know here?
	handled, err = receiver.HandleAction("add", args, context)

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

	case "insert":
		// intentionally do nothing?
		return true, nil
	}

	return false, nil
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
	itm, ok := p.GetItem(cfi).(*panel.BaseItem)
	if !ok {
		// not a type we care about (probably another panel)
		return false, nil
	}
	w := itm.Widget()
	switch t := w.(type) {
	case *playground.Playground:
		handled, err = t.HandleAction(action, args, context)
	case *browser.Browser:
		handled, err = t.HandleAction(action, args, context)
	}

	return handled, err
}

