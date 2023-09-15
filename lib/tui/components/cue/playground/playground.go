package playground

import (
	"fmt"
	"strings"
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
	seeScope bool
	useScope bool
	scope    *valPack

	// the editor box
	edit *tview.TextArea  // text
	editCfg *helpers.SourceConfig

	// the final value
	final    *valPack

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
		scope: &valPack{},
		final: &valPack{},

		debounceTime: time.Millisecond * 500,
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
	C.editCfg = &helpers.SourceConfig{}
	C.edit = tview.NewTextArea()
	C.edit.
		SetTitle("  expression(s)  ").
		SetBorder(true)

	C.edit.SetText(initialText, false)

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
	C.debouncer = watch.NewDebouncer(C.debounceTime)
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

func (C *Playground) GetScopeConfig() *helpers.SourceConfig {
	return C.scope.config
}

func (C *Playground) GetEditConfig() *helpers.SourceConfig {
	return C.editCfg
}

func (C *Playground) UseScope(use bool) {
	C.useScope = use
	C.final.viewer.SetUsingScope(use)
}

func (C *Playground) ToggleShowScope() {
	C.seeScope = !C.seeScope
	// when not showing scope and has scope
	// display in editory text
	s := ""
	if !C.seeScope {
		if len(C.scope.config.Args) > 0 {
			s += "[violet](" + strings.Join(C.scope.config.Args, " ") + ")[-] "
		}
	}
	C.edit.SetTitle(fmt.Sprintf("  %sexpression(s)  ", s))
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


