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

	// default overide to all panels
	// would it be better as a widget creator? (after refactor 1)
	// or a function that can take a widget creator with a default ItemBase++
	_creator ItemCreator

	// metadata
	_cnt  int
	_name string
}

var eval_count int

func NewEval() *Eval {
	M := &Eval{
		Panel: NewPanel(nil, nil),
		showPanel: true,
		showOther: true,
		_cnt: eval_count,
		_name: fmt.Sprintf("  Eval %v  ", eval_count),
	}
	eval_count++

	// do layout setup here
	M.Flex.SetDirection(tview.FlexColumn)
	M.Flex.SetBorder(true).SetTitle(M._name)

	return M
}

func (M *Eval) Mount(context map[string]any) error {
	tui.Log("crit", "example message")
	tui.Log("error", "example message")
	tui.Log("warn", "example message")
	tui.Log("info", "example message")
	tui.Log("extra", "example message")
	tui.Log("debug", "example message")
	tui.Log("trace", "example message")

	// this will mount the core element and all children
	M.Flex.Mount(context)
	// tui.Log("trace", "Eval.Mount")

	// probably want to do some self mount first?
	M.setupEventHandlers()

	// and then refresh?
	err := M.Refresh(context)
	if err != nil {
		tui.SendCustomEvent("/console/error", err)
		return err
	}

	return nil
}

func (M *Eval) Unmount() error {
	// remove keybinds
	M.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey { return event })

	// handle border display
	tui.RemoveWidgetHandler(M.Panel, "/sys/key/A-P")
	tui.RemoveWidgetHandler(M.Panel, "/sys/key/A-O")

	// this is where we can do some unloading, depending on the application
	M.Flex.Unmount()

	return nil
}

// todo, add more functions so that we can separate new command messages from refresh?

func (M *Eval) Refresh(context map[string]any) error {
	tui.Log("debug", fmt.Sprintf("Eval.refresh.1: %v", context ))

	// reprocess args, all commands should enter the Eval page first
	// needed for when we come in from the command line first time, or the command box later
	context = enrichContext(context)
	args := []string{}
	if _args, ok := context["args"]; ok {
		args = _args.([]string)
	}
	tui.Log("debug", fmt.Sprintf("Eval.Refresh.2: %v %# v", args, context))

	// handle any top-leval eval commands
	action := ""
	if _action, ok := context["action"]; ok {
		action = _action.(string)
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

	case "list":
		err := M.ListEval()
		if err != nil {
			tui.Tell("error", err)
			tui.Log("error", err)
			return err
		}
		return nil

	}

	// this should go away and be handled in the panel
	// we want Eval to be dumb as bricks
	if M.GetItemCount() == 0 {
		I, err := M.Panel.creator(context, M.Panel)
		if err != nil {
			tui.Log("error", err)	
			return err
		}
		M.AddItem(I, 0, 1, true)
		tui.Draw()
		return nil
	}

	panel := M.GetMostFocusedPanel()
	if panel == nil {
		panel = M.Panel
	}

	return panel.Refresh(context)
}


func (M *Eval) Focus(delegate func(p tview.Primitive)) {
	// tui.Log("warn", "Eval.Focus")
	delegate(M.Panel)
	// M.Panel.Focus(delegate)
}

// This is probably now what Wrap*Handlers is in tview, and Panel might now be a different kind of component, since others embed and extend it
func (M *Eval) setupEventHandlers() {

	//
	// Our message bus system (which probably needs some updating for nested handling
	//

	// handle border display
	tui.AddWidgetHandler(M.Panel, "/sys/key/A-P", func(e events.Event) {
		if M.HasFocus() {
			M.showPanel = !M.showPanel
			M.SetShowBordersR(M.showPanel, M.showOther)
		}
	})

	tui.AddWidgetHandler(M.Panel, "/sys/key/A-O", func(e events.Event) {
		if M.HasFocus() {
			M.showOther = !M.showOther
			M.SetShowBordersR(M.showPanel, M.showOther)
		}
	})

	//
	// tview.Primitive scoped key input handling
	//
	M.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		alt := event.Modifiers() & tcell.ModAlt == tcell.ModAlt
		//ctrl := event.Modifiers() & tcell.ModCtrl == tcell.ModCtrl
		//meta := event.Modifiers() & tcell.ModMeta == tcell.ModMeta
		//shift := event.Modifiers() & tcell.ModShift == tcell.ModShift

		ctx := make(map[string]any)

		// we only care about ALT+... keys at this level
		// tui.Log("trace", fmt.Sprintf("Panel.inputHandler.2 %v %v %v %v %v %q %v", P.Id(), alt, ctrl, meta, shift, string(event.Rune()), event.Key()))
		// tui.Log("warn", fmt.Sprintf("Eval.keyInput %v %v %v", alt, event.Key(), string(event.Rune())))

		panel := M.GetMostFocusedPanel()
		if panel != nil {
			ctx["panel"] = panel 
			ctx["panel-id"] = panel.Id()
			ctx["child-focus-index"] = panel.ChildFocus()
		}

		switch event.Key() {

		// give up focus to parent (this is meh, as it doesn't cross panel bounderies (but maybe easier after refactor?)
		case tcell.KeyESC:
			if panel._parent != nil {
				tui.SetFocus(panel._parent)
			}

		case tcell.KeyRune:
			handled := false
			if alt {
				handled = true
				switch event.Rune() {
				// lowercase = nav
				// upsercase = move/insert

				// left, prev
				case 'h':
					ctx["action"] = "nav.right"
				case 'H':
					ctx["action"] = "move"
					ctx["where"] = "prev"

				// down, prev
				case 'j':
					ctx["action"] = "nav.down"
				case 'J':
					ctx["action"] = "insert"
					ctx["where"] = "prev"
					ctx["item"] = "help"

				// up, next
				case 'k':
					ctx["action"] = "nav.up"
				case 'K':
					ctx["action"] = "insert"
					ctx["where"] = "next"
					ctx["item"] = "help"

				// up, right
				case 'l':
					ctx["action"] = "nav.left"
				case 'L':
					ctx["action"] = "move"
					ctx["where"] = "next"

				default:
					handled = false
				}
			}

			// mid := panel.Id()

			if !handled && alt {
				switch event.Rune() {

				case 'T':
					ctx["action"] = "split"
					ctx["item"] = "help"

				case 'D':
					ctx["action"] = "delete" // DELETE

				// flip flex orientation
				case 'F':
					panel.FlipFlexDirection()
					return nil

				// dev stuff
				case 'v':
					focus := panel.ChildFocus()
					count := panel.GetItemCount()
					tui.Log("trace", fmt.Sprintf("Panel(%s).focus at %v of %v", panel.Id(), focus, count))
					return nil

				default:
					return event

				}
				handled = true
			}

			if handled {
				// ???
				ctx["index"] = panel.ChildFocus()

				// M.Refresh(ctx)
				panel.Refresh(ctx)
				return nil
			}

		}

		return event
	})
}
