package form

import (
	"github.com/gdamore/tcell/v2"
	"github.com/hofstadter-io/hof/lib/tui/tview"
)

type Form interface {
	tview.Primitive

	Name() string

	GetValues() (values map[string]interface{})
	SetValues(values map[string]interface{})

	AddItem(name string, item FormItem, taborder, proportion int)
	GetItem(name string) FormItem
	GetItems(name string) []FormItem

	AddButton(name string, button FormButton, taborder, proportion int)
	GetButton(name string) FormButton
	GetButtons(name string) []FormButton
}

type FormItem interface {
	tview.Primitive

	Name() string

	GetValues() (values map[string]interface{})
	SetValues(values map[string]interface{})

	SetFinishedFunction(handler func(key tcell.Key))
}

type FormButton interface {
	tview.Primitive

	Name() string

	OnSubmit()

	SetBlurFunction(handler func(key tcell.Key))
}

type FormBase struct {
	*Flex

	focusedElement int

	// An optional function which is called when the user hits Escape.
	cancel func()
}

func New(name string) *FormBase {
	L := &FormBase{
		Flex: NewFlex(name),
	}

	return L
}

func (f *FormBase) Focus(delegate func(p tview.Primitive)) {
	items := f.GetItems()
	buttons := f.GetButtons()

	if len(items)+len(buttons) == 0 {
		return
	}

	// Hand on the focus to one of our child elements.
	if f.focusedElement < 0 || f.focusedElement >= len(items)+len(buttons) {
		f.focusedElement = 0
	}
	handler := func(key tcell.Key) {
		switch key {
		case tcell.KeyTab, tcell.KeyEnter:
			f.focusedElement++
			f.Focus(delegate)
		case tcell.KeyBacktab:
			f.focusedElement--
			if f.focusedElement < 0 {
				f.focusedElement = len(items) + len(buttons) - 1
			}
			f.Focus(delegate)
		case tcell.KeyEscape:
			if f.cancel != nil {
				f.cancel()
			} else {
				f.focusedElement = 0
				f.Focus(delegate)
			}
		}
	}

	if f.focusedElement < len(items) {
		// We're selecting an item.
		item := items[f.focusedElement]
		item.SetFinishedFunction(handler)
		delegate(item)
	} else {
		// We're selecting a button.
		button := buttons[f.focusedElement-len(items)]
		button.SetBlurFunction(handler)
		delegate(button)
	}
}
