package common

import (
	"fmt"
	"os"
	"io"

	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/quick"
	"github.com/gdamore/tcell/v2"

	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/tview"
)

type TextEditor struct {
	*tview.TextView

	W io.Writer

	OnChange func()
}

func NewTextEditor(onchange func()) *TextEditor {

	te := &TextEditor{
		TextView: tview.NewTextView(),
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

	body, err := os.ReadFile(path)
	if err != nil {
		tui.SendCustomEvent("/console/err", err.Error())
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
		tui.SendCustomEvent("/console/err", err.Error())
	}

	ED.Focus(func(p tview.Primitive){})

	ED.SetInputCapture(func(evt *tcell.EventKey) *tcell.EventKey {

		switch evt.Key() {

		case tcell.KeyRune:
			switch evt.Rune() {
				case '?':
				tui.SendCustomEvent("/console/err", err.Error())
				return nil
			default:
				return evt
			}

		default:
			return evt
		}

		// VB.Rebuild("")

		return nil
	})
}

func (ED *TextEditor) Focus(delegate func(p tview.Primitive)) {
	delegate(ED.TextView)
}
