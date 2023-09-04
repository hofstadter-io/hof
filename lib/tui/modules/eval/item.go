package eval

import (
	"fmt"

	"github.com/gdamore/tcell/v2"

	"github.com/hofstadter-io/hof/lib/tui/tview"
)


type Item struct {
	*tview.Frame

	_id   int
	_name string
	_item tview.Primitive

	_parent *Panel	
}
var text_count = 0

func NewItem(item tview.Primitive, parent *Panel) *Item {
	t := &Item{
		_item: item,
		_id: text_count,
		_parent: parent,
	}
	text_count++

	// setup fram with text
	t.Frame = tview.NewFrame(t._item)

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
 
func (t *Item) Item() tview.Primitive {
	return t._item
}

func (t *Item) SetItem(item tview.Primitive) {
	t._item = item
}
 
func (t *Item) Parent() *Panel {
	return t._parent
}

func (t *Item) SetParent(parent *Panel) {
	t._parent = parent
}
 
