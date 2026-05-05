package cmd

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/agent/cmd/tui"
)

func Main(args []string, rflags flags.RootPflagpole, aflags flags.AgentPflagpole, cflags flags.Agent__ChatPflagpole) error {
	R, AR, _, err := commonStart(args, rflags, aflags)
	if err != nil {
		return err
	}

	// HMMM, maybe we don't need to run this, it's just the server
	//   should also check in on the sqlite database and having two processes try to access it
	//   there will be a server mode if that's all you really want to do (which can serve up any #veg, but also submodules have extra endpoints, like here for dagger backed filesys)
	// run the agentic ~engine~ API server as long as the TUI is running
	// TODO, a way to shut this down cleanly, probably using context?
	// go func() {
	// 	AR.Run()
	// }()

	m, err := tui.InitialModel(R, AR, "list")
	if err != nil {
		return err
	}
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}
