package eval

import (
	"fmt"

	"cuelang.org/go/cue"
	"github.com/gdamore/tcell/v2"

	"github.com/hofstadter-io/hof/lib/runtime"
	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/tview"
)

type Panel struct {
	*tview.Flex

	_parent *Panel

	creator func (context map[string]any) tview.Primitive

	_Runtime *runtime.Runtime
	_Value   cue.Value
	_content string

	_cnt int
}

func (P *Panel) Id() string {
	return fmt.Sprintf("p:%d", P._cnt)
}

var panel_count = 0
func NewPanel(parent *Panel) *Panel {
	P := &Panel{
		Flex: tview.NewFlex(),
		creator: defaultCreator,
		_cnt: panel_count,
		_parent: parent,
	}
	panel_count++

	P.Flex.SetBorderColor(tcell.Color42)
	P.SetBorder(true).SetTitle(fmt.Sprintf("  %s  ↺  ", P.Id()))

	return P
}


func (P *Panel) Focus(delegate func(p tview.Primitive)) {
	tui.Log("warn", "Panel.Focus " + P.Id())
	if P.GetItemCount() > 0 {
		P.Flex.Focus(delegate)
	}
	tui.SetFocus(P.Flex)
}

func (P *Panel) Mount(context map[string]any) error {
	tui.Log("trace", fmt.Sprintf("Panel.Mount: %v", context))
	// this is where we can do some loading
	P.Flex.Mount(context)

	err := P.Refresh(context)
	if err != nil {
		tui.SendCustomEvent("/console/error", err)
		return err
	}

	// mount any other components
	// maybe we should have [...Children], so two pointers, one for dev, one for sys (Children)
	// then this call to mount can be handled without extra stuff by default?
	//M.View.Mount(context)
	//M.Eval.Mount(context)

	return nil
}

func (P *Panel) Unmount() error {
	// this is where we can do some unloading, depending on the application
	//M.View.Unmount()
	//M.Eval.Unmount()
	P.Flex.Unmount()

	return nil
}

func (P *Panel) Refresh(context map[string]any) error {
	// tui.Log("trace", fmt.Sprintf("Panel.Refresh: %v", context))

	// get and setup args
	args := []string{}
	if _args, ok := context["args"]; ok {
		args = _args.([]string)
	}
	args, context = processArgsAndContext(args, context)

	// extract some info from context
	action := ""
	if _action, ok := context["action"]; ok {
		action = _action.(string)
	}
	_index, _ := context["index"]
	index, _ := _index.(int)
	_where, _ := context["where"]
	where, _ := _where.(string)
	flexDir := P.Flex.GetDirection()

	tui.Log("trace", fmt.Sprintf("Panel.Refresh: %v %v %v %v %v", P.Id(), action, where, index, flexDir))

	// do things based on context info to build up a component
	switch action {
	case "insert":
		P.insertPanelItem(action, where, index, context)	
	case "move":
		P.movePanelItem(action, where, index, context)

	case "nav":

	case "conn":

	case "edit":

	case "mode":

	case "split":
		P.splitPanelItem(action, where, index, context)

	case "delete":
		P.deletePanelItem(action, where, index, context)
		if P.Flex.GetItemCount() == 0 {
			P.AddItem(P.creator(context), 0, 1, true)
		}

	default:
	}

	tui.Draw()
	return nil
}

func (P *Panel) insertPanelItem(action, where string, index int, context map[string]any) {
	tui.Log("trace", fmt.Sprintf("Panel.insertPanelItem %v %v %v %v %v", action, where, index, P.Id(), P.GetItemCount()))
	if i := P.ChildFocus(); i >= 0 {
		itm := P.GetItem(i)
		switch itm := itm.(type) {
		case *Panel:
			itm.insertPanelItem(action, where, index, context)
			return
		default:
			tui.Log("trace", fmt.Sprintf("Panel.childHorz %v %v %v %v %v", action, where, index, P.Id(), P.GetItemCount()))
		}
	}
	
	switch where {

	case "prev":
		t := P.creator(context)
		P.Flex.InsItem(index, t, 0, 1, true)
		tui.SetFocus(t)

	case "next":
		t := P.creator(context)
		P.Flex.InsItem(index+1, t, 0, 1, true)
		tui.SetFocus(t)

	case "index":
		t := P.creator(context)
		P.Flex.InsItem(index, t, 0, 1, true)
		tui.SetFocus(t)

	default:
		return

	} // end: switch where
}

func (P *Panel) movePanelItem(action, where string, index int, context map[string]any) {

	p := P.GetMostFocusedPanel()
	c := p.GetItemCount()
	i := p.ChildFocus()

	tui.Log(
		"trace",
		fmt.Sprintf(
			"Panel.movePanelItem %v %v %v %v>%v %v/%v",
			action, where, index, P.Id(),p.Id(), i,c, 
		),
	)
	if c < 2 {
		return 
	}

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
	tui.Log("trace", fmt.Sprintf("swapping %d & %d in %s", i,j,p.Id()))
	p.SwapIndexes(i,j)
}

