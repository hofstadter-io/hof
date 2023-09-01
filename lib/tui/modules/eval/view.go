package eval

import (
	"fmt"
	"os"
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

func (P *EvalPage) Mount(context map[string]interface{}) error {
	// this is where we can do some loading
	tui.SendCustomEvent("/console/warn", "eval mount")

	P.Refresh(context)

	P.View.Mount(context)
	P.Eval.Mount(context)


	return nil
}

func (P *EvalPage) Refresh(context map[string]interface{}) error {
	tui.SendCustomEvent("/console/warn", "eval refresh")
	args := []string{"eval"}

	if context != nil {
		if a, ok := context["args"]; ok {
			args = append(args, a.([]string)...)
		}
	}


	os.Args = args
	pflag.Parse()
	rflags := flags.RootPflags
	// cflags := flags.EvalFlags

	R, err := runtime.New(pflag.Args(), rflags)
	// write error to text area error?
	if err != nil {
		s := fmt.Sprintf("%v\n", cuetils.ExpandCueError(err))
		fmt.Fprint(P.Text, s)
		tui.SendCustomEvent("/console/error", s)
		return nil
	}

	err = R.Load()
	if err != nil {
		s := fmt.Sprintf("%v\n", cuetils.ExpandCueError(err))
		fmt.Fprint(P.Text, s)
		tui.SendCustomEvent("/console/error", s)
		return nil
	}

	P.Runtime = R
	
	// setup file browser
	onNodeSelect := func(path string) {
		// app.Logger("onNodeSelect: " + path)
	}

	P.View = components.NewValueBrowser(P.Runtime.Value, onNodeSelect)
	P.View.Rebuild("")

	P.View.SetTitle(strings.Join(P.Runtime.Entrypoints, " "))

	P.Eval = components.NewValueEvaluator(R)

	// remove Flex Item
	// add Items

	P.Flex.Clear()
	P.SetBorder(false)

	// P.Flex.RemoveItem(P.Text)
	P.Flex.
		AddItem(P.View, 0, 1, false).
		AddItem(P.Eval, 0, 1, true)


	return nil
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

	if P.IsMounted() {
		// should already be on page and visible (?)
		context["args"] = args

		P.Refresh(context)
	} else {
		// just navigate, 
		context["path"] = "/eval"
		context["args"] = args

		go tui.SendCustomEvent("/router/dispatch", context)
	}
}
