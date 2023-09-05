package eval

import (
	"fmt"

	"cuelang.org/go/cue"
	"github.com/gdamore/tcell/v2"

	"github.com/hofstadter-io/hof/lib/runtime"
	"github.com/hofstadter-io/hof/lib/tui/tview"
)

type ItemCreator func(context map[string]any, panel *Panel) (*Item, error)

// item will become an interface, for now, it is just the 

type Item struct {
	// base item should be a frame with settings for
	//  - the header dropdown
	//  - parent/widget
	//  - save/load
	*tview.Frame

	//
	// there is a base component in here
	//

	// meta
	_id   int
	_name string

	// tui
	_parent *Panel	
	_widget tview.Primitive  // Do we eventually build a iface & base type for this? Grafana for inspiration?

	// params++ this item was created with
	_context map[string]any

	//
	// unsure where these go yet, Item or Widget
	//

	// i/o connection (format & type tbd)
	_conns map[string]any

	//
	// cue  (eventually own component)
	//

	// scope
	//_scopeV cue.Value
	//_scopeR *runtime.Runtime
	// (note, cannot have 2 runtimes yet)

	// current
	_runtime *runtime.Runtime
	_value   cue.Value
	_text    string

	// final value
	_final cue.Value

}

// need a better way to do this, but uuid/cuid is a bit much
var item_count = 0

func NewItem(context map[string]any, parent *Panel) *Item {
	t := &Item{
		_id: item_count,
		_parent: parent,
		_widget: tview.NewBox(),
	}
	item_count++

	// setup frame with temp box
	t.Frame = tview.NewFrame(t._widget)

	// style fram
	t.SetBorders(0,0,0,0,0,0) // just the one-line header
	txt := fmt.Sprintf(" ☰  %s", t.Id())
	t.AddText(txt, true, tview.AlignLeft, tcell.ColorLimeGreen)
	t.SetBorder(true)
	return t
}

func (t *Item) Id() string {
	return fmt.Sprintf("t:%d", t._id)
}

func (t *Item) Name() string {
	return t._name
}

func (t *Item) SetName(name string) {
	t._name = name
}
 
func (t *Item) Widget() tview.Primitive {
	return t._widget
}

func (t *Item) SetWidget(widget tview.Primitive) {
	t._widget = widget
	t.Frame.SetPrimitive(t._widget)
}
 
func (t *Item) Parent() *Panel {
	return t._parent
}

func (t *Item) SetParent(parent *Panel) {
	t._parent = parent
} 
