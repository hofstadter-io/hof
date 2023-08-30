package cuecmd

import (
	"golang.org/x/crypto/ssh/terminal"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/runtime"
	"github.com/hofstadter-io/hof/lib/tui/app"
	"github.com/hofstadter-io/hof/lib/tui/pages"
)

func runTUI(R *runtime.Runtime, cflags flags.EvalFlagpole) (err error) {
	// stuff to ensure we don't mess up the user's terminal
	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		return err
	}
	defer terminal.Restore(0, oldState)

	// application
	App := app.NewApp()

	App.Runtime = R

	// setup pages
	App.Pages = App.Pages.
		AddPage("eval", pages.NewEvalPage(App), true, true)

	err = App.SetRoot(App.Pages, true).Run()

	return err
}
