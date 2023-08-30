package components

import (
	"fmt"
	"os"
	"io"

	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/quick"
	"github.com/rivo/tview"

	"github.com/hofstadter-io/hof/lib/tui/app"
)

type TextEditor struct {
	*tview.TextView

	App *app.App
	W io.Writer

	OnChange func()
}

func NewTextEditor(app *app.App, onchange func()) *TextEditor {

	te := &TextEditor{
		TextView: tview.NewTextView(),
		App: app,
		OnChange: onchange,
	}
	te.SetWordWrap(true).
		SetDynamicColors(true).
		SetBorder(true)
	te.SetChangedFunc(te.OnChange)

	te.W = tview.ANSIWriter(te)

	return te
}

func (ED *TextEditor) OpenFile(path string) {
	ED.App.Logger("ED.OpenFile: " + path + "\n")

	body, err := os.ReadFile(path)
	if err != nil {
		ED.App.Logger("error: " + err.Error())
	}

	l := lexers.Match(path)
	lexer := "text"	
	if l != nil {
		lexer = l.Config().Name
	} else {
		var s string
		if len(body) > 512 {
			s = string(body[:512])
		} else {
			s = string(body)
		}
			
		l = lexers.Analyse(s)
		if l != nil {
			lexer = l.Config().Name
		}
	}

	ED.SetTitle(fmt.Sprintf("%s (%s)", path, lexer))

	ED.Clear()

	err = quick.Highlight(ED.W, string(body), lexer, "terminal256", "solarized-dark")
	if err != nil {
		ED.App.Logger("error: " + err.Error())
	}

	ED.Focus(func(p tview.Primitive){})

}

func (ED *TextEditor) Focus(delegate func(p tview.Primitive)) {
	ED.App.Logger("ED.Focus\n")
	delegate(ED.TextView)
}
