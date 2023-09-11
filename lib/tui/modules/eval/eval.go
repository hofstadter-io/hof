package eval

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"

	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/components/panel"
	"github.com/hofstadter-io/hof/lib/tui/components/widget"
	"github.com/hofstadter-io/hof/lib/tui/tview"
)

type Eval struct {
	*panel.Panel

	// border display
	showPanel, showOther bool

	// default overide to all panels
	// would it be better as a widget creator? (after refactor 1)
	// or a function that can take a widget creator with a default ItemBase++
	_creator panel.ItemCreator

	// metadata
	_name string
}

func NewEval() *Eval {
	M := &Eval{
		showPanel: true,
		showOther: true,
		_name: fmt.Sprintf("  Eval  "),
	}
	M.Panel = panel.New(nil, M.creator)

	item, _ := helpItem(nil, M.Panel)
	M.Panel.AddItem(item, 0, 1, true)

	// do layout setup here
	M.Flex.SetDirection(tview.FlexColumn)
	M.Flex.SetBorder(true).SetTitle(M._name)

	return M
}

func (M *Eval) Mount(context map[string]any) error {

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

func (M *Eval) showError(err error) error {
	txt := widget.NewTextView()
	fmt.Fprint(txt, err)

	I := panel.NewBaseItem(nil, M.Panel)
	I.SetWidget(txt)

	M.Panel.AddItem(I, 0, 1, true)

	return err
}



func (M *Eval) Focus(delegate func(p tview.Primitive)) {
	// tui.Log("warn", "Eval.Focus")
	delegate(M.Panel)
	// M.Panel.Focus(delegate)
}

func (M *Eval) getItemByPath(path string) (panel.PanelItem, error) {
	parts := strings.Split(path, ".")

	// set at our panel
	curr := M.Panel

	for _, part := range parts {
		p := curr.GetItemByName(part)
		if p == nil {
			p = curr.GetItemById(part)
			if p == nil {
				return nil, fmt.Errorf("unable to find node %q in %q", part, path)
			}
		}
		switch t := p.(type) {
		case *panel.Panel:
			curr = t	
		case panel.PanelItem:
			return t, nil
		}
	}

	return nil, fmt.Errorf("did not find item at path %q", path)
}
