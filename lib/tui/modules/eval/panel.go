package eval

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"

	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/components"
	"github.com/hofstadter-io/hof/lib/tui/tview"
)

type Panel struct {
	*tview.Flex

	_parent *Panel

	_creator ItemCreator

	// this needs to go away?
	// or we want scope for a whole panel?
	// probably take away for now, so as to not confuse people
	// and hopefully item level is enough with connections, and send/recv value update handlers
	// we should probably send these updates around via the message bus
	//_Runtime *runtime.Runtime
	//_Value   cue.Value
	//_content string

	// todo, make some sort of embedded clase for meta + save/reload?
	_cnt  int
	_name string
}

func (P *Panel) EncodeMap() (map[string]any, error) {
	m := make(map[string]any)

	m["id"] = P._cnt
	m["name"] = P._name
	m["type"] = "panel"
	m["direction"] = P.Flex.GetDirection()

	items := []map[string]any{}

	for _, item := range P.GetItems() {
		var (
			i map[string]any
			err error
		)

		switch item := item.Item.(type) {
		case *Panel:
			i, err = item.EncodeMap()
			if err != nil {
				return m, err
			}
		case *Item:
			i, err = item.EncodeMap()
			if err != nil {
				return m, err
			}

		default:
			panic("unhandled item type in panel")	
		}

		items = append(items, i)
	}

	m["items"] = items

	return m, nil
}

func PanelDecodeMap(data map[string]any, parent *Panel, creator ItemCreator) (*Panel, error) {
	P := &Panel{
		Flex: tview.NewFlex(),
		_creator: creator,
		_parent: parent,
		_cnt: data["id"].(int),
		_name: data["name"].(string),
	}

	if items, ok := data["items"]; ok {
		for _, idata := range items.([]any) {
			imap := idata.(map[string]any)
			I, err := ItemDecodeMap(imap, P)
			if err != nil {
				return P, err
			}
			P.AddItem(I, 0, 1, true)
		}
	} else {
		txt := NewTextView()
		fmt.Fprint(txt, fmt.Sprintf("unhandled panel decode: \n%# v\n\n", data))
		fmt.Fprint(txt, EvalHelpText)
		I := NewItem(nil, parent)
		I.SetWidget(txt)
		P.AddItem(I, 0, 1, true)

	}

	// do layout setup here, once some children have been instantiated
	P.Flex.SetDirection(data["direction"].(int))
	P.Flex.SetBorder(true).SetTitle(P.Name())

	return P, nil
}

func (P *Panel) Id() string {
	return fmt.Sprintf("p:%d", P._cnt)
}

func (P *Panel) Name() string {
	return P._name
}

func (P *Panel) SetName(name string) {
	P._name = name
}

var panel_count = 0
func NewPanel(parent *Panel, creator ItemCreator) *Panel {
	P := &Panel{
		Flex: tview.NewFlex(),
		_creator: creator,
		_cnt: panel_count,
		_parent: parent,
	}
	panel_count++

	P.Flex.SetBorderColor(tcell.Color42)
	P.SetBorder(true).SetTitle(fmt.Sprintf("  %s  ↺  ", P.Id()))

	return P
}


func (P *Panel) Focus(delegate func(p tview.Primitive)) {
	// tui.Log("warn", "Panel.Focus " + P.Id())
	if P.GetItemCount() > 0 {
		P.Flex.Focus(delegate)
	}
	tui.SetFocus(P.Flex)
}

func (P *Panel) Mount(context map[string]any) error {
	// tui.Log("trace", fmt.Sprintf("Panel.Mount: %v", context))
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
	tui.Log("extra", fmt.Sprintf("Panel.Refresh: %v %v", P.Id(), context))

	// extract some info from context
	action := ""
	if _action, ok := context["action"]; ok {
		action = _action.(string)
	}
	args := []string{}
	if _args, ok := context["args"]; ok {
		args = _args.([]string)
	}

	cid := P.ChildFocus()

	// do things based on context info to build up a component
	switch action {
	// panel / item layout related
	case "insert":
		P.insertPanelItem(context)	
	case "move":
		P.movePanelItem(context)
	case "split":
		P.splitPanelItem(context)
	case "delete":
		P.deletePanelItem(context)

	// item state related
	case "update":
		P.updatePanelItem(context)	

	case "set.scope", "set.scope.runtime", "set.scope.value":
		P.setItemScope(cid, action, args, context)	

	case "set.value", "set.value.runtime", "set.value.value":
		P.setItemValue(cid, action, args, context)	

	case "set.text":
		P.setItemText(cid, action, args, context)	

	case "set.panel.name":
		if len(args) < 1 {
			return fmt.Errorf("%s requires an argument", action)
		}
		P.SetName(args[0])
		P.SetTitle("  "+P.Name()+"  ")

	case "set.name", "set.item.name":
		if len(args) < 1 {
			return fmt.Errorf("%s requires an argument", action)
		}
		if cid >= 0 {
			i := P.GetItem(cid)
			switch i := i.(type) {
			case *Panel:
				i.SetName(args[0])
				i.SetTitle("  "+i.Name()+"  ")
			case *Item:
				i.SetName(args[0])
				i.SetTitle("  "+i.Name()+"  ")
			}
		} else {
			P.SetName(args[0])
			P.SetTitle("  "+P.Name()+"  ")
		}

	case "conn":
		P.makeItemConnection(args, context)		


	default:
	}

	return nil
}

