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
	// cue  (eventually own component) (these are more like temp holders during processing of inputs
	//

	// scope
	_scopeR *runtime.Runtime
	_scopeV cue.Value
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

	// this needs to move down into the widget
	m["context"] = I._context

	m["widget"], err = I._widget.EncodeMap()
	if err != nil {
		return m, err
	}

	return m, nil
}

func ItemDecodeMap(data map[string]any, parent *Panel) (*Item, error) {
	I := &Item{
		_parent: parent,
		_widget: NewBox(),
		_cnt: data["id"].(int),
		_name: data["name"].(string),
	}

	// setup frame with temp box
	I.Frame = tview.NewFrame(I._widget)

	// style fram
	I.SetBorders(0,0,0,0,0,0) // just the one-line header
	txt := fmt.Sprintf(" %s ", I.Id())
	I.AddText(txt, true, tview.AlignLeft, tcell.ColorLimeGreen)
	I.SetBorder(true)

	var context map[string]any

	if c, ok := data["context"]; ok {
		context = c.(map[string]any)
	} else {
		return I, fmt.Errorf("context config not found in item: %# v", data)
	}

	i, err := parent.creator(context, parent)
	if err != nil {
		return i, err
	}

	//var wdata map[string]any
	//if w, ok := data["widget"]; ok {
	//  wdata = w.(map[string]any)
	//} else {
	//  return I, fmt.Errorf("widget config not found in item: %# v", data)
	//}

	return i, nil
}
