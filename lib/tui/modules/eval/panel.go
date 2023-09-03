package eval

import (
	"fmt"
	"net/url"
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

type Panel struct {
	*tview.Flex

	creator func (context map[string]any) *Panel

	_Runtime *runtime.Runtime
	_Value   cue.Value
	_content string
}

func NewPanel() *Panel {
	P := &Panel{
		Flex: tview.NewFlex(),
		creator: defaultCreator,
	}

	return P
}

func defaultCreator (context map[string]any) (*Panel) {
	P := NewPanel()
	if debug {
		P.Flex.SetBorder(true).SetTitle(fmt.Sprintf("p:%02d", count))
		setupInputHandler(P)
		count++
	} else {
		P.Mount(context)
	}

	return P
}


func (P *Panel) Id() string {
	return "panel"
}

func (P *Panel) Mount(context map[string]any) error {
	tui.Log("trace", fmt.Sprintf("Panel.Mount: %v", context))
	// this is where we can do some loading
	P.Flex.Mount(context)

	err := P.Refresh(context)
	if err != nil {
		tui.SendCustomEvent("/console/error", err)
		return err
	}

	// mount any other components
	// maybe we should have [...Children], so two pointers, one for dev, one for sys (Children)
	// then this call to mount can be handled without extra stuff by default?
	//M.View.Mount(context)
	//M.Eval.Mount(context)

	return nil
}

func (P *Panel) Unmount() error {
	// this is where we can do some unloading, depending on the application
	//M.View.Unmount()
	//M.Eval.Unmount()
	P.Flex.Unmount()

	return nil
}

func (P *Panel) Refresh(context map[string]any) error {
	tui.Log("trace", fmt.Sprintf("Panel.Refresh: %v", context))

	// get and setup args
	args := []string{}
	if _args, ok := context["args"]; ok {
		args = _args.([]string)
	}
	args, context = processArgsAndContext(args, context)

	// extract some info from context
	mode := ""
	if _mode, ok := context["mode"]; ok {
		mode = _mode.(string)
	}

	// do things based on context info to build up a component
	var t tview.Primitive

	// with ? already known?
	with := ""
	if _with, ok := context["with"]; ok {
		with = _with.(string)
	}
	switch with {
	case "http":
		var from string
		if _from, ok := context["from"]; ok {
			from = _from.(string)
		}
		P.loadHttpValue(mode, from)
		e := components.NewValueEvaluator(P._content, P._Value, cue.Value{})
		e.Mount(context)
		e.Rebuild(context)
		t = e


	default:
		P.loadRuntime(args)

		if mode == "play" {
			e := components.NewValueEvaluator("", cue.Value{}, P._Runtime.Value)
			e.SetScope(true)
			e.Mount(context)
			e.Rebuild(context)
			t = e
		} else {
			b := components.NewValueBrowser(P._Runtime.Value, "cue", func(string){})
			b.SetTitle(fmt.Sprintf("  %v  ", args)).SetBorder(true)
			b.Mount(context)
			b.Rebuild("")
			t = b
		}
	}


	// setup sub-componets
	P.Flex.AddItem(t, 0, 1, false)

	// this is where you update data and set in components
	// then at the end call tui.Draw()

	return nil
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
	if err != nil {
		tui.Log("error", cuetils.ExpandCueError(err))
		return err
	}

	err = R.Load()
	if err != nil {
		tui.Log("error", cuetils.ExpandCueError(err))
		return err
	}

	// tui.Log("warn", fmt.Sprintf("%#v", R.Value))

	P._Runtime = R
	P._Value = P._Runtime.Value

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

func NewItemText() *tview.TextView {
	t := tview.NewTextView()
	t.SetBorder(true)
	fmt.Fprint(t, newItemText)
	return t
}
 
const newItemText = `
Panel Controls:
<esc>        unfocus
alt-M        set mode
alt-[HhlL]   horz panel inserts
alt-[JjkK]   vert panel inserts
ctrl-[HhlL]  horz panel move
ctrl-[JjkK]  vert panel move
`
