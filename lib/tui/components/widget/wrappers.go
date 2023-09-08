package widget

import (
	"github.com/hofstadter-io/hof/lib/tui/tview"
)

// base and wrapped tview widgets, temporarily here

type Base struct {
	_typename string
}

func (W Base) TypeName() string {
	return W._typename
}

type Box struct {
	*tview.Box
	Base
}

func NewBox() *Box {
	return &Box{
		Box: tview.NewBox(),
		Base: Base{
			_typename: "Box",
		},
	}
}

func (*Box) Encode() (map[string]any, error) {
	return map[string]any{
		"type": "Box",
	}, nil
}


type TextView struct {
	*tview.TextView
	Base
}

func (T *TextView) Encode() (map[string]any, error) {
	return map[string]any{
		"type": "TextView",
		"text": T.GetText(false),
	}, nil
}

func NewTextView() *TextView {
	t := &TextView{
		TextView: tview.NewTextView(),
		Base: Base{
			_typename: "TextView",
		},
	}

	// default settings
	t.SetDynamicColors(true)

	return t
}

func WidgetDecodeMap(data map[string]any) (Widget, error) {


	return nil, nil
}

