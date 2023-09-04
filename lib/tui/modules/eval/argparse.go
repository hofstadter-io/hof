package eval

import (
	"strconv"
	"strings"
)

// this function helps parse args and context into richer information
// so that we can handle different kinds of input coming from different places
// and then have more consistent input to the components consuming the inputs
func processArgsAndContext(args []string, context map[string]any) ([]string, map[string]any) {
	// tui.Log("info", fmt.Sprintf("parse.args-n-ctx.BEG: %v %v", args, context))

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
		//case "s", "scoped":
		//  context["action"] = "scoped"

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

		case "help":
			context["with"] = "help"

		// UPDATE, these probably become ALT-... keys
		// connect (start and end modes)
		// set colors while in progress
		// or also support passing all at once
		// also support mouse clicking
		//case "x", "conn", "connect":
		//  context["action"] = "connect"
		//case "X", "sconn", "scoped-connect":
		//  context["action"] = "scoped-connect"

		case "sh", "bash":
			context["with"] = tok

		// user intention, end of args
		case "--":
			goto argsDone

		default:
			// is it an int?
			if i, err := strconv.Atoi(tok); err == nil {
				context["where"] = "index"
				// maybe make this a list, to support x,y row,col or more indexes
				// or maybe make them CSV each, like pos ~ 1,2
				// but maybe easier as a slice of ints that the function can interpret as it wishes, ordered like
				// we'd have a hard time assigning them to different keys here
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

	// tui.Log("info", fmt.Sprintf("parse.args-n-ctx.END: %v %v", args, context))

	return args, context
}
