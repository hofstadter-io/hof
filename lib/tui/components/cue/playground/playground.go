package playground

import (
	"fmt"
	"time"

	"github.com/hofstadter-io/hof/lib/tui/components/cue/browser"
	"github.com/hofstadter-io/hof/lib/tui/components/cue/helpers"
	"github.com/hofstadter-io/hof/lib/tui/tview"
	"github.com/hofstadter-io/hof/lib/watch"
)

type PlayMode string

const (
	ModeEval PlayMode = "eval"
	ModeFlow PlayMode = "flow"
)

type Playground struct {
	// *tview.Frame eventually?
	*tview.Flex

	// scope used during parsing / evaluation
	seeScope bool
	useScope bool
	scope    *browser.Browser

	// the editor box
	edit *tview.TextArea  // text
	editCfg *helpers.SourceConfig

	mode PlayMode

	// the final value
	final *browser.Browser // scope

	// for handling TUI inputs
	debouncer func(func()) // that's funky!
	debounceTime time.Duration
}

func (*Playground) TypeName() string {
	return "cue/playground"
}

func New(initialText string) (*Playground) {

	C := &Playground{
		Flex: tview.NewFlex(),
		debounceTime: time.Millisecond * 500,
		mode: ModeEval,
	}

	// our wrapper around the CUE widgets
	C.Flex = tview.NewFlex().SetDirection(tview.FlexColumn)

	// TODO, options form

	// scope viewer
	C.scope = browser.New()
	C.scope.SetName("scope")
	C.scope.SetBorder(true)

	// curr editor
	C.editCfg = &helpers.SourceConfig{}
	C.edit = tview.NewTextArea()
	C.edit.
		SetTitle("  expression(s)  ").
		SetBorder(true)

	C.edit.SetText(initialText, false)

	// results viewer
	C.final = browser.New()
	C.final.SetName("result")
	C.final.SetBorder(true)

	// usingScope?
	C.final.SetUsingScope(false)
	C.useScope = false

	// layout
	C.Flex.
		AddItem(C.scope, 0, 1, true).
		AddItem(C.edit, 0, 1, true).
		AddItem(C.final, 0, 1, true)

	C.setupKeybinds()

	// setup change response with douncer
	// to trigger rebuild on editor changes
	C.debouncer = watch.NewDebouncer(C.debounceTime)
	C.edit.SetChangedFunc(func() {
		C.debouncer(func(){
			C.Rebuild()
		})
	})
	return C
}

func (C *Playground) UseScope(use bool) {
	C.useScope = use
	C.final.SetUsingScope(use)
}

func (C *Playground) ToggleShowScope() {
	C.seeScope = !C.seeScope
}

func (C *Playground) RebuildEditTitle() {
	// when not showing scope and has scope
	// display in editory text
	s := ""
	if !C.seeScope {
		if l := len(C.scope.GetSourceConfigs()); l > 0 {
			s += fmt.Sprintf("[violet](srcs:%d)[-] ", l)
		}
	}
	C.edit.SetTitle(fmt.Sprintf("  [gold]%s[-] %sexpr(s)  ", C.mode, s))
}

func (C *Playground) SetFlexDirection(dir int) {
	C.SetDirection(dir)
}


func (C *Playground) Mount(context map[string]any) error {

	return nil
}

func (C *Playground) Focus(delegate func(p tview.Primitive)) {
	if C.scope.HasFocus() {
		delegate(C.scope)
		return
	}
	if C.final.HasFocus() {
		delegate(C.final)
		return
	}
	// otherwise, assume we want to keep the view focus
	delegate(C.edit)
	return
}

func (C *Playground) GetSourceConfigs() []*helpers.SourceConfig {
	return C.scope.GetSourceConfigs()
}

func (C *Playground) GetEditConfig() *helpers.SourceConfig {
	return C.editCfg
}
