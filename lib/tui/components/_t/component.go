{{ $Name := .name | title -}}
{{ $name := .name | lower -}}
package components

import (
	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/tview"
)


type {{ $Name }} struct {
	*tview.Flex // or whatever you want

	// whatever else you want
}

func New{{ $Name }}(/* ... */) *{{ $Name }} {
	c := &{{ $Name }} {
		Flex: tview.NewFlex(),
	}

	// other first time / layout setup

	return c
}

func (C *{{ $Name }}) Focus(delegate func(p tview.Primitive)) {
	delegate(C.Flex)
	// this is where you can choose how to focus
}

func (C *{{ $Name }}) Mount(context map[string]any) error {

	C.Flex.Mount(context)

	// do any setup, mount any subcomponents

	return nil
}

func (C *{{ $Name }}) Unmount() error {

	// do any teardown, unmount any subcomponents

	C.Flex.Unmount()

	return nil
}
