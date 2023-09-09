package root

import (
	"github.com/gdamore/tcell/v2"

	"github.com/hofstadter-io/hof/lib/connector"

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
	mainPanel *MainPanel
	lastCmd string

	//
	// Bottom Panel elements
	//
	devConsole *console.DevConsoleWidget
}

func New() *RootView {

	V := &RootView{
		Layout: panels.New(),
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

	V.mainPanel.Connect(C)
}

func (V *RootView) getLastCommand() (cmd string) {
	// tui.Log("trace", "GET LAST CMD: " + V.lastCmd)
	return V.lastCmd
}

func (V *RootView) setLastCommand(cmd string) {
	// tui.Log("trace", "SET LAST CMD: " + cmd)
	V.lastCmd = cmd
}

func (V *RootView) buildTopPanel() {
	V.cbox = cmdbox.New(V.getLastCommand, V.setLastCommand)
	V.cbox.
		SetTitle("  [gold]_[ivory]Hofstadter[-]  ").
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
	V.errConsole = console.NewErrConsoleWidget()

	// Top Panels
	V.AddFirstPanel("top-bar", topBar, 3, 0, 0, "", false, "")
	V.AddFirstPanel("err-console", V.errConsole, 0, 1, 0, "", true, "A-z")

}

func (V *RootView) buildMainPanel() {
	// A Horizontal Layout with a Router as the main element
	V.mainPanel = NewMainPanel(V.getLastCommand, V.setLastCommand)
	V.SetMainPanel("main-panel", V.mainPanel, 0, 1, 0, "A- ")
}

func (V *RootView) buildBotPanel() {
	// dev console
	V.devConsole = console.NewDevConsoleWidget()

	// Bottom Panels
	V.AddLastPanel("dev-console", V.devConsole, 0, 1, 1, "", true, "A-/")
}
