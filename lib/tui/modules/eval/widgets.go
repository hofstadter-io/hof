package eval

import (
	"github.com/hofstadter-io/hof/lib/tui/tview"
)

// base and wrapped tview widgets, temporarily here

type Widget interface {
	tview.Primitive

	TypeName() string

	EncodeMap() (map[string]any, error)
}

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

func (*Box) EncodeMap() (map[string]any, error) {
	return map[string]any{
		"type": "Box",
	}, nil
}


type TextView struct {
	*tview.TextView
	Base
}

func (*TextView) EncodeMap() (map[string]any, error) {
	return map[string]any{
		"type": "TextView",
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
