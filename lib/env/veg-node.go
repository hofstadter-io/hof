package env

import "github.com/hofstadter-io/hof/lib/hof"

type Env struct {
	*hof.Node[any]

	Name        string
	HumanName   string
	MachineName string

	Kind string

	Description        string
	HumanDescription   string
	MachineDescription string

	// stuff to match docker / compose / k8s
}
