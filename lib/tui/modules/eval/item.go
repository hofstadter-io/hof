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
	_cnt   int
	_name string

	// tui
	_parent *Panel	
	_widget Widget  // Do we eventually build a iface & base type for this? Grafana for inspiration?

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

	// args for en|decode (probably just one of them)
	_runtimeArgs []string
	_valueArgs []string

	// final value
	_final cue.Value

}

// need a better way to do this, but uuid/cuid is a bit much
var item_count = 0

func NewItem(context map[string]any, parent *Panel) *Item {
	t := &Item{
		_cnt: item_count,
		_parent: parent,
		_widget: NewBox(),
		_context: context,
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

func (I *Item) Id() string {
	return fmt.Sprintf("t:%d", I._cnt)
}

func (I *Item) Name() string {
	return I._name
}

func (I *Item) SetName(name string) {
	I._name = name
}
 
func (I *Item) Widget() tview.Primitive {
	return I._widget
}

func (I *Item) SetWidget(widget Widget) {
	I._widget = widget
	I.Frame.SetPrimitive(I._widget)
}
 
func (I *Item) Parent() *Panel {
	return I._parent
}

func (I *Item) SetParent(parent *Panel) {
	I._parent = parent
} 

func (I *Item) EncodeMap() (map[string]any, error) {
	var err error
	m := make(map[string]any)

	m["id"] = I._cnt
	m["name"] = I._name
	m["type"] = "item"
	m["context"] = I._context

	m["widget"], err = I._widget.EncodeMap()
	if err != nil {
		return m, err
	}

	return m, nil
}

func ItemDecodeMap(data map[string]any, parent *Panel, creator ItemCreator) (*Panel, error) {
	P := &Panel{
		Flex: tview.NewFlex(),
		_creator: creator,
		_parent: parent,
	}

	return P, nil
}
