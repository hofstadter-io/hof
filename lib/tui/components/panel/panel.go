package panel

import (
	"fmt"
	"strconv"
	"sync/atomic"

	"github.com/gdamore/tcell/v2"

	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/tview"
)


type Panel struct {
	// TODO, this is going to be the best practice, we also hope it will allow use to focus a panel without the item in the flex
	// *tview.Frame

	*tview.Flex

	_parent *Panel

	_creator ItemCreator

	// todo, make some sort of embedded clase for meta + save/reload?
	_cnt  int
	_name string
}

func (P *Panel) Id() string {
	return fmt.Sprintf("p%d", P._cnt)
}

func (P *Panel) Name() string {
	return P._name
}

func (P *Panel) TypeName() string {
	return "panel"
}

func (P *Panel) SetName(name string) {
	P._name = name
}

func (P *Panel) GetParent() *Panel {
	return P._parent
}

func (P *Panel) SetParent(parent *Panel) {
	P._parent = parent
}

func (P *Panel) GetItemById(id string) tview.Primitive {
	items := P.GetItems()
	for _, itm := range items {
		switch i := itm.Item.(type) {
		case *Panel:
			if i.Id() == id {
				return i
			}
		case PanelItem:
			if i.Id() == id {
				return i
			}
		}
	}

	return nil
}

func (P *Panel) GetItemByName(name string) tview.Primitive {
	items := P.GetItems()
	for _, itm := range items {
		switch i := itm.Item.(type) {
		case *Panel:
			if i.Name() == name {
				return i
			}
		case PanelItem:
			if i.Name() == name {
				return i
			}
		}
	}

	return nil
}

var panel_count *atomic.Int64
func init() {
	panel_count = new(atomic.Int64)
}

func New(parent *Panel, creator ItemCreator) *Panel {
	P := &Panel{
		Flex: tview.NewFlex(),
		_creator: creator,
		_cnt: int(panel_count.Add(1)),
		_parent: parent,
	}

	// fallback if needed & possible
	if P._creator == nil && P._parent != nil {
		P._creator = P._parent._creator
	}

	P.Flex.SetBorderColor(tcell.Color42).SetBorder(true)
	P.SetDirection(P.Flex.GetDirection())

	return P
}

func (P *Panel) SetDirection(d int) {
	P.Flex.SetDirection(d)
	P.SetTitle(P.TitleString())
}

func (P *Panel) TitleString() string {
	// glyphs to show flex direction
	dir := ""
	if P.GetDirection() == tview.FlexRow {
		dir = "="
	}	else {
		dir = "|"
	}

	// name (with space if set)
	n := ""
	if _n := P.Name(); _n != "" {
		n = _n + " "
	}

	return fmt.Sprintf("  %s(%s) %s ", n, P.Id(), dir)
}


func (P *Panel) Focus(delegate func(p tview.Primitive)) {
	// tui.Log("warn", "Panel.Focus " + P.Id())
	if P.GetItemCount() > 0 {
		P.Flex.Focus(delegate)
	}
	// tui.SetFocus(P)
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
	case "create":
		P.createPanelItem(context)	
	case "insert":
		P.insertPanelItem(context)	
	case "move":
		P.movePanelItem(context)
	case "split":
		P.splitPanelItem(context)
	case "delete":
		P.deletePanelItem(context)
	case "set.panel.name":
		if len(args) < 1 {
			return fmt.Errorf("%s requires an argument", action)
		}
		P.SetName(args[0])
		P.SetTitle(P.TitleString())

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
			case PanelItem:
				i.SetName(args[0])
				i.SetTitle("  "+i.Name()+"  ")
			}
		} else {
			P.SetName(args[0])
			P.SetTitle("  "+P.Name()+"  ")
		}
	case "set.size":
		tui.Log("alert", fmt.Sprintf("setting fixed size %v %v", cid, args))
		if len(args) < 1 {
			return fmt.Errorf("%s requires an int argument", action)
		}
		if cid >= 0 {
			i, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			P.Flex.SetItemFixedSize(cid, i)
		}

	case "set.ratio":
		if len(args) < 1 {
			return fmt.Errorf("%s requires an int argument", action)
		}
		if cid >= 0 {
			i, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			P.Flex.SetItemProportion(cid, i)
		}


	default:
	}

	return nil
}

func (P *Panel) SetShowBordersR(showPanel, showOther bool) {
	P.SetBorder(showPanel)

	for _, i := range P.GetItems() {
		switch t := i.Item.(type) {
		case *Panel:
			t.SetShowBordersR(showPanel, showOther)
		case PanelItem:
			t.SetBorder(showOther)
		case *tview.Box:
			t.SetBorder(showOther)
		}
	}

}

func (P *Panel) GetMostFocusedPanel() *Panel {
	// look for item with focus
	for _, i := range P.GetItems() {
		switch t := i.Item.(type) {
		case *Panel:
			p := t.GetMostFocusedPanel()
			if p != nil {
				return p
			}
		}
	}

	// otherwise, if we have focus, it is us
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
	P.SetDirection(flexDir)

	for _, i := range P.GetItems() {
		switch t := i.Item.(type) {
		case *Panel:
			t.FlipFlexDirection()
		}
	}

}

