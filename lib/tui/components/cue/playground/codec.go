package playground

import (
	"fmt"

	// "github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/components/cue/browser"
	"github.com/hofstadter-io/hof/lib/tui/components/cue/helpers"
	"github.com/hofstadter-io/hof/lib/tui/components/widget"
)

func (C *Playground) Encode() (map[string]any, error) {
	var err error
	m := map[string]any{
		"typename": C.TypeName(),
		"useScope": C.useScope,
		"seeScope": C.seeScope,
		"mode": C.mode,
		"text": C.edit.GetText(),
		"direction": C.GetDirection(),
	}

	m["scope"], err = C.scope.Encode()
	if err != nil {
		return nil, err
	}

	m["edit"], err = C.editCfg.Encode()
	if err != nil {
		return nil, err
	}

	m["final"], err = C.final.Encode()
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (C *Playground)	Decode(input map[string]any) (widget.Widget, error) {
	// tui.Log("extra", fmt.Sprintf("Play.Decode: %# v", input))
	var err error

	tn, ok := input["typename"]
	if !ok {
		return nil, fmt.Errorf("'typename' missing when calling playground.Decode: %#v", input)
	}

	if tn != C.TypeName() {
		return nil, fmt.Errorf("'typename' mismatch when calling playground.Decode: expected %s, got %s", C.TypeName(), tn)
	}

	text := ""
	if _text, ok := input["text"]; ok {
		text = _text.(string)
	}

	w := New(text)

	ec, ok := input["edit"]
	if !ok {
		return nil, fmt.Errorf("scope.config not found in input to Playground.Decode: %#v", input)
	}
	ecfg, err := (&helpers.SourceConfig{}).Decode(ec.(map[string]any))
	if err != nil {
		return nil, err
	}
	w.editCfg = ecfg

	sv, ok := input["scope"]
	if !ok {
		return nil, fmt.Errorf("scope not found in input to Playground.Decode: %#v", input)
	}
	b, err := (&browser.Browser{}).Decode(sv.(map[string]any))
	if err != nil {
		return nil, err
	}
	// todo, setter so we can setup to share scope.config with viewer.config automatically
	w.scope = b.(*browser.Browser)

	fv, ok := input["final"]
	if !ok {
		return nil, fmt.Errorf("final.viewer not found in input to Playground.Decode: %#v", input)
	}
	b, err = (&browser.Browser{}).Decode(fv.(map[string]any))
	if err != nil {
		return nil, err
	}
	w.final = b.(*browser.Browser)
	w.SetItem(2, w.final, 0, 1, true)

	s, ok := input["useScope"]
	if !ok {
		return nil, fmt.Errorf("'useScope' missing when calling playground.Decode: %#v", input)
	}
	useScope := s.(bool)
	w.UseScope(useScope)

	s, ok = input["seeScope"]
	if !ok {
		return nil, fmt.Errorf("'seeScope' missing when calling playground.Decode: %#v", input)
	}
	w.seeScope = s.(bool)

	s, ok = input["mode"]
	if !ok {
		return nil, fmt.Errorf("'mode' missing when calling playground.Decode: %#v", input)
	}
	w.mode = PlayMode(s.(string))

	w.SetDirection(input["direction"].(int))

	if useScope {
		w.scope.RebuildValue()
		w.scope.Rebuild()
	}
	w.Rebuild()

	return w, nil
}
