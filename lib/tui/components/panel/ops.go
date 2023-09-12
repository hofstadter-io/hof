package panel

import (
	"fmt"

	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/tview"
)

func (P *Panel) insertPanelItem(context map[string]any) {
	where := "tail"
	if _where, ok := context["where"]; ok {
		if w, sok := _where.(string); sok {
			where = w
		} else {
			tui.Log("error", fmt.Sprintf("unknown where in Panel.insertPanelItem: %v %#v", P.Id(), context))
		}	
	}

	panel := P
	if _panel, ok := context["panel"]; ok {
		panel = _panel.(*Panel)
	}
	cfi := -1
	if _cfi, ok := context["child-focus-index"]; ok {
		cfi = _cfi.(int)
		// tui.Log("trace", fmt.Sprintf("setting cfi.1 %d\n", cfi))
	}

	if cfi == -1 {
		tui.Log("error", fmt.Sprintf("nil child in Panel.insertPanelItem: %v %#v", panel.Id(), context))
		where = "tail"
	}

	t, _ := panel._creator(context, panel)
	
	switch where {

	case "head":
		panel.Flex.InsItem(0, t, 0, 1, true)

	case "prev":
		panel.Flex.InsItem(cfi, t, 0, 1, true)

	case "next":
		panel.Flex.InsItem(cfi+1, t, 0, 1, true)

	case "tail":
		panel.Flex.AddItem(t, 0, 1, true)

	case "index":
		panel.Flex.InsItem(cfi, t, 0, 1, true)

	default:
		return

	} // end: switch where

	tui.SetFocus(t)
}

func (P *Panel) updatePanelItem(context map[string]any) {
	panel := P
	if _panel, ok := context["panel"]; ok {
		panel = _panel.(*Panel)
	}
	cfi := -1
	if _cfi, ok := context["child-focus-index"]; ok {
		cfi = _cfi.(int)
		// tui.Log("trace", fmt.Sprintf("setting cfi.1 %d\n", cfi))
	}

	i := panel.ChildFocus()
	if i == -1 {
		// tui.Log("warn", fmt.Sprintf("using 0 for nil child in Panel.updatePanelItem: %v %#v", P.Id(), context))
		i = 0
	} else {
		cfi = i
		// tui.Log("trace", fmt.Sprintf("setting cfi.2 %d\n", cfi))
	}
	
	t, _ := panel._creator(context, panel)

	// just insert, this happens on first load and such
	if P.GetItemCount() == 0 {
		panel.Flex.AddItem(t, 0, 1, true)
	}

	if cfi < 0 {
		// tui.Log("error", fmt.Sprintf("negative cfi %# v\n", context))

		// a bit of hackery, seems this happens on startup, because there is no focus yet
		// we should probably solve this by setting a focus / initial component correctly
		cnt := panel.GetItemCount()
		if cnt == 0 {
			panel.Flex.AddItem(t, 0, 1, true)
		} else if cnt == 1 {
			panel.Flex.SetItem(0, t, 0, 1, true)
		}
		return
	}

	// update a position
	panel.Flex.SetItem(cfi, t, 0, 1, true)

	tui.SetFocus(t)
}

func (P *Panel) movePanelItem(context map[string]any) {

	panel := P
	if _panel, ok := context["panel"]; ok {
		panel = _panel.(*Panel)
	}
	cfi := -1
	if _cfi, ok := context["child-focus-index"]; ok {
		cfi = _cfi.(int)
		// tui.Log("trace", fmt.Sprintf("setting cfi.1 %d\n", cfi))
	}

	c := panel.GetItemCount()
	i := cfi

	if c < 2 {
		return 
	}

	_where, _ := context["where"]
	where, _ := _where.(string)

	j := i
	switch where {
	case "prev":
		j--
	case "next":
		j++	
	default:
		tui.Log("error", "unknown movePanel where: " + where)
		return
	}

	// j is out of bounds, do nothing
	if j < 0 || j >= c {
		return
	}

	// otherwise, we should be good to swap
	// tui.Log("trace", fmt.Sprintf("swapping %d & %d in %s", i,j,p.Id()))
	panel.SwapIndexes(i,j)
}

func (P *Panel) deletePanelItem(context map[string]any) {

	panel := P
	if _panel, ok := context["panel"]; ok {
		panel = _panel.(*Panel)
	}
	cfi := -1
	if _cfi, ok := context["child-focus-index"]; ok {
		cfi = _cfi.(int)
		// tui.Log("trace", fmt.Sprintf("setting cfi.1 %d\n", cfi))
	}

	// do the removal
	if cfi >= 0 {
		panel.RemoveIndex(cfi)
	} else {
		pp := panel._parent
		pp.RemoveItem(panel)
		panel = pp
	}

	// do some cleanup
	if panel.GetItemCount() == 0 {

		// unwind towards the root, deleting nested panels with only a single child panel
		// this works by first removing ourself, since we have no children, and then
		// checking after to see if the panel we removed ourself from has no children afterwards
		// we also need to stop when we reach the root
		for panel.GetItemCount() == 0 && panel._parent != nil {
			panel._parent.RemoveItem(panel)
			panel = panel._parent
		}

		// add default item, if we are in an empty panel
		// (which should only be the root at this point)
		if panel.GetItemCount() == 0 {
			// if panel._parent == nil { // old check, new one probably equivalent
			// we don't want to be here if the deletion process landed us in a panel with other elements
			// this code should only add back the default help text when there are no other widgets left
			context["item"] = "default"
			t, _ := panel._creator(context, panel)
			panel.AddItem(t, 0, 1, true)	
		}
	}

	tui.SetFocus(panel)
}

func (P *Panel) splitPanelItem(context map[string]any) {

	panel := P
	if _panel, ok := context["panel"]; ok {
		panel = _panel.(*Panel)
	}
	cfi := -1
	if _cfi, ok := context["child-focus-index"]; ok {
		cfi = _cfi.(int)
		// tui.Log("trace", fmt.Sprintf("setting cfi.1 %d\n", cfi))
	}

	// tui.Log("error", fmt.Sprintf("Panel.split: %v %v", p.Id(), i))

	// there is a child that we are going to split
	if cfi >= 0 {
		// shortcut, just add if there aren't enough children
		// they can hit it twice to get the next split
		if panel.GetItemCount() < 2 {
			t, _ := panel._creator(context, panel)
			panel.AddItem(t, 0, 1, true)
			tui.SetFocus(t)
			return
		}

		c := panel.GetItem(cfi)
		d := panel.GetDirection()
		if d == tview.FlexColumn {
			d = tview.FlexRow
		} else {
			d = tview.FlexColumn
		}

		switch c.(type) {
		case PanelItem:
			// make a new panel, opposite dir
			n := New(panel, nil)
			n.Flex.SetDirection(d)
			n.SetBorder(panel.GetBorder())
			n.AddItem(c, 0, 1, true)
			context["item"] = "default"
			t, _ := n._creator(context, panel)
			n.AddItem(t, 0, 1, true)
			// setupEventHandlers(n, nil, nil)

			panel.SetItem(cfi, n, 0, 1, true)
			tui.SetFocus(n)
		}

	} else {
		// otherwise 0,1 children, so just add
		// not sure we will get here...
		context["item"] = "default"
		t, _ := panel._creator(context, panel)
		panel.AddItem(t, 0, 1, true)
		tui.SetFocus(t)
	}

}
