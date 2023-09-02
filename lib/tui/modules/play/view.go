package play

import (
	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/tview"
	"github.com/hofstadter-io/hof/lib/tui/hoc/router"
)

// Both a Play and a Layout and a Switcher.SubLayout
type Play struct {
	*tview.Flex
}

func NewPlay() *Play {
	p := &Play{
		Flex: tview.NewFlex(),
	}

	// do layout setup here

	return p
}

func (P *Play) Id() string {
	return "play"
}

func (P *Play) Routes() []router.RoutePair {
	return []router.RoutePair{
		router.RoutePair{"/play", P},
	}
}

func (P *Play) Name() string {
	return "Play"
}

func (P *Play) HotKey() string {
	return ""
}

func (P *Play) CommandName() string {
	return "play"
}

func (P *Play) CommandUsage() string {
	return "play"
}

func (P *Play) CommandHelp() string {
	return "open hof's playground"
}

// CommandCallback is invoked when the user runs your module
// return the object you want in mount or refresh
func (P *Play) CommandCallback(args []string, context map[string]interface{}) {
	if context == nil {
		context = make(map[string]any)
	}
	context["args"] = args

	if P.IsMounted() {
		// just refresh with new args
		P.Refresh(context)
	} else {
		// need to navigate, mount will do the rest
		context["path"] = "/play"
		go tui.SendCustomEvent("/router/dispatch", context)
	}
}

func (P *Play) Mount(context map[string]any) error {
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

func (P *Play) Unmount() error {
	// this is where we can do some unloading, depending on the application
	//P.View.Unmount()
	//P.Eval.Unmount()
	P.Flex.Unmount()

	return nil
}

func (P *Play) Refresh(context map[string]any) error {

	// this is where you update data and set in components
	// then at the end call tui.Draw()

	return nil
}