func (P *Panel) deletePanelItem(action, where string, index int, context map[string]any) {

	p := P.GetMostFocusedPanel()
	i := p.ChildFocus()

	// do the removal
	if i >= 0 {
		p.RemoveIndex(i)
	} else {
		pp := p._parent
		pp.RemoveItem(p)
		p = pp
	}

	// do some cleanup
	if p.GetItemCount() == 0 {
		if p._parent != nil {
			// remove ourself if parented
			p._parent.RemoveItem(p)
		} else {
			// add default item, we are the root
			p.AddItem(p.creator(context), 0, 1, true)
		}
	}

}

func (P *Panel) splitPanelItem(action, where string, index int, context map[string]any) {

	p := P.GetMostFocusedPanel()
	i := p.ChildFocus()

	tui.Log("error", fmt.Sprintf("Panel.split: %v %v", p.Id(), i))

	// there is a child that we are going to split
	if i >= 0 {
		// shortcut, just add if there aren't enough children
		// they can hit it twice to get the next split
		if p.GetItemCount() < 2 {
			p.AddItem(p.creator(context), 0, 1, true)
		}

		c := p.GetItem(i)
		d := p.GetDirection()
		if d == tview.FlexColumn {
			d = tview.FlexRow
		} else {
			d = tview.FlexColumn
		}

		switch c.(type) {
		case *Item:
			// make a new panel, opposite dir
			n := NewPanel(p)
			n.Flex.SetDirection(d)
			n.AddItem(c, 0, 1, true)
			n.AddItem(n.creator(context), 0, 1, true)
			setupEventHandlers(n, nil, nil)

			p.SetItem(i, n, 0, 1, true)

		}


	} else {
		// pp := p._parent
		// otherwise 0,1 children, so just add
		// not sure we will get here...
		p.AddItem(p.creator(context), 0, 1, true)
	}

}

func (P *Panel) splitPanelItemOld(action, where string, index int, context map[string]any) {
	tui.Log("trace", fmt.Sprintf("Panel.split %v %v %v %v %v", action, where, index, P.Id(), P.GetItemCount()))
	var childItem tview.Primitive
	itm := P.GetChildFocusItem()
	switch itm := itm.(type) {
	case *Panel:
		// there's a child panel with focus, so recurse and delegate
		itm.splitPanelItem(action, where, index, context)
		return
	default:
		// there's a child with focus, it should be a leaf (though could be a complex component itself, there are no panels below)
		tui.Log("trace", fmt.Sprintf("Panel.splitChild %v %v %v %v %v", action, where, index, P.Id(), P.GetItemCount()))
		childItem = itm
	}

	if P.GetItemCount() < 2 {
		P.AddItem(P.creator(context), 0, 1, true)
		return
	}

	// action == "split"
	d := P.Flex.GetDirection()
	if d == tview.FlexColumn {
		d = tview.FlexRow
	} else {
		d = tview.FlexColumn
	}

	// new panel
	p := NewPanel(P)
	p.Flex.SetDirection(d)

	if childItem != nil {

	}

	// opposite direction
	p.AddItem(P.Flex.GetItem(index), 0, 1, true)
	p.AddItem(P.creator(context), 0, 1, true)
	setupEventHandlers(p, nil, nil)

	P.Flex.SetItem(index, p, 0, 1, true)
		// InsItem(index, P.creator(context), 0, 1, true)
}


func (P *Panel) SetShowBordersR(showPanel, showOther bool) {
	P.SetBorder(showPanel)

	for _, i := range P.GetItems() {
		switch t := i.Item.(type) {
		case *Panel:
			t.SetShowBordersR(showPanel, showOther)
		case *Item:
			t.SetBorder(showOther)
		case *tview.Box:
			t.SetBorder(showOther)
		}
	}

}

func (P *Panel) GetMostFocusedPanel() *Panel {

	for _, i := range P.GetItems() {
		switch t := i.Item.(type) {
		case *Panel:
			p := t.GetMostFocusedPanel()
			if p != nil {
				return p
			}
		//case *Item:
		//  // we have a non-panel item that is focused
		//  // so it is us
		//  if t.HasFocus() {
		//    return P
		//  }
		}
	}

	// in theory, we could get here if
	// a panel could have focus without items having focus, not sure this is possible
	if P.HasFocus() {
		return P
	}
	return nil
}

func (P *Panel) GetMostFocusedParent() *Panel {

	for _, i := range P.GetItems() {
		switch t := i.Item.(type) {
		case *Panel:
			p := t.GetMostFocusedParent()
			if p != nil {
				return P
			}

		case tview.Primitive:
			if t.HasFocus() {
				return P
			}
		}
	}

	return nil
}

func (P *Panel) FlipFlexDirection() {
	flexDir := P.Flex.GetDirection()
	if flexDir == tview.FlexRow {
		flexDir = tview.FlexColumn
	} else {
		flexDir = tview.FlexRow
	}
	P.Flex.SetDirection(flexDir)
}
