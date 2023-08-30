package layouts

import (
	"github.com/rivo/tview"
)

type DefaultLayout struct {
	*tview.Flex

	Main *tview.Flex

	Tree tview.Primitive
	Text tview.Primitive
	Repl tview.Primitive
}

func NewDefaultLayout(tree, text, repl tview.Primitive) *DefaultLayout {
	layout := &DefaultLayout{
		Tree: tree,
		Text: text,
		Repl: repl,
	}

	// flexbox
	layout.Flex = tview.NewFlex()
	layout.Main = tview.NewFlex().SetDirection(tview.FlexRow)

	layout.Flex.
		AddItem(layout.Tree, 42, 1, false).
		AddItem(layout.Main, 0, 2, false)

	layout.Main.
		AddItem(layout.Text, 0, 3, false).
		AddItem(layout.Repl, 10, 1, false)

	return layout
}
