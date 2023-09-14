package browser

import (
	"fmt"

	"github.com/hofstadter-io/hof/lib/tui"
)

// local handler type
type HandlerFunc func (B *Browser, action string, args []string, context map[string]any) (handled bool, err error)

// action registry
var actions = map[string]HandlerFunc{
	"create": handleSet,
	"set":    handleSet,
	"add":    handleAdd,
	"conn":   handleAdd,
	"watch":  handleWatchConfig,
	"globs":  handleWatchConfig,
}

// implementation of widget.ActionHandler interface
func (B *Browser) HandleAction(action string, args []string, context map[string]any) (bool, error) {
	tui.Log("warn", fmt.Sprintf("Browser.HandleAction: %v %v %v", action, args, context))

	handler, ok := actions[action]
	if ok {
		return handler(B, action, args, context)
	}

	return false, nil
}
