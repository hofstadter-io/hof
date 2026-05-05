package runtime

import (
	"fmt"
	"time"

	"github.com/hofstadter-io/hof/lib/env"
)

type EnvEnricher func(*Runtime, *env.Env) error

func (R *Runtime) EnrichEnv(envs []string, enrich EnvEnricher) error {
	start := time.Now()
	defer func() {
		end := time.Now()
		R.Stats.Add("enrich/env", end.Sub(start))
	}()

	if R.Flags.Verbosity > 1 {
		fmt.Println("Runtime.Env: ", envs)
		for _, node := range R.Nodes {
			node.Print()
		}
	}

	// Find only the datamodel nodes
	// TODO, dedup any references
	cs := []*env.Env{}
	for _, node := range R.Nodes {
		// check for DM root
		if node.Hof.Env.Root {

			cs = append(cs, &env.Env{Node: node})
		}
	}

	R.Envs = cs

	for _, c := range R.Envs {
		err := enrich(R, c)
		if err != nil {
			return err
		}
	}

	return nil
}
