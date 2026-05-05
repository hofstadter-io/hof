package agent

import "github.com/hofstadter-io/hof/lib/hof"

type Agentic struct {
	*hof.Node[any]

	Name        string
	HumanName   string
	MachineName string

	Description        string
	HumanDescription   string
	MachineDescription string

	// stuff to match docker / compose / k8s
}
