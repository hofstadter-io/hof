package components

import (
	"fmt"

	"github.com/hofstadter-io/hof/lib/tui/components/widget"

	// cue widgets
	"github.com/hofstadter-io/hof/lib/tui/components/cue/browser"
	"github.com/hofstadter-io/hof/lib/tui/components/cue/playground"
)

var _registry = map[string]func (input map[string]any) (widget.Widget, error){
	// common widgets
	"widget/Box": (widget.NewBox()).Decode,
	"widget/TextView": (widget.NewTextView()).Decode,

	// cue widgets
	"cue/browser": (&browser.Browser{}).Decode,
	"cue/playground": (&playground.Playground{}).Decode,
}

func DecodeWidget(input map[string]any) (widget.Widget, error) {

	typename, ok := input["typename"]
	if !ok {
		return nil, fmt.Errorf("input to DecodeWidget did not contain 'typename'")
	}

	decoder, ok := _registry[typename.(string)]
	if !ok {
		return nil, fmt.Errorf("unknown 'typename': %q in DecodeWidget", typename)
	}

	return decoder(input)
}
