package dash

import (
	"github.com/hofstadter-io/hof/lib/tui/tview"
)

type Dash struct {
	*tview.Frame

	Flex *tview.Flex

}

func New() *Dash {
	return &Dash{
		Frame: tview.NewFrame(),
		Flex:  tview.NewFlex(),
	}
}

func (C *Dash) Focus(delegate func(p tview.Primitive)) {
	// this is where you can choose how to focus

	// we just delegate to the Flex
	delegate(C.Flex)
}

func (C *Dash) Mount(context map[string]any) error {

	C.Flex.Mount(context)

	// do any setup, mount any subcomponents

	return nil
}

func (C *Dash) Unmount() error {

	// do any teardown, unmount any subcomponents

	C.Flex.Unmount()

	return nil
}
