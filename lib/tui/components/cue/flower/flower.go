package flower

import (
	"bytes"

	"cuelang.org/go/cue"

	"github.com/hofstadter-io/hof/cmd/hof/flags"

	flowcontext "github.com/hofstadter-io/hof/flow/context"
	"github.com/hofstadter-io/hof/flow/middleware"
	"github.com/hofstadter-io/hof/flow/tasks"
	"github.com/hofstadter-io/hof/flow/flow"

	"github.com/hofstadter-io/hof/lib/tui/components/cue/browser"
	"github.com/hofstadter-io/hof/lib/tui/components/cue/helpers"
	"github.com/hofstadter-io/hof/lib/tui/tview"
)

type valPack struct {
	config  *helpers.SourceConfig
	value   cue.Value
	viewer  *browser.Browser // scope
}


type Flower struct {
	*tview.Flex

	scope    *valPack

	edit *tview.TextArea  // text

	// value from which flow is created
	value cue.Value

	// TODO, flow DAG viewer

	// TODO, stdout/stderr viewer

	// FOR NOW, something to explore the dag during/after running
	final    *valPack

	// something to run it
	flowctx *flowcontext.Context
	stdin, stdout, stderr bytes.Buffer


	// def want to try terraform/dag package

}

func (F *Flower) TypeName() string {
	return "cue/flower"
}

func New() (*Flower) {

	F := &Flower{
		Flex: tview.NewFlex(),
		scope: &valPack{},
		final: &valPack{},
	}

	// our wrapper around the CUE widgets
	F.Flex = tview.NewFlex().SetDirection(tview.FlexColumn)

	// scope viewer
	F.scope.config = &helpers.SourceConfig{}
	F.scope.viewer = browser.New(F.scope.config, "cue")
	F.scope.viewer.SetName("scope")
	F.scope.viewer.SetBorder(true)

	// curr editor
	F.edit = tview.NewTextArea()
	F.edit.
		SetTitle("  expression(s)  ").
		SetBorder(true)

	F.edit.SetText("", false)

	// TODO, options form

	// TODO, flow DAG viewer

	// for now, results viewer
	F.final.config = &helpers.SourceConfig{}
	F.final.viewer = browser.New(F.final.config, "cue")
	F.final.viewer.SetName("result")
	F.final.viewer.SetBorder(true)
	F.final.viewer.SetUsingScope(false)

	// layout
	F.Flex.
		AddItem(F.scope.viewer, 0, 1, true).
		AddItem(F.edit, 0, 1, true).
		AddItem(F.final.viewer, 0, 1, true)

	return F
}


func (F *Flower) runFlow() error {

	postFlow := F.value

		


	ctx := flowcontext.New()
	ctx.RootValue = postFlow
	ctx.Stdin = &F.stdin
	ctx.Stdout = &F.stdout
	ctx.Stderr = &F.stderr

	F.flowctx = ctx

	// how to inject tags into original value
	// fill / return value
	ropts := flags.RootPflagpole{}
	fopts := flags.FlowPflagpole{}
	middleware.UseDefaults(ctx, ropts, fopts)
	tasks.RegisterDefaults(ctx)

	p, err := flow.OldFlow(ctx, postFlow)
	if err != nil {
		return err
	}

	err = p.Start()
	if err != nil {
		return err
	}

	return nil
}
