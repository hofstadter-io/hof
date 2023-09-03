package eval

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"

	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/hoc/router"
	"github.com/hofstadter-io/hof/lib/tui/tview"
)

type Eval struct {
	*tview.Flex

	creator func (context map[string]any) *Panel
}

func NewEval() *Eval {
	m := &Eval{
		Flex: tview.NewFlex(),
		creator: defaultCreator,
	}

	// do layout setup here
	m.Flex.SetBorder(true).SetTitle("  Eval  ")
	m.Flex.SetDirection(tview.FlexColumn)

	return m
}

func (M *Eval) Id() string {
	return "eval"
}

func (M *Eval) Routes() []router.RoutePair {
	return []router.RoutePair{
		router.RoutePair{"/eval", M},
	}
}

func (M *Eval) Name() string {
	return "Eval"
}

func (M *Eval) HotKey() string {
	return ""
}

func (M *Eval) CommandName() string {
	return "eval"
}

func (M *Eval) CommandUsage() string {
	return "eval"
}

func (M *Eval) CommandHelp() string {
	return "help for eval module"
}

// CommandCallback is invoked when the user runs your module
// your goal is to enrich the context with the args
// return the object you want in Refresh
func (M *Eval) CommandCallback(args []string, context map[string]any) {
	tui.Log("error", fmt.Sprintf("eval cmd: %v %v", args, context))
	// strip of own command
	args, context = processArgsAndContext(args, context)

	if M.IsMounted() {
		tui.Log("error", fmt.Sprintf("eval mounted->refresh: %v %v", args, context))
		// just refresh with new args
		// maybe we need to be more intelligent here, make a different function(s)
		M.Refresh(context)
	} else {
		tui.Log("error", fmt.Sprintf("eval unmounted->router: %v %v", args, context))
		// need to navigate, mount will do the rest
		context["path"] = "/eval"
		go tui.SendCustomEvent("/router/dispatch", context)
	}
}


func processArgsAndContext(args []string, context map[string]any) ([]string, map[string]any) {
	tui.Log("info", fmt.Sprintf("parse.args-n-ctx.BEG: %v %v", args, context))

	if len(args) > 0 && args[0] == "eval" {
		args = args[1:]
	}

	// setup context
	if context == nil {
		context = make(map[string]any)
	}

	// find any modifiers
	for len(args) > 0 {
		tok := args[0]
		switch tok {

		case "i", "insert", "add":
			context["action"] = "insert"
		case "s", "scoped":
			context["action"] = "scoped"

		case "h", "head", "f", "first":
			context["where"] = "head"
		case "l", "last", "t", "tail", "e", "end":
			context["where"] = "last"

		case "b", "before", "p", "prev":
			context["where"] = "before"
		case "a", "after", "n", "next":
			context["where"] = "next"

		// for within, to be passed down to current panel
		case "K":
			context["where"] = "top"
		case "k":
			context["where"] = "above"
		case "j":
			context["where"] = "below"
		case "J":
			context["where"] = "bottom"

		case "P", "play":
			context["mode"] = "play"
		case "N", "new":
			context["mode"] = "new"
		case "T", "tree":
			context["mode"] = "tree"
		case "C", "cue":
			context["mode"] = "cue"
		case "yaml":
			context["mode"] = "yaml"
		case "json":
			context["mode"] = "json"
		case "flow":
			context["mode"] = "flow"
		case "eval":
			context["mode"] = "eval"

		// connect (start and end modes)
		// set colors while in progress
		// or also support passing all at once
		// also support mouse clicking
		case "x", "conn", "connect":
			context["action"] = "connect"
		case "X", "sconn", "scoped-connect":
			context["action"] = "scoped-connect"

		case "sh", "bash":
			context["with"] = tok

		case "--":
			goto argsDone
		default:
			if i, err := strconv.Atoi(tok); err == nil {
				context["where"] = "index"
				context["index"] = i
			} else {
				// break doens't break here, but we want to stop processing args here
				goto argsDone
			}
		}
		args = args[1:]
	}
argsDone:

	// handle some special cases
	if len(args) > 0 && strings.HasPrefix(args[0], "http") {
		context["with"] = "http"
		context["from"] = args[0]
		args = args[1:]
	}

	// make sure we update the context args
	context["args"] = args

	tui.Log("info", fmt.Sprintf("parse.args-n-ctx.END: %v %v", args, context))
	return args, context
}

