package runtime

import (
	"fmt"
	"time"

	"github.com/hofstadter-io/hof/flow/flow"
	"github.com/hofstadter-io/hof/lib/hof"
)

type FlowEnricher func(*Runtime, *flow.Flow) error

func (R *Runtime) EnrichFlows(flows []string, enrich FlowEnricher) error {
	start := time.Now()
	defer func() {
		end := time.Now()
		R.Stats.Add("enrich/flow", end.Sub(start))
	}()

	if R.Flags.Verbosity > 1 {
		fmt.Println("Runtime.Flow: ", flows)
		for _, node := range R.Nodes {
			node.Print()
		}
	}

	// Find only the datamodel nodes
	// TODO, dedup any references
	fs := []*flow.Flow{}
	for _, node := range R.Nodes {
		// check for Chat root
		if node.Hof.Flow.Root {

			if !keepFilter(node, flows) {
				continue
			}
			upgrade := func(n *hof.Node[flow.Flow]) *flow.Flow {
				v := flow.NewFlow(n)
				return v
			}
			u := hof.Upgrade[any, flow.Flow](node, upgrade, nil)
			// we'd like this line in upgrade, but...
			// how do we make T a Node[T] type (or ensure that it has a hof)
			// u.T.Hof = u.Hof
			f := u.T
			f.Node = u
			fs = append(fs, f)
		}
	}

	R.Workflows = fs

	for _, c := range R.Workflows {
		err := enrich(R, c)
		if err != nil {
			return err
		}
	}


	return nil
}
