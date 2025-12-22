package runtime

import (
	"fmt"
	"time"

	"github.com/hofstadter-io/hof/lib/agent"
)

type AgentEnricher func(*Runtime, *agent.Agent) error

func (R *Runtime) EnrichAgent(agents []string, enrich AgentEnricher) error {
	start := time.Now()
	defer func() {
		end := time.Now()
		R.Stats.Add("enrich/agent", end.Sub(start))
	}()

	if R.Flags.Verbosity > 1 {
		fmt.Println("Runtime.Agent: ", agents)
		for _, node := range R.Nodes {
			node.Print()
		}
	}

	// Find only the datamodel nodes
	// TODO, dedup any references
	cs := []*agent.Agent{}
	for _, node := range R.Nodes {
		// check for DM root
		if node.Hof.Agent.Root {

			cs = append(cs, &agent.Agent{Node: node})
		}
	}

	R.Agents = cs

	for _, c := range R.Agents {
		err := enrich(R, c)
		if err != nil {
			return err
		}
	}

	return nil
}
