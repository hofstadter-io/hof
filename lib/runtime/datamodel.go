package runtime

import (
	"fmt"
	"regexp"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/datamodel"
	"github.com/hofstadter-io/hof/lib/hof"
)

func (R *Runtime) FindDatamodels(dflags flags.DatamodelPflagpole) error {
	if R.Flags.Verbosity > 1 {
		for _, node := range R.Nodes {
			node.Print()
		}
	}

	keep := func(hn *hof.Node[any]) bool {
		// filter by name
		if len(dflags.Datamodels) > 0 {
			for _, d := range dflags.Datamodels {
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

	// Find only the datamodel nodes
	// TODO, dedup any references
	dms := []*datamodel.Datamodel{}
	for _, node := range R.Nodes {
		// check for DM root
		if node.Hof.Datamodel.Root {
			if !keep(node) {
				continue
			}
			t := func(n *hof.Node[datamodel.Value]) *datamodel.Value {
				v := new(datamodel.Value)
				v.Node = n
				v.Snapshot = new(datamodel.Snapshot)
				return v
			}
			u := hof.Upgrade[any, datamodel.Value](node, t, nil)
			// we'd like this line in upgrade, but...
			// how do we make T a Node[T] type (or ensure that it has a hof)
			// u.T.Hof = u.Hof
			dms = append(dms, &datamodel.Datamodel{Node: u})
		}
	}

	R.Datamodels = dms

	// filter datamodel if flag set?
	// which flags do we handle here
	//   vs in the various commands?

	// load history
	// calc diffs? (or load)

	return nil
}
