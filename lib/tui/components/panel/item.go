package panel

import (
	"fmt"
	"sync/atomic"

	// "github.com/gdamore/tcell/v2"

	"github.com/hofstadter-io/hof/lib/tui/components"
	"github.com/hofstadter-io/hof/lib/tui/components/cue/browser"
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
	// Decode(data ItemContext, parent *Panel) (PanelItem, error)
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
// there are issues on reload potentially, and creating new elements after
var item_count *atomic.Int64
func init() {
	item_count = new(atomic.Int64)
}

func NewBaseItem(parent *Panel) *BaseItem {
	t := &BaseItem{
		Frame: tview.NewFrame(),
		_cnt: int(item_count.Add(1)),
		_parent: parent,
	}

	// setup frame with temp box
	t.Frame.SetPrimitive(widget.NewBox())

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
	m["typename"] = "baseitem"

	m["widget"], err = I._widget.Encode()
	if err != nil {
		return m, err
	}

	return m, nil
}

// dummy
func (I *BaseItem) Decode(map[string]any) (widget.Widget, error) {
	return nil, nil
}

func ItemDecodeMap(data map[string]any, parent *Panel, creator ItemCreator) (PanelItem, error) {
	var err error
	id, _ := data["id"].(int)
	name, _ := data["name"].(string)
	wmap, _ := data["widget"].(map[string]any)

	// typename should always be 'baseitem' for now
	typename, _ := data["typename"].(string)
	var w widget.Widget

	switch typename {
	case "baseitem":
		w, err = components.DecodeWidget(wmap)
		switch t := w.(type) {
		case *browser.Browser:
			t.RebuildValue()
			t.Rebuild()
		}
	case "panel":
		w, err = PanelDecodeMap(wmap, parent, creator)
	}
	if err != nil {
		return nil, err
	}

	// todo, handle PanelItem here, probably need a registry here?
	// for now, just create a BaseItem

	I := &BaseItem{
		Frame: tview.NewFrame(),
		_widget: w,
		_cnt: id,
		_name: name,
	}

	// setup frame with temp box
	I.Frame.SetPrimitive(I._widget)

	// style frame
	I.SetBorders(0,0,0,0,0,0) // just the one-line header
	txt := fmt.Sprintf(" %s ", I.Id())
	I.SetTitle(txt).SetTitleAlign(tview.AlignLeft)
	// t.AddText(txt, true, tview.AlignLeft, tcell.ColorLimeGreen)
	I.SetBorder(true)

	return I, nil
}
