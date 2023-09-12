package widget

import (
	"fmt"

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
			_typename: "widget/Box",
		},
	}
}

func (W *Box) Encode() (map[string]any, error) {
	return map[string]any{
		"typename": W.TypeName(),
	}, nil
}

func (W *Box)	Decode(input map[string]any) (Widget, error) {
	tn, ok := input["typename"]
	if !ok {
		return nil, fmt.Errorf("'typename' missing when calling widget.Box.Decode: %#v", input)
	}

	if tn != W.TypeName() {
		return nil, fmt.Errorf("'typename' mismatch when calling widget.Box.Decode: expected %s, got %s", W.TypeName(), tn)
	}

	return NewBox(), nil
}

type TextView struct {
	*tview.TextView
	Base
}

func NewTextView() *TextView {
	t := &TextView{
		TextView: tview.NewTextView(),
		Base: Base{
			_typename: "widget/TextView",
		},
	}

	// default settings
	t.SetDynamicColors(true)

	return t
}

func (W *TextView) Encode() (map[string]any, error) {
	return map[string]any{
		"typename": W.TypeName(),
		"text": W.GetText(false),
	}, nil
}

func (W *TextView)	Decode(input map[string]any) (Widget, error) {
	tn, ok := input["typename"]
	if !ok {
		return nil, fmt.Errorf("'typename' missing when calling widget.Box.Decode: %#v", input)
	}

	if tn != W.TypeName() {
		return nil, fmt.Errorf("'typename' mismatch when calling widget.Box.Decode: expected %s, got %s", W.TypeName(), tn)
	}

	text := ""
	if _text, ok := input["text"]; ok {
		text = _text.(string)
	}

	w := NewTextView()
	w.SetText(text)

	return w, nil
}
