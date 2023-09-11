package panel

import (
	"fmt"

	"github.com/hofstadter-io/hof/lib/tui/components/widget"
	"github.com/hofstadter-io/hof/lib/tui/tview"
)

// the context format is determined by the user of Panel
// you will need to construct this from user input and then process it in items
// Panel will handle moving it to the right Item, or use the ItemCreator to make a new one
type ItemContext map[string]any


type PanelItem interface {
	tview.Primitive

	// some functions we need to add to primitive
	SetBorder(show bool) *tview.Box
	SetTitle(title string) *tview.Box

	// some metadata
	Name() string
	SetName(string)

	Parent() *Panel
	SetParent(*Panel)

	Widget() widget.Widget
	SetWidget(widget.Widget)

	// this is how you can pass data built up in top-level pages, commands, and keybindings
	// and then modify or rebuild the component locally, you decide how to refresh
	Handle(context map[string]any) (handled bool, err error)

	// should rebuild the item from the current or latest context
	Rebuild() error

	// Encode should turn the PanelItem into ItemContext.
	// It should contain enough detail to rebuild the item
	// after serialization to JSON and back.
	Encode() (ItemContext, error)

	// Decode should create a PanelItem from an ItemContext
	Decode(context ItemContext, creator ItemCreator, parent *Panel) (PanelItem, error)
}

// This is a function that builds a new PanelItem from an ItemContext
type ItemCreator func(context ItemContext, parent *Panel) (PanelItem, error)

// BaseItem can be used as a starting point for your custom PanelItems.
// It wraps a Widget in some metadata and a tview.Frame
type BaseItem struct {
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
	_widget widget.Widget
}

// need a better way to do this, but uuid/cuid is a bit much
var item_count = 0

func NewBaseItem(context map[string]any, parent *Panel) *BaseItem {
	t := &BaseItem{
		_cnt: item_count,
		_parent: parent,
		_widget: widget.NewBox(),
	}
	item_count++

	// setup frame with temp box
	t.Frame = tview.NewFrame(t._widget)

	// style fram
	t.SetBorders(0,0,0,0,0,0) // just the one-line header
	txt := fmt.Sprintf(" %s ", t.Id())
	t.SetTitle(txt).SetTitleAlign(tview.AlignLeft)
	// t.AddText(txt, true, tview.AlignLeft, tcell.ColorLimeGreen)
	t.SetBorder(true)
	return t
}

func (I *BaseItem) Id() string {
	return fmt.Sprintf("b%d", I._cnt)
}

func (I *BaseItem) Name() string {
	return I._name
}

func (I *BaseItem) SetName(name string) {
	I._name = name
}
 
func (I *BaseItem) Parent() *Panel {
	return I._parent
}

func (I *BaseItem) SetParent(parent *Panel) {
	I._parent = parent
} 

func (I *BaseItem) Widget() widget.Widget {
	return I._widget
}

func (I *BaseItem) SetWidget(w widget.Widget) {
	I._widget = w
	I.Frame.SetPrimitive(I._widget)
}

	// this is how you can pass data built up in top-level pages, commands, and keybindings
	// and then modify or rebuild the component locally, you decide how to refresh
func (I *BaseItem) Handle(context map[string]any) (handled bool, err error) {
	return false, nil
}

	// should rebuild the item from the current or latest context
func (I *BaseItem) Rebuild() error {
	return nil
}
 
func (I *BaseItem) Encode() (ItemContext, error) {
	var err error
	m := make(map[string]any)

	m["id"] = I._cnt
	m["name"] = I._name
	m["type"] = "item"

	m["widget"], err = I._widget.Encode()
	if err != nil {
		return m, err
	}

	return m, nil
}

func (I *BaseItem) Decode(context ItemContext, creator ItemCreator, parent *Panel) (PanelItem, error) {
	//I := &BaseItem{
	//  _widget: widget.NewBox(),
	//  _cnt: data["id"].(int),
	//  _name: data["name"].(string),
	//}

	//// setup frame with temp box
	//I.Frame = tview.NewFrame(I._widget)

	//// style fram
	//I.SetBorders(0,0,0,0,0,0) // just the one-line header
	//txt := fmt.Sprintf(" %s ", I.Id())
	//I.AddText(txt, true, tview.AlignLeft, tcell.ColorLimeGreen)
	//I.SetBorder(true)

	//var context map[string]any

	//if c, ok := data["widget"]; ok {
	//  context = c.(map[string]any)
	//} else {
	//  return I, fmt.Errorf("context config not found in item: %# v", data)
	//}

	//i, err := creator(context)
	//if err != nil {
	//  return i, err
	//}

	////var wdata map[string]any
	////if w, ok := data["widget"]; ok {
	////  wdata = w.(map[string]any)
	////} else {
	////  return I, fmt.Errorf("widget config not found in item: %# v", data)
	////}

	return nil, nil
}
