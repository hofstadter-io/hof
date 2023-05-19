package chat

import (
	"github.com/hofstadter-io/hof/lib/hof"
)

type Chat struct {
	*hof.Node[any]

	Name string
	HumanName   string
	MachineName string

	Description        string
	HumanDescription   string
	MachineDescription string

	Model   string

	Prompts  map[string]string
	Examples map[string]Example
	Messages map[string]Message
	Paramaters map[string]map[string]any

	Question string
	Output   string

	// internal
	Session Session
}

type Session struct {
	System string
	Examples []Example
	Messages []Message
}

type Example struct {
	Input  string
	Output string
}

type Message struct {
	Role    string
	Content string
}
