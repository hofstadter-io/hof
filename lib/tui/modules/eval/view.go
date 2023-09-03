package eval

import (
	"fmt"

	"github.com/gdamore/tcell/v2"

	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/tview"
)

type Eval struct {
	*Panel
}

func NewEval() *Eval {
	M := &Eval{
		Panel: NewPanel(),
	}

	// do layout setup here
	M.Flex.SetDirection(tview.FlexColumn)
	M.Flex.SetBorder(true).SetTitle("  Eval (flex)  ")

	return M
}

var count = 0
var debug = true

func (M *Eval) Mount(context map[string]any) error {
	// this is where we can do some loading
	M.Flex.Mount(context)
	setupInputHandler(M.Panel)

	err := M.Refresh(context)
	if err != nil {
		tui.SendCustomEvent("/console/error", err)
		return err
	}

	return nil
}

func (M *Eval) Unmount() error {
	// this is where we can do some unloading, depending on the application
	M.Flex.Unmount()

	// TODO, unmount all items, or will the above handle it for us?

	// remove keybinds
	M.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey { return event })

	return nil
}

func (M *Eval) Refresh(context map[string]any) error {
	return refresh(M, context)
}

func refresh(M *Eval, context map[string]any) error {
	flexDir := M.Flex.GetDirection()

	// strip off the command name
	_args, _ := context["args"]
	args, _ := _args.([]string)
	if len(args) > 0 && args[0] == "eval" {
		args = args[1:]
		context["args"] = args
	}

	_action, _ := context["action"]
	action, _ := _action.(string)
	_index, _ := context["index"]
	index, _ := _index.(int)
	_where, _ := context["where"]
	where, _ := _where.(string)

	switch action {
	case "delete":
	  M.Flex.RemoveIndex(index)
	case "split":
		M.splitPanelItem(action, where, index, context)		

	case "insert", "scoped":
		if flexDir == tview.FlexColumn {
			M.insertHorz(action, where, index, context)		
		} else {
			M.insertVert(action, where, index, context)		
		}

	default:
		// no action, probably coming to eval for the first time
		M.Flex.AddItem(M.creator(context), 0, 1, true)
	}

	// only set border when no elements
	// M.Flex.SetBorder(M.Flex.GetItemCount() == 0)
	dir := "row"
	if flexDir == tview.FlexColumn {
		dir = "col"
	}

	M.Flex.SetBorder(true).SetTitle(fmt.Sprintf("  Eval (flex-%s)  ", dir))
	tui.Draw()

	return nil
}

func (P *Panel) insertVert(action, where string, index int, context map[string]any) {
	switch where {
	case "top":
		P.Flex.InsItem(0, P.creator(context), 0, 1, true)

	case "above":
		P.Flex.InsItem(index, P.creator(context), 0, 1, true)

	case "below":
		P.Flex.InsItem(index+1, P.creator(context), 0, 1, true)

	case "bottom":
		P.Flex.AddItem(P.creator(context), 0, 1, true)

	case "index":
		P.Flex.InsItem(index, P.creator(context), 0, 1, true)
	}
}

func (P *Panel) insertHorz(action, where string, index int, context map[string]any) {
	switch where {
	case "head":
		P.Flex.InsItem(0, P.creator(context), 0, 1, true)

	case "before":
		P.Flex.InsItem(index, P.creator(context), 0, 1, true)

	case "next":
		P.Flex.InsItem(index+1, P.creator(context), 0, 1, true)

	case "last":
		P.Flex.AddItem(P.creator(context), 0, 1, true)

	case "index":
		P.Flex.InsItem(index, P.creator(context), 0, 1, true)
	}
}

func (P *Panel) splitPanelItem(action, where string, index int, context map[string]any) {
	tui.Log("trace", fmt.Sprint("Panel.slipt", action, where, index))
	// action == "split"
	d := P.Flex.GetDirection()
	if d == tview.FlexColumn {
		d = tview.FlexRow
	} else {
		d = tview.FlexColumn
	}
	// new panel
	p := NewPanel()
	// opposite direction
	p.Flex.SetDirection(d)
	p.Flex.SetBorder(true)
	p.AddItem(P.Flex.GetItem(index), 0, 1, false)
	P.Flex.SetItem(index, p, 0, 1, true)
		// InsItem(index, P.creator(context), 0, 1, true)
}
