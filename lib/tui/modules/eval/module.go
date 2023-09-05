package eval

import (
	"fmt"

	"github.com/hofstadter-io/hof/lib/connector"
	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/hoc/router"
)

func New() connector.Connector {
	items := []any{
		NewEval(),
	}
	m := connector.New("Eval")
	m.Add(items)

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
func (M *Eval) CommandCallback(context map[string]any) {
	// strip of own command
	tui.Log("error", fmt.Sprintf("Eval.CmdCallback: %# v", context))

	if M.IsMounted() {
		// tui.Log("error", fmt.Sprintf("eval mounted->refresh: %v %v", args, context))
		// just refresh with new args
		// maybe we need to be more intelligent here, make a different function(s)
		M.Refresh(context)
	} else {
		// tui.Log("error", fmt.Sprintf("eval unmounted->router: %v %v", args, context))
		// need to navigate, mount will do the rest
		context["path"] = "/eval"
		go tui.SendCustomEvent("/router/dispatch", context)
	}
}
