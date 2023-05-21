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

	// user inputs
	Args []string
	Files map[string]string
  Question string

	Model   string

	System   string
	Examples []Example
	Messages []Message
	Parameters map[string]any

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
