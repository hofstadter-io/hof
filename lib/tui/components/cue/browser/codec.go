package browser

import (
	"fmt"

	"cuelang.org/go/cue/format"

	"github.com/hofstadter-io/hof/lib/tui/components/cue/helpers"
	"github.com/hofstadter-io/hof/lib/tui/components/widget"
	"github.com/hofstadter-io/hof/lib/gen"
)

func (W *Browser) Encode() (map[string]any, error) {
	m := map[string]any{
		"typename": W.TypeName(),
		"mode": W.mode,
		"usingScope": W.usingScope,
		"docs": W.docs,
		"attrs": W.attrs,
		"defs": W.defs,
		"optional": W.optional,
		"ignore": W.ignore,
		"inline": W.inline,
		"resolve": W.resolve,
		"concrete": W.concrete,
		"hidden": W.hidden,
		"final": W.final,
		"validate": W.validate,
	}

	var err error
	m["source"], err = W.source.Encode()
	return m, err
}

func (W *Browser)	Decode(input map[string]any) (widget.Widget, error) {
	s, ok := input["source"]
	if !ok {
		return nil, fmt.Errorf("source not found in input to Browser.Decode: %#v", input)
	}
	sc, err := (&helpers.SourceConfig{}).Decode(s.(map[string]any))
	if err != nil {
		return nil, err
	}

	w := New(sc, "cue")
	w.mode = input["mode"].(string)
	w.usingScope = input["usingScope"].(bool)
	w.docs = input["docs"].(bool)
	w.attrs = input["attrs"].(bool)
	w.defs = input["defs"].(bool)
	w.optional = input["optional"].(bool)
	w.ignore = input["ignore"].(bool)
	w.inline = input["inline"].(bool)
	w.resolve = input["resolve"].(bool)
	w.concrete = input["concrete"].(bool)
	w.hidden = input["hidden"].(bool)
	w.final = input["final"].(bool)
	w.validate = input["validate"].(bool)

	return w, nil
}

func (C *Browser) GetValueText(mode string) (string, error) {
	var (
		b []byte
		err error
	)
	switch mode {
	case "cue":
		syn := C.value.Syntax(C.Options()...)

		b, err = format.Node(syn)
		if !C.ignore {
			if err != nil {
				return "", err
			}
		}

	case "json":
		f := &gen.File{}
		b, err = f.FormatData(C.value, "json")
		if err != nil {
			return "", err
		}

	case "yaml":
		f := &gen.File{}
		b, err = f.FormatData(C.value, "yaml")
		if err != nil {
			return "", err
		}

	default:
		return "", fmt.Errorf("unknown file type %q", mode)

	}

	return string(b), err
}
