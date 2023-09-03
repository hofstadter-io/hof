package eval

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"

	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/tview"
)

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


func setupInputHandler(P *Panel) {
	P.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		alt := event.Modifiers() & tcell.ModAlt == tcell.ModAlt
		ctrl := event.Modifiers() & tcell.ModCtrl == tcell.ModCtrl

		ctx := make(map[string]any)
		flexDir := P.Flex.GetDirection()

		// we only care about ALT+... keys at this level
		tui.Log("trace", fmt.Sprintf("Panel.inputHandler %v %v %v", alt, ctrl, string(event.Rune())))

		switch event.Key() {

		case tcell.KeyESC:
			focus := P.Flex.ChildFocus()
			count := P.Flex.GetItemCount()
			tui.Log("trace", fmt.Sprintf("Eval flex focus at %d of %d", focus, count))

		case tcell.KeyRune:
			handled := false
			{ // scope for action & handled clarity
				action := ""
				if alt {
					action = "insert"
				}
				if ctrl {
					action = "move"
				}

				if flexDir == tview.FlexColumn {
					handled = true
					switch event.Rune() {
					// insert head
					case 'H':
						ctx["action"] = action
						ctx["where"] = "head"
					// insert before
					case 'h':
						ctx["action"] = action
						ctx["where"] = "before"
					// insert next
					case 'l':
						ctx["action"] = action
						ctx["where"] = "next"
					// insert last
					case 'L':
						ctx["action"] = action
						ctx["where"] = "last"
					// not interested
					default:
						handled = false
					}
				} else {
					handled = true
					switch event.Rune() {
					// insert head
					case 'K':
						ctx["action"] = action
						ctx["where"] = "top"
					// insert before
					case 'k':
						ctx["action"] = action
						ctx["where"] = "above"
					// insert next
					case 'j':
						ctx["action"] = action
						ctx["where"] = "below"
					// insert last
					case 'J':
						ctx["action"] = action
						ctx["where"] = "bottom"
					// not interested
					default:
						handled = false
					}
				}
			}

			if !handled {
				switch event.Rune() {

				case 't':
					// split the opposite direction
					ctx["action"] = "split"
					if flexDir == tview.FlexRow {
						ctx["where"] = "vert"
					} else {
						ctx["where"] = "horz"
					}

				// DELETE
				case 'D':
					ctx["action"] = "delete"

				// flip flex orientation
				case 'F':
					if flexDir == tview.FlexRow {
						flexDir = tview.FlexColumn
					} else {
						flexDir = tview.FlexRow
					}
					P.Flex.SetDirection(flexDir)
					tui.Draw()
					return nil

				// dev stuff
				case 'v':
					focus := P.Flex.ChildFocus()
					count := P.Flex.GetItemCount()
					tui.Log("trace", fmt.Sprintf("Eval flex focus at %d of %d", focus, count))
					return nil

				default:
					return event

				}
				handled = true
			}

			if handled {
				// ???
				ctx["index"] = P.Flex.ChildFocus()

				P.Refresh(ctx)
				tui.Draw()
				return nil
			}

		}

		return event
	})
}
