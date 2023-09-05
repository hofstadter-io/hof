package cmd

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"strconv"
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
	//sigc := make(chan os.Signal, 10)
	//signal.Notify(sigc,
	//    // syscall.SIGHUP,
	//    // syscall.SIGTERM,
	//    syscall.SIGINT,
	//    syscall.SIGQUIT,
	//    syscall.SIGURG,
	//)
	//go func() {
	//  for {
	//    s := <-sigc
	//    fmt.Printf("got sig: %v, %d\n", s, s)
	//  }
	//}()

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
	tui.AddGlobalHandler("/sys/key/C-A-c", func(e events.Event) {
		App.Stop()
	})

	// Log Key presses (if you want to)
	if tkl := os.Getenv("HOF_TUI_KEYLOGGER"); tkl != "" {
		tklb, _ := strconv.ParseBool(tkl)
		if tklb {
			logKeys()
		}
	}

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


	// some special cases to deal with no-args startup
	path := "eval"
	if len(args) > 0 {
		path = args[0]
		args = os.Args[2:]
	} else {
		args = []string{"eval"}
	}

	// some special case for `hof tui eval` vs `hof eval --tui`
	if len(args) == 1 && args[0] == "eval" {
		// start on the help screen
		args = append(args, "help")
	}

	d := map[string]any{
		"path": path,
		"args": args,
	}

	fmt.Println(path, args)

	go func() {
		// some latent locksups occur randomly
		time.Sleep(time.Millisecond * 23)
		tui.SendCustomEvent("/router/dispatch", d)
		tui.SendCustomEvent("/status/message", "Welcome to [gold]_[ivory]Hofstadter[-]!!")
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
	tui.AddGlobalHandler("/sys/mouse", func(e events.Event) {
		if k, ok := e.Data.(events.EventMouse); ok {
			b := k.Buttons()
			if 0 < b && b < 256 {
				go tui.SendCustomEvent("/console/info", fmt.Sprintf("key: %d", k.Buttons()))
			}
		}
	})
}

func runPprofServer() {
	runtime.SetMutexProfileFraction(1)
	http.ListenAndServe(":8888", nil)
}
