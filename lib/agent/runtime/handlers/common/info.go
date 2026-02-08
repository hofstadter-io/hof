package common

import (
	aruntime "github.com/hofstadter-io/hof/lib/agent/runtime"
	"github.com/hofstadter-io/hof/lib/agent/services/environ"
	"github.com/hofstadter-io/hof/lib/cuetils"
)

func ReloadConfig(ar *aruntime.Runtime) error {
	err := ar.ReadEnvConfig()
	if err != nil {
		return cuetils.ExpandCueError(err)
	}
	return nil
}

func ListEnvirons() ([]string, error) {
	envs, err := environ.Client().ListEnvirons()
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range envs {
		names = append(names, e.Name)
	}
	return names, nil
}
