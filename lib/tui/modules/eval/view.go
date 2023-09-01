package eval

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	// "github.com/hofstadter-io/hof/lib/connector"
	"github.com/hofstadter-io/hof/lib/cuetils"
	"github.com/hofstadter-io/hof/lib/runtime"
	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/components"
	"github.com/hofstadter-io/hof/lib/tui/hoc/router"
	"github.com/hofstadter-io/hof/lib/tui/tview"
)

// Both a Module and a Layout and a Switcher.SubLayout
type EvalPage struct {
	*tview.Flex

	Runtime *runtime.Runtime

	Text *tview.TextView
	View *components.ValueBrowser
	Eval *components.ValueEvaluator

}

func NewEvalPage() *EvalPage {
	page := &EvalPage{
		Flex:	tview.NewFlex(),
		Text:	tview.NewTextView(),
	}

	// temp filler
	//page.View = tview.NewTextView()
	//page.Eval = tview.NewTextView()

	// main layout
	page.AddItem(page.Text, 0, 1, false)
		//AddItem(page.View, 0, 1, false).
		//AddItem(page.Eval, 0, 1, false)

	page.SetBorder(true).SetTitle("  eval  ")

	return page
}

func (P *EvalPage) Id() string {
	return "eval"
}

func (P *EvalPage) Name() string {
	return "Eval"
}

func (P *EvalPage) Mount(context map[string]any) error {
	// this is where we can do some loading
	P.Flex.Mount(context)

	err := P.Refresh(context)
	if err != nil {
		tui.SendCustomEvent("/console/error", err)
		return err
	}

	P.View.Mount(context)
	P.Eval.Mount(context)


	return nil
}

func (P *EvalPage) Unmount() error {
	// this is where we can do some loading
	//P.View.Unmount()
	//P.Eval.Unmount()
	P.Flex.Unmount()

	return nil
}

func (P *EvalPage) Refresh(context map[string]interface{}) error {
	tui.SendCustomEvent("/console/trace", fmt.Sprintf("eval refresh %#v", context))

	_args, _ := context["args"]
	args, _ := _args.([]string)
	if len(args) > 0 && args[0] == "eval" {
		args = args[1:]
	}

	var (
		rflags flags.RootPflagpole
		cflags flags.EvalFlagpole
	)
	fset := pflag.NewFlagSet("root", pflag.ContinueOnError)
	flags.SetupRootPflags(fset, &rflags)
	flags.SetupEvalFlags(fset, &cflags)

	err := fset.Parse(args)
	if err != nil {
		P.ShowTextError(err)
		return err
	}

	args = fset.Args()

	d := map[string]any{
		"cmd": "eval",
		"args": args,
		"rflags": rflags,
		"cflags": cflags,
	}
	tui.SendCustomEvent("/console/trace", fmt.Sprintf("eval refresh %#v", d))


	R, err := runtime.New(args, rflags)
	// write error to text area error?
	if err != nil {
		P.ShowTextError(err)
		return err
	}

	err = R.Load()
	if err != nil {
		P.ShowTextError(err)
		return err
	}

	P.Runtime = R
	
	// setup file browser
	onNodeSelect := func(path string) {
		// app.Logger("onNodeSelect: " + path)
	}

	if P.View == nil {
		P.View = components.NewValueBrowser(P.Runtime.Value, onNodeSelect)
		P.Eval = components.NewValueEvaluator(R)
	} else {
		P.View.Value = R.Value
		P.Eval.Runtime = R
	}
	P.View.SetTitle(strings.Join(P.Runtime.Entrypoints, " "))
	P.View.Rebuild("")
	P.Eval.Rebuild()

	// remove Flex Item
	// add Items

	P.Flex.Clear()
	P.SetBorder(false)

	// P.Flex.RemoveItem(P.Text)
	P.Flex.
		AddItem(P.View, 0, 1, false).
		AddItem(P.Eval, 0, 1, true)


	tui.Draw()


	return nil
}

func (P *EvalPage) ShowTextError(err error) {
	s := fmt.Sprintf("%v\n", cuetils.ExpandCueError(err))
	P.Text.Clear()
	fmt.Fprint(P.Text, s)

	P.Flex.Clear()
	P.SetBorder(false)
	P.Text.SetBorder(true).SetTitle("  ERROR  ")

	// P.Flex.RemoveItem(P.Text)
	P.Flex.AddItem(P.Text, 0, 1, false)
	if P.Eval != nil {
		P.Flex.AddItem(P.Eval, 0, 1, true)
	}

	tui.Draw()
}

func (P *EvalPage) Routes() []router.RoutePair {
	return []router.RoutePair{
		router.RoutePair{"/eval", P},
	}
}

func (P *EvalPage) CommandName() string {
	return "eval"
}

func (P *EvalPage) CommandUsage() string {
	return "eval <args> [flags]"
}

func (P *EvalPage) CommandHelp() string {
	return "explore cue values"
}

// this takes the args from the command bar and turns it into a router dispatch
func (P *EvalPage) CommandCallback(args []string, context map[string]any) {
	if context == nil {
		context = make(map[string]any)
	}
	context["args"] = args

	if P.IsMounted() {
		// just refresh with new args
		P.Refresh(context)
	} else {
		// need to navigate, mount will do the rest
		context["path"] = "/eval"
		tui.SendCustomEvent("/console/trace", fmt.Sprintf("eval navigate %v", context))
		go tui.SendCustomEvent("/router/dispatch", context)
	}
}
