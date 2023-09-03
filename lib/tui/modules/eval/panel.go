package panel

import (
	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/hoc/router"
	"github.com/hofstadter-io/hof/lib/tui/tview"
)

type Panel struct {
	*tview.Flex
}

func NewPanel() *Panel {
	m := &Panel{
		Flex: tview.NewFlex(),
	}

	// do layout setup here

	return m
}

func (M *Panel) Id() string {
	return "panel"
}

func (M *Panel) Routes() []router.RoutePair {
	return []router.RoutePair{
		router.RoutePair{"/panel", M},
	}
}

func (M *Panel) Name() string {
	return "Panel"
}

func (M *Panel) HotKey() string {
	return ""
}

func (M *Panel) CommandName() string {
	return "panel"
}

func (M *Panel) CommandUsage() string {
	return "panel"
}

func (M *Panel) CommandHelp() string {
	return "help for panel module"
}

// CommandCallback is invoked when the user runs your module
// return the object you want in mount or refresh
func (M *Panel) CommandCallback(args []string, context map[string]interface{}) {
	if context == nil {
		context = make(map[string]any)
	}
	context["args"] = args

	if M.IsMounted() {
		// just refresh with new args
		M.Refresh(context)
	} else {
		// need to navigate, mount will do the rest
		context["path"] = "/panel"
		go tui.SendCustomEvent("/router/dispatch", context)
	}
}

func (M *Panel) Mount(context map[string]any) error {
	// this is where we can do some loading
	M.Flex.Mount(context)

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

func (M *Panel) Unmount() error {
	// this is where we can do some unloading, depending on the application
	//M.View.Unmount()
	//M.Eval.Unmount()
	M.Flex.Unmount()

	return nil
}

func (M *Panel) Refresh(context map[string]any) error {

	// this is where you update data and set in components
	// then at the end call tui.Draw()

	return nil
}
