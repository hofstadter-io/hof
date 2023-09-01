package cmd

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"time"

	// "golang.org/x/crypto/ssh/terminal"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/app"
	"github.com/hofstadter-io/hof/lib/tui/events"

	"github.com/hofstadter-io/hof/lib/tui/modules"
)

func Cmd(args []string, rflags flags.RootPflagpole) error {
	// stuff to ensure we don't mess up the user's terminal
	//oldState, err := terminal.MakeRaw(0)
	//if err != nil {
	//  return err
	//}
	//defer terminal.Restore(0, oldState)

	// setup new app 
	App, err := app.NewApp()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	tui.SetApp(App)

	// initialize our modules
	modules.Init()

	// Set the root view
	root := modules.RootView()
	App.SetRootView(root)

	// Ctrl-c to quit program
	tui.AddGlobalHandler("/sys/key/C-c", func(e events.Event) {
		App.Stop()
	})

	// Log Key presses (if you want to)
	logKeys()

	// Run PProf (useful for catching hangs)
	// go runPprofServer()

	// catch panics and exit, vermui will catch, clean up, format error, print, and repanic
	defer func() {
		err := recover()
		if err != nil {
			App.Stop()
			panic(err)
		}
	}()

	path := "home"
	if len(args) > 0 {
		path = args[0]
		args = args[1:]
	}

	d := map[string]any{
		"path": path,
		"args": args,
	}

	go func() {
		// some latent locksups occur randomly
		time.Sleep(time.Millisecond * 10)
		tui.SendCustomEvent("/router/dispatch", d)
		tui.SendCustomEvent("/status/message", "Welcome to [lime]VermUI[white]!!")
	}()

	// Start the Main (Blocking) Loop
	return App.Start()
}

func logKeys() {
	tui.AddGlobalHandler("/sys/key", func(e events.Event) {
		if k, ok := e.Data.(events.EventKey); ok {
			go tui.SendCustomEvent("/console/info", "key: " + k.KeyStr)
		}
	})
}

func runPprofServer() {
	runtime.SetMutexProfileFraction(1)
	http.ListenAndServe(":8888", nil)
}
