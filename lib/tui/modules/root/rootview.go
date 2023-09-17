package root

import (
	"fmt"
	"reflect"

	"github.com/gdamore/tcell/v2"

	"github.com/hofstadter-io/hof/lib/connector"

	"github.com/hofstadter-io/hof/lib/tui"
	"github.com/hofstadter-io/hof/lib/tui/events"
	"github.com/hofstadter-io/hof/lib/tui/tview"
	"github.com/hofstadter-io/hof/lib/tui/hoc/cmdbox"
	"github.com/hofstadter-io/hof/lib/tui/hoc/console"
	"github.com/hofstadter-io/hof/lib/tui/hoc/layouts/panels"
	"github.com/hofstadter-io/hof/lib/tui/hoc/statusbar"
)

type RootView struct {
	*panels.Layout

	//
	// Top Panel elements
	//
	// Always Visible
	cbox *cmdbox.CmdBoxWidget
	sbar *statusbar.StatusBar
	// Hidden
	errConsole *console.ErrConsoleWidget

	//
	// Main Panel element
	//
	pages *tview.Pages
	currPage string

	// hmm...
	lastCmd string

	//
	// Bottom Panel elements
	//
	devConsole *console.DevConsoleWidget
}

func New() *RootView {

	V := &RootView{
		Layout: panels.New(),
		pages: tview.NewPages(),
	}

	V.SetDirection(tview.FlexRow)

	V.buildTopPanel()

	V.buildMainPanel()

	V.buildBotPanel()

	return V
}

func (V *RootView) Connect(C connector.Connector) {
	cmds := C.Get((*cmdbox.Command)(nil))
	for _, Cmd := range cmds {
		cmd := Cmd.(cmdbox.Command)
		// fmt.Println("Command: ", cmd.CommandName())
		V.cbox.AddCommand(cmd)
	}

	V.cbox.AddCommand(&cmdbox.DefaultCommand{
		Name: "q",
		Usage: "q",
		Help: "quit hof",
		Callback: func(context map[string]any) {
			tui.Log("trace", "got here") // turns out we don't, so we have a havk below and tech debt to figure out here
			tui.Stop()
		},
	})

	cs := C.Get((*Commander)(nil))
	for _, c := range cs {
		cc := c.(Commander)
		V.pages.AddPage(cc.CommandName(), cc, true, false)
	}
}

type Commander interface {
	tview.Primitive
	CommandName() string
}

func (V *RootView) getLastCommand() (cmd string) {
	// tui.Log("trace", "GET LAST CMD: " + V.lastCmd)
	return V.lastCmd
}

func (V *RootView) setLastCommand(cmd string) {
	// tui.Log("trace", "SET LAST CMD: " + cmd)
	V.lastCmd = cmd
}

func (V *RootView) Mount(context map[string]any) error {
	tui.Log("debug", fmt.Sprintf("RootView.Mount %v", context))

	tui.AddGlobalHandler("/cmdbox/:cmd", func(e events.Event) {
		// this is wonky, why nested events
		ev, ok := e.Data.(*events.EventCustom)
		ctx, cok := ev.Data().(map[string]any)
		tui.Log("alert", fmt.Sprintf("RootView.CmdHandle %v %v %v %v %v", reflect.TypeOf(e.Data), e.Data, ok, cok, ctx))
		if ok {
			_ = V.tryPageChange(ctx)
		}
	})

	// set first page on mount
	_ = V.tryPageChange(context)

	// mount sub components
	V.Layout.Mount(context)
	// V.cbox.Mount(context)
	// V.pages.Mount(context)

	return nil
}

func (V *RootView) tryPageChange(context map[string]any) error {
	page := V.currPage
	if p, ok := context["page"]; ok {
		s := p.(string)
		page = s
	}

	// shouldn't happen, throw error
	if page == "" {
		return fmt.Errorf("missing page parameter in: %v", context)
	}

	if page == "q" {
		tui.Stop()
		return nil
	}

	// don't need to navigate
	if page == V.currPage {
		return nil
	}

	unmount := func() {
		cp := V.pages.GetPage(V.currPage)
		if cp != nil {
			V.pages.HidePage(V.currPage)
			cp.Item.Unmount()
		}
	}

	if p := V.pages.GetPage(page); p != nil {
		unmount()
		// probably a better function for switching pages
		V.currPage = page
		V.lastCmd = page
		V.pages.ShowPage(page)
		p.Item.Mount(context)
	} else {
		unmount()
		V.currPage = "not found"
		V.lastCmd = ""
		V.pages.ShowPage("not found")
		return fmt.Errorf("unknown page: %s", page)
	}

	return nil
}

func (V *RootView) buildTopPanel() {
	V.cbox = cmdbox.New(V.getLastCommand, V.setLastCommand)
	V.cbox.
		SetTitle(" [gold::bi]_[ivory]Hofstadter[-::-] ").
		SetTitleAlign(tview.AlignLeft).
		SetTitleColor(tcell.ColorIvory).
		SetBorder(true).
		SetBorderColor(tcell.ColorDodgerBlue)
	V.cbox.SetFieldTextColor(tcell.ColorIvory)

	V.sbar = statusbar.New()
	V.sbar.SetBorderColor(tcell.ColorIvory)

	// topBar is a Flex with 2 columns
	topBar := tview.NewFlex().SetDirection(tview.FlexColumn)
	topBar.AddItem(V.cbox, 0, 1, false)
	topBar.AddItem(V.sbar, 0, 1, true)

	// error console
	// V.errConsole = console.NewErrConsoleWidget()

	// Top Panels
	V.AddFirstPanel("top-bar", topBar, 3, 0, 0, "", false, "")
	// V.AddFirstPanel("err-console", V.errConsole, 0, 1, 0, "", true, "A-z")

}

func (V *RootView) buildMainPanel() {
	// A Horizontal Layout with a Router as the main element
	V.SetMainPanel("main-panel", V.pages, 0, 1, 0, "A- ")
}

func (V *RootView) buildBotPanel() {
	// dev console
	V.devConsole = console.NewDevConsoleWidget()

	// Bottom Panels
	V.AddLastPanel("dev-console", V.devConsole, 0, 1, 1, "", true, "A-/")
}
