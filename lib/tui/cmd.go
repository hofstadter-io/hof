package tui

import (
	"os"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/tui/app"
	"github.com/hofstadter-io/hof/lib/tui/pages"
)



func Cmd(args []string, rflags flags.RootPflagpole) error {
	// stuff to ensure we don't mess up the user's terminal
	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		return err
	}
	defer terminal.Restore(0, oldState)

	// only arg is the default file to open, for now...
	path := ""
	if len(args) > 0 {
		path = args[0]
	}
	if path == "" {
		p, err := os.Getwd()
		if err != nil {
			return err
		}
		path = p
	}

	// setup new app 
	App := app.NewApp()

	// setup pages
	App.Pages = App.Pages.
		AddPage("vem", pages.NewVemPage(App, path), true, true)

	return App.SetRoot(App.Pages, true).Run()
}





//func Cmd(args []string, rflags flags.RootPflagpole) error {

//  rl, err := readline.New("> ")
//  if err != nil {
//    panic(err)
//  }
//  defer rl.Close()

//  for {
//    line, err := rl.Readline()
//    if err != nil { // io.EOF
//      break
//    }
//    line = strings.TrimSpace(line)
//    fmt.Println(line)
//  }

//  return nil
//}
