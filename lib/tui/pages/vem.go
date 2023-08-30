package pages

import (
	"os"

	"github.com/rivo/tview"

	"github.com/hofstadter-io/hof/lib/tui/app"
	"github.com/hofstadter-io/hof/lib/tui/components"
	"github.com/hofstadter-io/hof/lib/tui/layouts"
)

type VemPage struct {
	*Page

	App *app.App

	Tree *components.FileBrowser
	Text *components.TextEditor
	Repl *components.Shell
}

func NewVemPage(app *app.App, dir string) *VemPage {
	page := &VemPage{
		Page: new(Page),
		App:  app,
	}

	// setup shell
	page.Repl = components.NewShell(app)
	app.Logger = page.Repl.Append

	// setup text editor
	onTextUpdate := func() {
		// page.App.Logger("onTextUpdate")
	}
	page.Text = components.NewTextEditor(app, onTextUpdate)

	// setup file browser
	onFileSelect := func(path string) {
		page.Text.OpenFile(path)
	}
	page.Tree = components.NewFileBrowser(app, dir, onFileSelect)


	page.Name = "vem"
	layout := layouts.NewDefaultLayout(
		page.Tree, page.Text, page.Repl,
	)
	page.Flex = layout.Flex

	return page
}

func (P *VemPage) Focus(delegate func(p tview.Primitive)) {
	P.App.Logger("VemPage.Focus\n")
	info, _ := os.Lstat(P.Tree.Dir)
	if info.IsDir() {
		delegate(P.Tree)
	} else {
		delegate(P.Text)
	}

}
