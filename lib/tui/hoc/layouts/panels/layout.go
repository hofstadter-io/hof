// Panels Layout is a Flex widget with
// hidable panels and a main content.
// Can be vert or horz oriendted throught the Flex widget.
// Panels can be jumped to with <focuskey> and hidden with <shift>-<focuskey>
// Recommend making the <focuskey> an: '<alt>-<key>' and hidden will be '<shift>-<alt>-<key>'
// Normal movement and interaction keys within the focussed panel.
//
// main (middle) panel, can be anything, including...
// - the router (when this is the root view)
// - another DashAndPanels
// - a pager, grid, or any other primitive
package panels

import (
	"sort"

	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/events"
	"github.com/hofstadter-io/hof/lib/tui/tview"
)

type Panel struct {
	*tview.Flex

	Name       string
	Item       tview.Primitive
	FixedSize  int
	Proportion int
	Focus      int
	FocusKey   string
	Hidden     bool
	HiddenKey  string
}

type Layout struct {
	*tview.Flex

	// first (left/top) panels, can be almost anything and hidden.
	fPanels map[string]*Panel

	// main (middle) panel, can be anything, I think.
	mPanel *Panel

	// last (right/bottom) panels, can be almost anything and hidden.
	lPanels map[string]*Panel
}

func New() *Layout {
	L := &Layout{
		Flex:    tview.NewFlex(),
		fPanels: map[string]*Panel{},
		lPanels: map[string]*Panel{},
	}

	return L
}

// AddFirstPanel adds a Panel to the left or top, depending on orientation.
func (L *Layout) AddFirstPanel(name string, item tview.Primitive, fixedSize, proportion,
	focus int, focuskey string, hidden bool, hiddenkey string) {
	panel := &Panel{
		Flex:       tview.NewFlex(),
		Name:       name,
		Item:       item,
		FixedSize:  fixedSize,
		Proportion: proportion,
		Focus:      focus,
		FocusKey:   focuskey,
		Hidden:     hidden,
		HiddenKey:  hiddenkey,
	}

	L.fPanels[name] = panel
}

// AddLastPanel adds a Panel to the right or bottom, depending on orientation.
func (L *Layout) AddLastPanel(name string, item tview.Primitive, fixedSize, proportion,
	focus int, focuskey string, hidden bool, hiddenkey string) {
	panel := &Panel{
		Flex:       tview.NewFlex(),
		Name:       name,
		Item:       item,
		FixedSize:  fixedSize,
		Proportion: proportion,
		Focus:      focus,
		FocusKey:   focuskey,
		Hidden:     hidden,
		HiddenKey:  hiddenkey,
	}

	L.lPanels[name] = panel
}

func (L *Layout) SetMainPanel(name string, item tview.Primitive, fixedSize, proportion, focus int, focuskey string) {
	panel := &Panel{
		Flex:       tview.NewFlex(),
		Name:       name,
		Item:       item,
		FixedSize:  fixedSize,
		Proportion: proportion,
		Focus:      focus,
		FocusKey:   focuskey,
	}

	panel.Flex.AddItem(item, 0, 1, false)
	L.mPanel = panel
}

func (L *Layout) Mount(context map[string]interface{}) error {
	err := L.build()
	if err != nil {
		return err
	}

	// Setup focuskeys
	for _, panel := range L.fPanels {
		panel.Item.Mount(context)
		if panel.FocusKey != "" {
			localPanel := panel
			tui.AddWidgetHandler(L, "/sys/key/"+localPanel.FocusKey, func(e events.Event) {
				// go tui.SendCustomEvent("/console/trace", "Focus: "+localPanel.Name)
				tui.SetFocus(localPanel.Item)
			})
		}
		if panel.HiddenKey != "" {
			localPanel := panel
			tui.AddWidgetHandler(L, "/sys/key/"+localPanel.HiddenKey, func(e events.Event) {
				localPanel.Hidden = !localPanel.Hidden
				// go tui.SendCustomEvent("/console/trace", fmt.Sprintf("Hidden: %s (%v)", localPanel.Name, localPanel.Hidden))
				L.build()
				if localPanel.Hidden {
					tui.SetFocus(L.mPanel.Item)
				} else {
					tui.SetFocus(localPanel.Item)
				}
				tui.Draw()
			})
		}
	}
	if L.mPanel.FocusKey != "" {

		L.mPanel.Item.Mount(context)
		localPanel := L.mPanel
		tui.AddWidgetHandler(L, "/sys/key/"+localPanel.FocusKey, func(e events.Event) {
			// go tui.SendCustomEvent("/console/trace", "Focus: "+localPanel.Name)
			tui.SetFocus(localPanel.Item)
		})
	}
	for _, panel := range L.lPanels {
		panel.Item.Mount(context)
		if panel.FocusKey != "" {
			localPanel := panel
			tui.AddWidgetHandler(L, "/sys/key/"+localPanel.FocusKey, func(e events.Event) {
				// go tui.SendCustomEvent("/console/trace", "Focus: "+localPanel.Name)
				tui.SetFocus(localPanel.Item)
			})
		}
		if panel.HiddenKey != "" {
			localPanel := panel
			tui.AddWidgetHandler(L, "/sys/key/"+localPanel.HiddenKey, func(e events.Event) {
				localPanel.Hidden = !localPanel.Hidden
				// go tui.SendCustomEvent("/console/trace", fmt.Sprintf("Hidden: %s (%v)", localPanel.Name, localPanel.Hidden))
				L.build()
				if localPanel.Hidden {
					tui.SetFocus(L.mPanel.Item)
				} else {
					tui.SetFocus(localPanel.Item)
				}
				tui.Draw()
			})
		}
	}

	return nil
}

// puke.... we call this far too often, probably the source of the ordering issue
func (L *Layout) build() error {
	// get and order the fPanels
	fPs := []*Panel{}
	for _, panel := range L.fPanels {
		if panel.Hidden {
			continue
		}
		fPs = append(fPs, panel)
	}
	sort.Slice(fPs, func(i, j int) bool {
		return fPs[i].Focus < fPs[j].Focus
	})

	// get and order the lPanels
	lPs := []*Panel{}
	for _, panel := range L.lPanels {
		if panel.Hidden {
			continue
		}
		lPs = append(lPs, panel)
	}
	sort.Slice(lPs, func(i, j int) bool {
		return lPs[i].Focus < lPs[j].Focus
	})

	// Start a fresh Flex item
	orient := L.GetDirection()
	L.Flex = tview.NewFlex().SetDirection(orient)

	for _, p := range fPs {
		L.AddItem(p.Item, p.FixedSize, p.Proportion, false)
	}

	p := L.mPanel
	L.AddItem(p.Item, p.FixedSize, p.Proportion, true)

	for _, p := range lPs {
		L.AddItem(p.Item, p.FixedSize, p.Proportion, false)
	}

	return nil
}