func (P *Panel) insertPanelItem(context map[string]any) {
	where := "tail"
	if _where, ok := context["where"]; ok {
		if w, sok := _where.(string); sok {
			where = w
		} else {
			tui.Log("error", fmt.Sprintf("unknown where in Panel.insertPanelItem: %v %#v", P.Id(), context))
		}	
	}

	i := P.ChildFocus()
	if i == -1 {
		tui.Log("error", fmt.Sprintf("nil child in Panel.insertPanelItem: %v %#v", P.Id(), context))
		where = "tail"
	}
	
	switch where {

	case "head":
		t, _ := P.creator(context, P)
		P.Flex.InsItem(0, t, 0, 1, true)
		tui.SetFocus(t)

	case "prev":
		t, _ := P.creator(context, P)
		P.Flex.InsItem(i, t, 0, 1, true)
		tui.SetFocus(t)

	case "next":
		t, _ := P.creator(context, P)
		P.Flex.InsItem(i+1, t, 0, 1, true)
		tui.SetFocus(t)

	case "tail":
		t, _ := P.creator(context, P)
		P.Flex.AddItem(t, 0, 1, true)
		tui.SetFocus(t)

	case "index":
		t, _ := P.creator(context, P)
		P.Flex.InsItem(i, t, 0, 1, true)
		tui.SetFocus(t)

	default:
		return

	} // end: switch where
}

func (P *Panel) setItemScope(cid int, action string, args []string, context map[string]any) error {
	i := P.GetItem(cid)
	I := i.(*Item)

	if action == "set.scope.value" {
		if len(args) != 1 {
			return fmt.Errorf("set.scope.value expects a single filepath argument")
		}
		return I.loadValue(args, true)
	}

	return I.loadRuntime(args, true)
}

func (P *Panel) setItemValue(cid int, action string, args []string, context map[string]any) error {
	i := P.GetItem(cid)
	I := i.(*Item)

	switch action {
	case "set.value", "set.value.value":
		if len(args) != 1 {
			return fmt.Errorf("%s expects a single filepath argument", action)
		}
		return I.loadValue(args, true)
	}

	return I.loadRuntime(args, true)
}

func (P *Panel) setItemText(cid int, action string, args []string, context map[string]any) error {
	if len(args) < 1 {
		return fmt.Errorf("set.text requires a filepath argument")
	}
	fp := args[0]

	b, err := os.ReadFile(fp)
	if err != nil {
		return err
	}

	i := P.GetItem(cid)
	I := i.(*Item)

	I._text = string(b)
	switch w := I.Widget().(type) {
	case *components.ValueEvaluator:
		w.SetText(I._text)
	case *tview.TextArea:
		w.SetText(I._text, false)
	case *tview.TextView:
		w.Clear()
		fmt.Fprintln(w, I._text)
	default:
		return fmt.Errorf("setting text is not supported for current item")
	}

	return nil
}

func (P *Panel) makeItemConnection(args []string, context map[string]any) error {


	return nil
}

func (P *Panel) updatePanelItem(context map[string]any) {
	panel := P
	if _panel, ok := context["panel"]; ok {
		panel = _panel.(*Panel)
	}
	cfi := -1
	if _cfi, ok := context["child-focus-index"]; ok {
		cfi = _cfi.(int)
		tui.Log("trace", fmt.Sprintf("setting cfi.1 %d\n", cfi))
	}

	i := panel.ChildFocus()
	if i == -1 {
		tui.Log("warn", fmt.Sprintf("using 0 for nil child in Panel.updatePanelItem: %v %#v", P.Id(), context))
		i = 0
	} else {
		cfi = i
		tui.Log("trace", fmt.Sprintf("setting cfi.2 %d\n", cfi))
	}
	
	t, _ := panel.creator(context, panel)
	tui.SetFocus(t)

	// just insert, this happens on first load and such
	if P.GetItemCount() == 0 {
		panel.Flex.AddItem(t, 0, 1, true)
	}

	if cfi < 0 {
		tui.Log("error", fmt.Sprintf("negative cfi %# v\n", context))
		return
	}

	// update a position
	panel.Flex.SetItem(cfi, t, 0, 1, true)
}

func (P *Panel) movePanelItem(context map[string]any) {

	p := P.GetMostFocusedPanel()
	c := p.GetItemCount()
	i := p.ChildFocus()

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
	p.SwapIndexes(i,j)
}

func (P *Panel) deletePanelItem(context map[string]any) {

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
			tui.SetFocus(p._parent)
		} else {
			// add default item, we are the root
			context["action"] = "insert"
			context["item"] = "help"
			t, _ := p.creator(context, p)
			p.AddItem(t, 0, 1, true)	
			tui.SetFocus(t)
		}
	}

}

func (P *Panel) splitPanelItem(context map[string]any) {

	p := P.GetMostFocusedPanel()
	i := p.ChildFocus()

	// tui.Log("error", fmt.Sprintf("Panel.split: %v %v", p.Id(), i))

	// there is a child that we are going to split
	if i >= 0 {
		// shortcut, just add if there aren't enough children
		// they can hit it twice to get the next split
		if p.GetItemCount() < 2 {
			t, _ := p.creator(context, p)
			p.AddItem(t, 0, 1, true)
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
			n := NewPanel(p, nil)
			n.Flex.SetDirection(d)
			n.AddItem(c, 0, 1, true)
			context["action"] = "insert"
			context["item"] = "help"
			t, _ := n.creator(context, p)
			n.AddItem(t, 0, 1, true)
			// setupEventHandlers(n, nil, nil)

			p.SetItem(i, n, 0, 1, true)
			tui.SetFocus(n)
		}

	} else {
		// otherwise 0,1 children, so just add
		// not sure we will get here...
		context["action"] = "insert"
		context["item"] = "help"
		t, _ := p.creator(context, p)
		p.AddItem(t, 0, 1, true)
		tui.SetFocus(t)
	}

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
