package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/agent/cmd/tui"
)

func Chat(args []string, rflags flags.RootPflagpole, aflags flags.AgentPflagpole, cflags flags.Agent__ChatPflagpole) error {
	// hard code the kind here, as we only support one at this point for chat
	aflags.Kind = []string{"agent"}
	R, AR, matches, err := commonStart(args, rflags, aflags)
	if err != nil {
		return err
	}

	if len(matches) != 1 {
		return fmt.Errorf("expected one match as an agent, matched %d of %d", len(matches), len(R.Agentics))
	}

	m, err := tui.InitialModel(R, AR, "chat")
	if err != nil {
		return err
	}

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}
