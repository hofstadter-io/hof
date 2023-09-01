package components

import (
	"github.com/hofstadter-io/hof/lib/tui/app"
	"github.com/hofstadter-io/hof/lib/tui/tview"
)

type Shell struct {
	*tview.TextArea

	text string

	App *app.App
}

func NewShell(app *app.App) *Shell {
	s := &Shell{
		TextArea: tview.NewTextArea(),
		App: app,
	}

	// lower-level setup
	s.SetTitle("Shell").
		SetBorder(true)

	return s
}

func (S *Shell) Write(text string) {
	S.text = text
	S.SetText(S.text, true)
}

func (S *Shell) Append(text string) {
	S.text += text
	S.SetText(S.text, true)
}
