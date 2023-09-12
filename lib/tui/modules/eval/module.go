package eval

import (
	"github.com/hofstadter-io/hof/lib/connector"
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
	// tui.Log("extra", fmt.Sprintf("Eval.CmdCallback: %# v", context))
	M.Refresh(context)
	return
}
