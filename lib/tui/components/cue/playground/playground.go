package playground

import (
	"time"

	"cuelang.org/go/cue"

	"github.com/hofstadter-io/hof/lib/tui/components/cue/browser"
	"github.com/hofstadter-io/hof/lib/tui/components/cue/helpers"
	"github.com/hofstadter-io/hof/lib/tui/tview"
	"github.com/hofstadter-io/hof/lib/watch"
)

type valPack struct {
	config  *helpers.SourceConfig
	value   cue.Value
	viewer  *browser.Browser // scope
}

type Playground struct {
	// *tview.Frame eventually?
	*tview.Flex

	// scope used during parsing / evaluation
	useScope bool
	scope    *valPack

	// the editor box
	text   string    // text entered by the user to make the final value
	edit *tview.TextArea  // text

	// the final value
	final    *valPack

	// that's funky!
	debouncer func(func())
}

func (*Playground) TypeName() string {
	return "cue/playground"
}

func (V *Playground) Encode() (map[string]any, error) {
	var err error
	m := map[string]any{
		"type": V.TypeName(),
		"useScope": V.useScope,
		"text": V.text,
	}

	m["scope.config"], err = V.scope.config.Encode()
	if err != nil {
		return m, err
	}

	m["scope.viewer"], err = V.scope.viewer.Encode()
	if err != nil {
		return m, err
	}

	m["final.viewer"], err = V.final.viewer.Encode()
	if err != nil {
		return m, err
	}

	return m, nil
}


func New(initialText string) (*Playground) {

	C := &Playground{
		Flex: tview.NewFlex(),
		text: initialText,
		scope: &valPack{},
		final: &valPack{},
	}
	// our wrapper around the CUE widgets
	C.Flex = tview.NewFlex().SetDirection(tview.FlexColumn)

	// TODO, options form

	// scope viewer
	C.scope.config = &helpers.SourceConfig{}
	C.scope.viewer = browser.New(C.scope.config, "cue")
	C.scope.viewer.SetName("scope")
	C.scope.viewer.SetBorder(true)

	// curr editor
	C.edit = tview.NewTextArea()
	C.edit.
		SetTitle("  expression(s)  ").
		SetBorder(true)

	C.edit.SetText(C.text, false)

	// results viewer
	C.final.config = &helpers.SourceConfig{}
	C.final.viewer = browser.New(C.final.config, "cue")
	C.final.viewer.SetName("result")
	C.final.viewer.SetBorder(true)

	// usingScope?
	C.final.viewer.SetUsingScope(false)
	C.useScope = false

	// layout
	C.Flex.
		AddItem(C.scope.viewer, 0, 1, true).
		AddItem(C.edit, 0, 1, true).
		AddItem(C.final.viewer, 0, 1, true)

	C.setupKeybinds()

	// setup change response with douncer
	// to trigger rebuild on editor changes
	C.debouncer = watch.NewDebouncer(time.Millisecond * 333)
	C.edit.SetChangedFunc(func() {
		C.debouncer(func(){
			C.Rebuild(false)
		})
	})
	return C
}

func (C *Playground) SetText(s string) {
	C.edit.SetText(s, false)
}

func (C *Playground) SetScopeConfig(sc *helpers.SourceConfig) {
	C.scope.config = sc
	C.scope.viewer.SetSourceConfig(sc)
}

func (C *Playground) UseScope(visible bool) {
	C.useScope = visible
}

func (C *Playground) SetFlexDirection(dir int) {
	C.SetDirection(dir)
}


func (C *Playground) Mount(context map[string]any) error {

	return nil
}

func (C *Playground) Focus(delegate func(p tview.Primitive)) {
	if C.final.viewer.HasFocus() {
		delegate(C.final.viewer)
		return
	}
	// otherwise, assume we want to keep the view focus
	delegate(C.edit)
	return
}


