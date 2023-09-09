package playground

import (
	"fmt"
	"strings"
	"time"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"github.com/gdamore/tcell/v2"
	"github.com/parnurzeal/gorequest"

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


func New(initialText string, scopeSourceConfig helpers.SourceConfig) (*Playground) {

	C := &Playground{
		Flex: tview.NewFlex(),
		text: initialText,
		scope: &valPack{
			config: scopeSourceConfig,
		},
		final: &valPack{},
	}
	// our wrapper around the CUE widgets
	C.Flex = tview.NewFlex().SetDirection(tview.FlexColumn)

	// TODO, options form

	// scope viewer
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
	C.final.viewer = browser.New(helpers.SourceConfig{}, "cue")
	C.final.viewer.SetName("result")
	C.final.viewer.SetBorder(true)

	// usingScope?
	C.final.viewer.SetUsingScope(true)
	C.useScope = true
	//if sourceConfig.Source != helpers.EvalNone {

	//} else {
	//  // add empty scope box
	//  C.Flex.AddItem(nil, 0, 0, false)
	//}

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

func (C *Playground) SetScopeConfig(sc helpers.SourceConfig) {
	C.scope.config = sc
}

func (C *Playground) UseScope(visible bool) {
	C.useScope = visible
}

func (C *Playground) SetFlexDirection(dir int) {
	C.SetDirection(dir)
}


const HTTP2_GOAWAY_CHECK = "http2: server sent GOAWAY and closed the connection"

func (C *Playground) PushToPlayground() (string, error) {
	src := C.edit.GetText()

	url := "https://cuelang.org/.netlify/functions/snippets"
	req := gorequest.New().Post(url)
	req.Set("Content-Type", "text/plain")
	req.Send(src)

	resp, body, errs := req.End()

	if len(errs) != 0 && !strings.Contains(errs[0].Error(), HTTP2_GOAWAY_CHECK) {
		fmt.Println("errs:", errs)
		fmt.Println("resp:", resp)
		fmt.Println("body:", body)
		return body, errs[0]
	}

	if len(errs) != 0 || resp.StatusCode >= 500 {
		return body, fmt.Errorf("Internal Error: " + body)
	}
	if resp.StatusCode >= 400 {
		return body, fmt.Errorf("Bad Request: " + body)
	}

	return body, nil
}

func (C *Playground) Rebuild(rebuildScope bool) error {
	// tui.Log("info", fmt.Sprintf("Play.rebuildScope %v %v %v", rebuildScope, C.useScope, C.scope.config))
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

		if err != nil {
			tui.Log("error", err)
		}
		// we shouldn't have to worry about this, but we aren't catching all the ways
		// that we get into this code, in particular, hotkey can set scope to true when none exists
		if !sv.Exists() {
			tui.Log("error", "scope value does not exist")
			err = fmt.Errorf("scope value does not exist")
		}

		if err == nil && sv.Exists() {
			if rebuildScope {
				// C.scope.config.Rebuild()
				cfg := helpers.SourceConfig{Value: sv}
				C.scope.viewer.SetSourceConfig(cfg)
				C.scope.viewer.Rebuild()
			}

			// tui.Log("warn", fmt.Sprintf("recompile with scope: %v", rebuildScope))
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
	C.final.viewer.SetUsingScope(C.useScope)
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

				case 'R':
					C.Rebuild(true)

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
