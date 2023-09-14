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
		"text": C.edit.GetText(),
		"direction": C.GetDirection(),
	}

	m["edit.config"], err = C.editCfg.Encode()
	if err != nil {
		return nil, err
	}

	m["scope.config"], err = C.scope.config.Encode()
	if err != nil {
		return nil, err
	}

	m["scope.viewer"], err = C.scope.viewer.Encode()
	if err != nil {
		return nil, err
	}

	m["final.viewer"], err = C.final.viewer.Encode()
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
		return nil, fmt.Errorf("'typename' missing when calling widget.Box.Decode: %#v", input)
	}

	if tn != C.TypeName() {
		return nil, fmt.Errorf("'typename' mismatch when calling widget.Box.Decode: expected %s, got %s", C.TypeName(), tn)
	}

	text := ""
	if _text, ok := input["text"]; ok {
		text = _text.(string)
	}

	w := New(text)

	ec, ok := input["edit.config"]
	if !ok {
		return nil, fmt.Errorf("scope.config not found in input to Playground.Decode: %#v", input)
	}
	ecfg, err := (&helpers.SourceConfig{}).Decode(ec.(map[string]any))
	if err != nil {
		return nil, err
	}
	w.editCfg = ecfg

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
	w.seeScope = useScope

	w.SetDirection(input["direction"].(int))

	if w.scope.config.Source != helpers.EvalConn {
		w.UseScope(true)
		w.Rebuild(true)
	}

	// hack, first time after loading edge cases
	//if w.scope.config.WatchTime > 0 && w.scope.config.WatchQuit == nil {
		//callback := func() {
			//w.Rebuild(true)
		//}
		//err = w.scope.config.Watch(w.Name(), callback, w.scope.config.WatchTime)
	//}
	//if w.editCfg.WatchTime > 0 && w.editCfg.WatchQuit == nil {
		//callback := func() {
			//txt, err := w.editCfg.GetText()
			//// tui.Log("crit", "got here again:\n"+txt)
			//if err != nil {
				//tui.Log("error", err)
			//}
			//w.SetText(txt)
			//w.Rebuild(false)
		//}
		//err = w.editCfg.Watch(w.Name(), callback, w.editCfg.WatchTime)
	//}
	return w, nil
}
