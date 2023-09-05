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

// we need this to be able to handle both making new AND updating existing components

// this is going to take a bit more time...


// this function is responsable for creating the components that fill slots in the panel
// these are the widgets that make up the application and should have their own operation
func (P *Panel) creator(context map[string]any, parent *Panel) (*Item, error) {
	tui.Log("extra", fmt.Sprintf("Panel.creator: %v", context ))

	// short-circuit for developer mode (first, before user custom)
	if panel_debug {
		t := NewTextView()
		i := NewItem(context, parent)
		i.SetWidget(t)
		return i, nil
	}

	// run a user custome panel, if provided
	if P._creator != nil {
		return P._creator(context, parent)
	}

	args := []string{}
	if _args, ok := context["args"]; ok {
		// because in-mem vs decode-yaml...
		switch _args := _args.(type) {
		case []string:
			args = _args
		case []any:
			for _, a := range _args {
				args = append(args, a.(string))
			}
		}
	}

	item := ""
	if _item, ok := context["item"]; ok {
		item = _item.(string)
	}

	I := NewItem(context, parent)

	switch item {
	case "help":
		tui.Log("debug", "Panel.creator: help")
		txt := NewTextView()
		fmt.Fprint(txt, EvalHelpText)
		I.SetWidget(txt)

	case "play":
		tui.Log("debug", "Panel.creator: play")
		I._runtimeArgs = args
		_ = I.loadRuntime(args)
		e := components.NewValueEvaluator("", cue.Value{}, I._runtime.Value)
		e.SetScope(true)
		e.Mount(context)
		e.Rebuild(context)
		I.SetWidget(e)

	case "tree":
		tui.Log("debug", "Panel.creator: tree")
		I._runtimeArgs = args
		_ = I.loadRuntime(args)
		b := components.NewValueBrowser(I._runtime.Value, "cue", func(string){})
		b.SetTitle(fmt.Sprintf("  %v  ", args)).SetBorder(true)
		b.Mount(context)
		b.Rebuild("")
		I.SetWidget(b)

	default:
		return I.defaultWidget(context, parent)

	}

	return I, nil
}

func (I *Item) defaultWidget(context map[string]any, parent *Panel) (*Item, error) {
	tui.Log("debug", fmt.Sprintf("Panel.defaultWidget: %v", context ))

	// with ? already known?
	source := ""
	if _src, ok := context["source"]; ok {
		source = _src.(string)
	}

	// decide what to build
	switch source {

	case "http":
		var from string
		if _from, ok := context["from"]; ok {
			from = _from.(string)
		}

		I.loadHttpValue(from)
		e := components.NewValueEvaluator(I._text, I._value, cue.Value{})
		e.Mount(context)
		e.Rebuild(context)
		I.SetWidget(e)

	case "bash":
		args := []string{}
		if _args, ok := context["args"]; ok {
			args = _args.([]string)
		}

		I.loadBashValue(args)
		e := components.NewValueEvaluator(I._text, I._value, cue.Value{})
		e.Mount(context)
		e.Rebuild(context)
		I.SetWidget(e)

	default:
		txt := NewTextView()
		fmt.Fprint(txt, fmt.Sprintf("unhandled item create: \n%# v\n\n", context))
		fmt.Fprint(txt, EvalHelpText)
		I.SetWidget(txt)

	}

	return I, nil
}

func (I *Item) loadRuntime(args []string) error {
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
	I._runtime = R
	if err != nil {
		tui.Log("error", cuetils.ExpandCueError(err))
		return err
	}

	err = R.Load()
	I._value = I._runtime.Value
	if err != nil {
		tui.Log("error", cuetils.ExpandCueError(err))
		return err
	}

	return nil
}

func (I *Item) loadHttpValue(from string) error {
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
		I._text = err.Error()
		return err
	}

	// rebuild, TODO, if scope, use that value and scope.Context() here
	ctx := cuecontext.New()
	I._value = ctx.CompileString(content, cue.InferBuiltins(true))

	content = "// from: " + from + "\n\n" + content
	I._text = content

	return nil
}

func (I *Item) loadBashValue(args []string) error {

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
	I._value = ctx.CompileString(out, cue.InferBuiltins(true))

	script = "// bash: " + " " + script
	I._text = script + "\n\n" + out

	return nil
}
