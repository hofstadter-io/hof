package flower

import (
	"github.com/hofstadter-io/hof/lib/tui/components/widget"
)


func (F *Flower) Encode() (map[string]any, error) {
	m := make(map[string]any)

	return m, nil
}

func (F *Flower) Decode(map[string]any) (widget.Widget, error) {
	f := New()

	return f, nil
}
