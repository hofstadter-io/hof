package eval

import (
	"fmt"

	"github.com/gdamore/tcell/v2"

	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/events"
	"github.com/hofstadter-io/hof/lib/tui/tview"
)

type Eval struct {
	*Panel

	// border display
	showPanel, showOther bool
}

func NewEval() *Eval {
	M := &Eval{
		Panel: NewPanel(nil),
		showPanel: true,
		showOther: true,
	}

	// do layout setup here
	M.Flex.SetDirection(tview.FlexColumn)
	M.Flex.SetBorder(true).SetTitle("  Eval (flex)  ")

	return M
}

func (M *Eval) Mount(context map[string]any) error {

	// this will mount the core element and all children
	M.Flex.Mount(context)
	tui.Log("trace", "Eval.Mount")

	// handle border display
	tui.AddWidgetHandler(M.Panel, "/sys/key/A-P", func(e events.Event) {
		M.showPanel = !M.showPanel
		M.SetShowBordersR(M.showPanel, M.showOther)
	})

	tui.AddWidgetHandler(M.Panel, "/sys/key/A-O", func(e events.Event) {
		M.showOther = !M.showOther
		M.SetShowBordersR(M.showPanel, M.showOther)
	})

	// probably want to do some self mount first?
	setupEventHandlers(
		M.Panel,
		nil,
		nil,
	)

	// and then refresh?
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

	//// strip off the command name
	//_args, _ := context["args"]
	//args, _ := _args.([]string)
	//if len(args) > 0 && args[0] == "eval" {
	//  args = args[1:]
	//  context["args"] = args
	//}

	//_action, _ := context["action"]
	//action, _ := _action.(string)
	//_index, _ := context["index"]
	//index, _ := _index.(int)
	//_where, _ := context["where"]
	//where, _ := _where.(string)

	//switch action {
	//case "delete":
	//  M.Flex.RemoveIndex(index)
	//case "split":
	//  M.splitPanelItem(action, where, index, context)		

	//case "insert", "scoped":
	//  if flexDir == tview.FlexColumn {
	//    M.insertHorz(action, where, index, context)		
	//  } else {
	//    M.insertVert(action, where, index, context)		
	//  }

	//default:
	//  // no action, probably coming to eval for the first time
	//  M.Flex.AddItem(M.creator(context), 0, 1, true)
	//}

	// only set border when no elements
	// M.Flex.SetBorder(M.Flex.GetItemCount() == 0)
	dir := "row"
	if flexDir == tview.FlexColumn {
		dir = "col"
	}

	M.Flex.SetBorder(true).SetTitle(fmt.Sprintf("  Eval (flex-%s) %s  ", dir, M.Panel.Id()))

	// add the default text if not child elements
	if M.Flex.GetItemCount() == 0 {
		// an initial text element, will want to do better here
		M.Flex.AddItem(M.creator(context), 0, 1, true)
	}

	tui.Draw()

	return nil
}

func (M *Eval) Focus(delegate func(p tview.Primitive)) {
	tui.Log("warn", "Eval.Focus")
	delegate(M.Panel)
	// M.Panel.Focus(delegate)
}
