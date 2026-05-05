package runtime

import (
	"fmt"
	"time"

	"github.com/hofstadter-io/hof/lib/agent"
)

type AgenticEnricher func(*Runtime, *agent.Agentic) error

func (R *Runtime) EnrichAgentic(agents []string, enrich AgenticEnricher) error {
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

	// Find only the agentic nodes
	// TODO, dedup any references, extract memoization from env into runtime
	cs := []*agent.Agentic{}
	for _, node := range R.Nodes {
		// check for DM root
		if node.Hof.Agentic.Root {

			cs = append(cs, &agent.Agentic{Node: node})
		}
	}

	R.Agentics = cs

	for _, c := range R.Agentics {
		err := enrich(R, c)
		if err != nil {
			return err
		}
	}

	return nil
}
