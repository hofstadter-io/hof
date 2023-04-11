package runtime

import (
	"fmt"
	"regexp"
	"time"

	"github.com/hofstadter-io/hof/lib/gen"
	"github.com/hofstadter-io/hof/lib/hof"
)

type GeneratorEnricher func(*Runtime, *gen.Generator) error

func (R *Runtime) EnrichGenerators(generators []string, enrich GeneratorEnricher) error {
	if R.Flags.Verbosity > 1 {
		fmt.Println("Runtime.EnrichGenerators: ", generators)
		for _, node := range R.Nodes {
			node.Print()
		}
	}

	keep := func(hn *hof.Node[any]) bool {
		// filter by name
		if len(generators) > 0 {
			for _, d := range generators {
				match, err := regexp.MatchString(d, hn.Hof.Metadata.Name)
				if err != nil {
					fmt.Println("error:", err)
					return false
				}
				if match {
					return true
				}
			}
			return false
		}

		// filter by time

		// filter by version?

		// default to true, should include everything when no checks are needed
		return true
	}

	// Find only the generator nodes
	// TODO, dedup any references
	gens := []*gen.Generator{}
	for _, node := range R.Nodes {
		// check for DM root
		if node.Hof.Gen.Root {
			if !keep(node) {
				continue
			}
			upgrade := func(n *hof.Node[gen.Generator]) *gen.Generator {
				v := gen.NewGenerator(n)
				return v
			}
			u := hof.Upgrade[any, gen.Generator](node, upgrade, nil)
			// we'd like this line in upgrade, but...
			// how do we make T a Node[T] type (or ensure that it has a hof)
			// u.T.Hof = u.Hof
			gen := u.T
			gen.Node = u
			gens = append(gens, gen)
		}
	}

	R.Generators = gens

	// what do we do to enrich a generator?
	// load & validate?
	// add datamodel history to input data?

	start := time.Now()
	defer func() {
		end := time.Now()
		R.Stats.GenLoadingTime = end.Sub(start)
	}()
	for _, gen := range R.Generators {
		err := enrich(R, gen)
		if err != nil {
			return err
		}
	}

	return nil
}
