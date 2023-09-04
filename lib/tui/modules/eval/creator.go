package eval

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"github.com/spf13/pflag"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/cuetils"
	"github.com/hofstadter-io/hof/lib/runtime"
	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/components"
	"github.com/hofstadter-io/hof/lib/tui/tview"
	"github.com/hofstadter-io/hof/lib/yagu"
)

// used for debugging panel CRUD & KEYS
var panel_debug = false

func init() {
	if v := os.Getenv("HOF_TUI_PANEL_DEBUG"); v != "" {
		vb, _ := strconv.ParseBool(v)
		if vb {
			panel_debug = true
		}
	}
}


// this function is responsable for creating the components that fill slots in the panel
// these are the widgets that make up the application and should have their own operation
func defaultCreator (context map[string]any, parent *Panel) (tview.Primitive) {
	tui.Log("trace", fmt.Sprintf("Panel.defaultCreator: %v", context ))

	// short-circuit for develop mode
	if panel_debug {
		t := tview.NewTextView()
		t.SetDynamicColors(true)
		fmt.Fprint(t, evalHelpText)
		return NewItem(t, parent)
	}

	// get infromation from context
	args := []string{}
	if _args, ok := context["args"]; ok {
		args = _args.([]string)
	}

	// needed for when we come in from the command line first time
	args, context = processArgsAndContext(args, context)

	mode := ""
	if _mode, ok := context["mode"]; ok {
		mode = _mode.(string)
	}

	// with ? already known?
	with := ""
	if _with, ok := context["with"]; ok {
		with = _with.(string)
	}

	var t tview.Primitive

	// decide what to build
	switch with {
	case "help":
		txt := tview.NewTextView()
		txt.SetDynamicColors(true)
		fmt.Fprint(txt, evalHelpText)
		t = txt

	case "http":
		var from string
		if _from, ok := context["from"]; ok {
			from = _from.(string)
		}
		parent.loadHttpValue(mode, from)
		e := components.NewValueEvaluator(parent._content, parent._Value, cue.Value{})
		e.Mount(context)
		e.Rebuild(context)
		t = e

	case "bash":
		var from string
		if _from, ok := context["from"]; ok {
			from = _from.(string)
		}
		parent.loadBashValue(mode, from, args)
		e := components.NewValueEvaluator(parent._content, parent._Value, cue.Value{})
		e.Mount(context)
		e.Rebuild(context)
		t = e

	default:
		// intentionally ignore, it should show up in the component below
		_ = parent.loadRuntime(args)

		if mode == "play" {
			e := components.NewValueEvaluator("", cue.Value{}, parent._Runtime.Value)
			e.SetScope(true)
			e.Mount(context)
			e.Rebuild(context)
			t = e
		} else {
			tui.Log("trace", fmt.Sprintf("defaultCreator.Runtime: %v %v", parent == nil, parent._Runtime == nil ))
			b := components.NewValueBrowser(parent._Runtime.Value, "cue", func(string){})
			b.SetTitle(fmt.Sprintf("  %v  ", args)).SetBorder(true)
			b.Mount(context)
			b.Rebuild("")
			t = b
		}
	}

	// set fallback error message
	if t == nil {
		txt := tview.NewTextView()
		txt.SetDynamicColors(true)
		fmt.Fprintf(txt, "unable to create element from context:\n%# v\n\n", context)
		fmt.Fprint(txt, evalHelpText)
		t = txt
	}

	return NewItem(t, parent)
}

func (P *Panel) loadRuntime(args []string) error {
	// tui.Log("trace", fmt.Sprintf("Panel.loadRuntime.inputs: %v", args))

	// build eval args & flags from the input args
	var (
		rflags flags.RootPflagpole
		cflags flags.EvalFlagpole
	)
	fset := pflag.NewFlagSet("panel", pflag.ContinueOnError)
	flags.SetupRootPflags(fset, &rflags)
	flags.SetupEvalFlags(fset, &cflags)
	fset.Parse(args)
	args = fset.Args()

	// tui.Log("trace", fmt.Sprintf("Panel.loadRuntime.parsed: %v %v", args, rflags))

	R, err := runtime.New(args, rflags)
	P._Runtime = R
	if err != nil {
		tui.Log("error", cuetils.ExpandCueError(err))
		return err
	}

	err = R.Load()
	P._Value = P._Runtime.Value
	if err != nil {
		tui.Log("error", cuetils.ExpandCueError(err))
		return err
	}

	return nil
}

func (P *Panel) loadHttpValue(mode, from string) error {
	// tui.Log("trace", fmt.Sprintf("Panel.loadHttpValue: %s %s", mode, from))

	// rework any cue/play links
	f := from
	if strings.Contains(from, "cuelang.org/play") {
		u, err := url.Parse(from)
		if err != nil {
			tui.Log("error", err)
			return err
		}
		id := u.Query().Get("id")
		f = fmt.Sprintf("https://%s/.netlify/functions/snippets?id=%s", u.Host, id)
	}

	// fetch content
	content, err := yagu.SimpleGet(f)
	if err != nil {
		tui.Log("error", err)
		P._content = err.Error()
		return err
	}

	// rebuild, TODO, if scope, use that value and scope.Context() here
	ctx := cuecontext.New()
	P._Value = ctx.CompileString(content, cue.InferBuiltins(true))

	content = "// from: " + from + "\n\n" + content
	P._content = content

	return nil
}

func (P *Panel) loadBashValue(mode, from string, args []string) error {

	wd, err := os.Getwd()
	if err != nil {
		tui.Log("error", err)
		return err
	}

	script := strings.Join(args, " ")
	out, err := yagu.Bash(script, wd)
	if err != nil {
		tui.Log("error", err)
		return err
	}

	// rebuild, TODO, if scope, use that value and scope.Context() here
	ctx := cuecontext.New()
	P._Value = ctx.CompileString(out, cue.InferBuiltins(true))

	script = "// from: " + from + " " + script
	P._content = script + "\n\n" + out


	return nil
}
