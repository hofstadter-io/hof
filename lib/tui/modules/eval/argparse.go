package eval

import (
	"fmt"
	"strings"

	"github.com/hofstadter-io/hof/lib/tui"
)

// this function helps parse args and context into richer information
// so that we can handle different kinds of input coming from different places
// and then have more consistent input to the components consuming the inputs
func enrichContext(context map[string]any) (map[string]any) {
	// tui.Log("trace", fmt.Sprintf("enrichContext.BEG: %# v", context))

	// ensure non-nil context
	if context == nil {
		context = make(map[string]any)
	}

	// extract and record original args
	args := []string{}
	if _args, ok := context["args"]; ok {
		args = _args.([]string)
	}
	context["orig-args"] = args

	//// so we can special case create a default element for naked eval
	// hadEval := false
	if len(args) > 0 && args[0] == "eval" {
		// hadEval = true
		args = args[1:]
	}

	// TODO, look to see if there is a ' -- '
	// and if exists, always parse here until it is reached
	// this will allow a user to be more exact and get better error messages

	// this is meant to be called in the switch in the middle of the loop
	// to consume an extra token when an @target.path is found
	maybeActionTarget := func() {
		if len(args) > 1 && strings.HasPrefix(args[1], "@") {
			context["target-panel"] = args[1][1:]  // also remove @
			// consume an extra arg, technically this is the action
			// but the target will be consumed right after anyway (just after the swith his is used in
			args = args[1:]
		}
		if len(args) > 1 && strings.HasPrefix(args[1], "#") {
			context["target-index"] = args[1][1:]  // also remove @
			// consume an extra arg, technically this is the action
			// but the target will be consumed right after anyway (just after the swith his is used in
			args = args[1:]
		}

	}

	// find any modifiers
	for len(args) > 0 {
		tok := args[0]
		switch tok {

		//
		// top-level eval (dash) commands
		//
		case 
			"preview",
			"save",
			"list",
			"show",
			"load":
			context["action"] = tok 

		//
		// actions
		//

		// panel actions
	 	case
			"create",
			"insert",
			"move",
			"split",
			"delete",
			"reload",
			"set.panel.name",
			"set.item.name",
			"set.name",
			"set.size",
			"set.ratio",
			"update":
			context["action"] = tok
			maybeActionTarget()

		// navigation between items & panels
		case "nav.left", "nl":
			context["action"] = "nav.left"
		case "nav.up", "nu":
			context["action"] = "nav.up"
		case "nav.down", "nd":
			context["action"] = "nav.down"
		case "nav.right", "nr":
			context["action"] = "nav.right"

		// in-item commands (these are for scope, but could be more)
		case 
			"add",
			"set":
			context["action"] = tok
			maybeActionTarget()

		// playground & (some)viewer commands
		case 
			"push",
			"export",
			"write",
			"set.value",
			"set.scope",
			// "set.text",

			"refresh", "watch", "watchGlobs",
			"get.refresh", "get.scope.watch", "get.value.watch",
			"set.refresh", "set.scope.watch", "set.value.watch",
			"set.scope.watchGlobs", "set.value.watchGlobs":
			context["action"] = tok 
			maybeActionTarget()

		case "connect", "conn":
			context["action"] = "conn"
			maybeActionTarget()

		//
		// items / widgets / pane
		//

		// default when nothing
		case "help":
			context["item"] = "help"
			context["action"] = "create"  // probably the default?

		// dual-pane eval'r (default when doing eval things)
		case "play":
			context["item"] = "play"
			context["action"] = "create"  // probably the default?

		// value viewer
		case "view":
			context["item"] = "view"
			context["action"] = "create"  // probably the default?

		case "flow":
			context["item"] = "flow"
			context["action"] = "create"  // probably the default?

		// should this be handled lower too?
		// we might want a more general 
		// sources?
		case "sh", "bash":
			context["source"] = "bash"

		//
		// what's acted upon, location within
		//
		
		// cue item action targets
		case "S", "scope":
			context["target"] = "scope"
		case "V", "value":
			context["target"] = "value"

		// eventually, which dashboard
		// here, there, new, extend

		// panel relative
		case "h", "head", "f", "first":
			context["where"] = "head"

		case "l", "last", "t", "tail", "e", "end":
			context["where"] = "tail"

		case "b", "before", "p", "prev":
			context["where"] = "prev"

		case "a", "after", "n", "next":
			context["where"] = "next"

		case "index", "i", "pos":
			context["where"] = "index"
			maybeActionTarget()

		// for within, to be passed down to current panel
		case "K":
			context["where"] = "top"
		case "k":
			context["where"] = "above"
		case "j":
			context["where"] = "below"
		case "J":
			context["where"] = "bottom"

		// default is the browser
		// setup the play with the eval as scope

		// these are more options, should be handled lower down, with item local parser
		//case "C", "cue":
		//  context["mode"] = "cue"
		//case "yaml":
		//  context["mode"] = "yaml"
		//case "json":
		//  context["mode"] = "json"


		// UPDATE, these probably become ALT-... keys
		// connect (start and end modes)
		// set colors while in progress
		// or also support passing all at once
		// also support mouse clicking

		// user intention, end of args
		case "--":
			goto argsDone

		default:
			// break doens't break here, but we want to stop processing args here
			goto argsDone
		}
		args = args[1:]
	}
argsDone:

	// handle some special cases
	if len(args) > 0 {
		first := args[0]
		if strings.HasPrefix(first, "http") {
			context["source"] = "http"
		}
	}

	// make sure we update the context args
	context["args"] = args

	tui.Log("trace", fmt.Sprintf("enrichContext.END: %# v", context))
	return context
}
