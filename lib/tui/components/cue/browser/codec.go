package browser

import (
	"fmt"

	"cuelang.org/go/cue/format"

	"github.com/hofstadter-io/hof/lib/gen"
)

func (V *Browser) Encode() (map[string]any, error) {
	m := map[string]any{
		"type": V.TypeName(),
		"mode": V.mode,
		"usingScope": V.usingScope,
		"docs": V.docs,
		"attrs": V.attrs,
		"defs": V.defs,
		"optional": V.optional,
		"ignore": V.ignore,
		"inline": V.inline,
		"resolve": V.resolve,
		"concrete": V.concrete,
		"hidden": V.hidden,
		"final": V.final,
		"validate": V.validate,
	}

	var err error
	m["source"], err = V.source.Encode()
	return m, err
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
