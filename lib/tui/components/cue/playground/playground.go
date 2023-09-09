package playground

import (
	"time"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"github.com/gdamore/tcell/v2"

	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/components/cue/browser"
	"github.com/hofstadter-io/hof/lib/tui/components/cue/helpers"
	"github.com/hofstadter-io/hof/lib/tui/tview"
	"github.com/hofstadter-io/hof/lib/watch"
)


type valPack struct {
	config  helpers.SourceConfig
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


func New(initialText string, sourceConfig helpers.SourceConfig) (*Playground) {

	C := &Playground{
		Flex: tview.NewFlex(),
		text: initialText,
		scope: &valPack{
			config: sourceConfig,
		},
		final: &valPack{},
	}
	// our wrapper around the CUE widgets
	C.Flex = tview.NewFlex().SetDirection(tview.FlexColumn)

	// TODO, options form

	// editor
	C.edit = tview.NewTextArea()
	C.edit.
		SetTitle("  expression(s)  ").
		SetBorder(true)

	C.edit.SetText(C.text, false)

	// results
	C.final.viewer = browser.New(helpers.SourceConfig{}, "cue")
	C.final.viewer.SetName("result")
	C.final.viewer.SetBorder(true)

	// usingScope?
	if sourceConfig.Source != helpers.EvalNone {
		C.final.viewer.SetUsingScope(true)
		C.useScope = true

		C.scope.viewer = browser.New(C.scope.config, "cue")
		C.scope.viewer.SetName("scope")
		C.scope.viewer.SetBorder(true)
		C.Flex.AddItem(C.scope.viewer, 0, 1, true)
	} else {
		// add empty scope box
		C.Flex.AddItem(nil, 0, 0, false)
	}

	// layout
	C.Flex.
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

func (C *Playground) SetScopeConfig(sc helpers.SourceConfig) {
	C.scope.config = sc
}

func (C *Playground) UseScope(visible bool) {
	C.useScope = visible
}

func (C *Playground) SetFlexDirection(dir int) {
	C.SetDirection(dir)
}

func (C *Playground) Rebuild(rebuildScope bool) error {
	var (
		v cue.Value
		err error
	)

	ctx := cuecontext.New()
	src := C.edit.GetText()

	// compile a value
	if !C.useScope {
		// just compile the text
		v = ctx.CompileString(src, cue.InferBuiltins(true))
	} else {
		// compile the text with a scope

		// tui.Log("warn", fmt.Sprintf("%#v", s))
		sv, serr := C.scope.config.GetValue()
		err = serr

		if err == nil {
			if rebuildScope {
				// C.scope.config.Rebuild()
				cfg := helpers.SourceConfig{Value: sv}
				C.scope.viewer.SetSourceConfig(cfg)
				C.scope.viewer.Rebuild()
			}

			ctx := sv.Context()
			v = ctx.CompileString(src, cue.InferBuiltins(true), cue.Scope(sv))
		}
	}

	cfg := helpers.SourceConfig{Value: v}
	if err != nil {
		tui.Log("error", err)
		cfg = helpers.SourceConfig{Text: err.Error()}
	}
	// only update view value, that way, if we erase everything, we still see the value
	C.final.viewer.SetSourceConfig(cfg)
	C.final.viewer.Rebuild()

	// show/hide scope as needed
	if C.useScope {
		C.SetItem(0, C.scope.viewer, 0, 1, true)
	} else {
		C.SetItem(0, nil, 0, 0, false)
	}


	// tui.Draw()
	return nil
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


func (C *Playground) setupKeybinds() {
	// events (hotkeys)
	C.SetInputCapture(func(evt *tcell.EventKey) *tcell.EventKey {
		switch evt.Key() {
		case tcell.KeyRune:
			if (evt.Modifiers() & tcell.ModAlt) == tcell.ModAlt {
				switch evt.Rune() {
				case 'f':
					flexDir := C.GetDirection()
					if flexDir == tview.FlexRow {
						C.SetDirection(tview.FlexColumn)
					} else {
						C.SetDirection(tview.FlexRow)
					}

				case 'S':
					C.useScope = !C.useScope
					C.Rebuild(false)

				default: 
					return evt
				}

				return nil
			}

			return evt

		default:
			return evt
		}
	})	
}
