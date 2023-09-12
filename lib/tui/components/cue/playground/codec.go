package playground

import (
	"fmt"

	"github.com/hofstadter-io/hof/lib/tui/components/cue/browser"
	"github.com/hofstadter-io/hof/lib/tui/components/cue/helpers"
	"github.com/hofstadter-io/hof/lib/tui/components/widget"
)

func (W *Playground) Encode() (map[string]any, error) {
	var err error
	m := map[string]any{
		"typename": W.TypeName(),
		"useScope": W.useScope,
		"text": W.edit.GetText(),
		"direction": W.GetDirection(),
	}

	m["scope.config"], err = W.scope.config.Encode()
	if err != nil {
		return nil, err
	}

	m["scope.viewer"], err = W.scope.viewer.Encode()
	if err != nil {
		return nil, err
	}

	m["final.viewer"], err = W.final.viewer.Encode()
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (W *Playground)	Decode(input map[string]any) (widget.Widget, error) {
	// tui.Log("extra", fmt.Sprintf("Play.Decode: %# v", input))
	var err error

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

	w := New(text)

	sc, ok := input["scope.config"]
	if !ok {
		return nil, fmt.Errorf("scope.config not found in input to Playground.Decode: %#v", input)
	}
	scfg, err := (&helpers.SourceConfig{}).Decode(sc.(map[string]any))
	if err != nil {
		return nil, err
	}
	w.SetScopeConfig(scfg)

	sv, ok := input["scope.viewer"]
	if !ok {
		return nil, fmt.Errorf("scope.viewer not found in input to Playground.Decode: %#v", input)
	}
	b, err := (&browser.Browser{}).Decode(sv.(map[string]any))
	if err != nil {
		return nil, err
	}
	w.scope.viewer = b.(*browser.Browser)

	fv, ok := input["final.viewer"]
	if !ok {
		return nil, fmt.Errorf("final.viewer not found in input to Playground.Decode: %#v", input)
	}
	b, err = (&browser.Browser{}).Decode(fv.(map[string]any))
	if err != nil {
		return nil, err
	}
	w.final.viewer = b.(*browser.Browser)
	w.SetItem(2, w.final.viewer, 0, 1, true)

	s, ok := input["useScope"]
	if !ok {
		return nil, fmt.Errorf("'typename' missing when calling widget.Box.Decode: %#v", input)
	}
	useScope := s.(bool)
	w.UseScope(useScope)

	w.SetDirection(input["direction"].(int))

	if w.scope.config.Source != helpers.EvalConn {
		w.Rebuild(true)
	}

	return w, nil
}