func (M *Eval) Mount(context map[string]any) error {
	tui.Log("trace", fmt.Sprintf("Eval.Mount: %v", context))
	// this is where we can do some loading
	M.Flex.Mount(context)
	M.setupKeybinds()

	err := M.Refresh(context)
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

func (M *Eval) Unmount() error {
	// this is where we can do some unloading, depending on the application
	//M.View.Unmount()
	//M.Eval.Unmount()
	M.Flex.Unmount()

	// remove keybinds
	M.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey { return event })

	return nil
}

func (M *Eval) Refresh(context map[string]any) error {
	tui.Log("trace", fmt.Sprintf("Eval.Refresh: %v", context))
	// strip off the command name
	_args, _ := context["args"]
	args, _ := _args.([]string)
	if len(args) > 0 && args[0] == "eval" {
		args = args[1:]
		context["args"] = args
	}

	_action, _ := context["action"]
	action, _ := _action.(string)
	_index, _ := context["index"]
	index, _ := _index.(int)
	_where, _ := context["where"]
	where, _ := _where.(string)

	if action == "delete" {
	  M.Flex.RemoveIndex(index)
	} else if action != "" {

		switch where {
		case "head":
			M.Flex.InsItem(0, M.creator(context), 0, 1, true)

		case "before":
			M.Flex.InsItem(index, M.creator(context), 0, 1, true)

		case "next":
			M.Flex.InsItem(index+1, M.creator(context), 0, 1, true)

		case "last":
			M.Flex.AddItem(M.creator(context), 0, 1, true)

		case "index":
			M.Flex.InsItem(index, M.creator(context), 0, 1, true)

		default:
			M.Flex.AddItem(M.creator(context), 0, 1, true)

		}
	} else {
		// no action, probably coming to eval for the first time
		M.Flex.AddItem(M.creator(context), 0, 1, true)
	}

	// only set border when no elements
	M.Flex.SetBorder(M.Flex.GetItemCount() == 0)

	tui.Draw()

	return nil
}

func (M *Eval) setupKeybinds() {
	M.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		alt := event.Modifiers() & tcell.ModAlt == tcell.ModAlt
		//meta := event.Modifiers() & tcell.ModMeta == tcell.ModMeta

		ctx := make(map[string]any)

		// we only care about ALT+... keys at this level
		if alt {
			switch event.Key() {

			case tcell.KeyRune:
				switch event.Rune() {

				// insert head
				case 'h':
					ctx["action"] = "insert"
					ctx["where"] = "head"
				// scoped head
				case 'H':
					ctx["action"] = "scoped"
					ctx["where"] = "head"
				// insert before
				case 'b':
					ctx["action"] = "insert"
					ctx["where"] = "before"
				// scoped before
				case 'B':
					ctx["action"] = "scoped"
					ctx["where"] = "before"
				// insert next
				case 'n':
					ctx["action"] = "insert"
					ctx["where"] = "next"
				// scoped next
				case 'N':
					ctx["action"] = "scoped"
					ctx["where"] = "next"
				// insert last
				case 'l':
					ctx["action"] = "insert"
					ctx["where"] = "last"
				// scoped last
				case 'L':
					ctx["action"] = "scoped"
					ctx["where"] = "last"

				// DELETE
				case 'D':
					ctx["action"] = "delete"

				// flip flex orientation
				case 'F':
					d := M.Flex.GetDirection()
					if d == tview.FlexRow {
						d = tview.FlexColumn
					} else {
						d = tview.FlexRow
					}
					M.Flex.SetDirection(d)
					tui.Draw()
					return nil

				// dev stuff
				case 'v':
					focus := M.Flex.ChildFocus()
					count := M.Flex.GetItemCount()
					tui.Log("trace", fmt.Sprintf("Eval flex focus at %d of %d", focus, count))
					return nil

				default:
					return event

				}

				ctx["index"] = M.Flex.ChildFocus()

				M.Refresh(ctx)
				return nil

			}
		}
		return event
	})
}

func defaultCreator (context map[string]any) *Panel {
	tui.Log("trace", fmt.Sprintf("Eval.defaultCreator: %v", context))
	p := NewPanel()
	p.Mount(context)

	return p
}
