package gen

import (
	"fmt"
	"strings"

	"cuelang.org/go/cue"

	"github.com/hofstadter-io/hof/lib/hof"
	"github.com/hofstadter-io/hof/lib/datamodel"
)

const noRootFmt = `warning: root datamodel not found for child with %s at %s
please pass the full datamodel to In values and select out later`

// TODO, we may be able to move the upgrade logic to the datamodel package
// and then unify the value with the ones we discover here.
// The reason we are doing it here now is that the references get lost during decoding
// and any changes we make to the datamodel are not reflected within the In value of the Genartor
// Ideally, we could do this once and update values as needed.
// It is likely we are doing work more than once.
// We also need to ensure some order of operations here,
//   as we need to have the history injected prior to decoding the majority of the Generator
//   so that users can reference the history within the generator.
//   The notable use case is making a file per snapshot for database migrations.

func (G *Generator) upgradeDMs(dms []*datamodel.Datamodel) error {
	if len(dms) == 0 {
		return nil
	}

	val := G.CueValue

	// build a hof.Node tree from the gen.CueValue
	gNs, err := hof.FindHofs(val)
	if err != nil {
		return err
	}

	// fmt.Println(G.Hof.Path, len(dms), len(dms[0].Node.T.History()), len(gNs))

	// assert that there is only 1 gN
	if len(gNs) != 1 {
		return fmt.Errorf("%s at %s created multiple $hof Nodes, you should not mix things like generators and datamodels withing the (exactly) same value, nesting is ok", G.Hof.Label, G.Hof.Path)
	}

	gN := gNs[0]
	// gN.Print()

	G.upgradeDMsR(gN, dms, nil)

	return nil
}

func (G *Generator) upgradeDMsR(hn *hof.Node[any], dms []*datamodel.Datamodel, root *datamodel.Datamodel) {
	if root == nil && hn.Hof.Datamodel.Root {
		for _, dm := range dms {
			if dm.Hof.Metadata.Name == hn.Hof.Metadata.Name {
				root = dm
			}
		}
	}

	// recursion into children
	for _, c := range hn.Children {
		G.upgradeDMsR(c, dms, root)
	}
	// we do all the real work post-order recursion
	// so that parent enrichments include child enrichments

	// return if we are not within the DM root
	if root == nil {
		return
	}

	// here, we need to inject any datamodel history(s)
	// this is going to need some fancy, recursive processing
	// so we will call out to a helper of some kind
	if hn.Hof.Datamodel.History {
		G.injectHistory(hn, dms, root)
	}

	//
	// what about the diffs? (or more generally the lens)
	//

	// here we create an ordered version of the node at the same level
	if hn.Hof.Datamodel.Ordered {
		G.injectOrdered(hn, dms, root)
	}

}

func (G *Generator) injectHistory(hn *hof.Node[any], dms []*datamodel.Datamodel, root *datamodel.Datamodel) {
	if root == nil {
		fmt.Println(noRootFmt, "@history", hn.Hof.Path)
		return
	}

	// fmt.Println("found @history at: ", hn.Hof.Path)

	// Instead of this check, we should try to walk the root to find where it aligns with the current hn.
	// In this way, we can hopefully write code that is ignorant of where in the node tree it is.

	if hn.Hof.Datamodel.Root {
		hist := root.Node.T.History()
		// fmt.Println("injecting hist at: ", hn.Hof.Path, len(hist), hist[0].Timestamp)
		p := hn.Hof.Path
		start := G.Hof.Label + "."
		p = strings.TrimPrefix(p, start)

		if root.Node.T.Snapshot.Lense.CurrDiff.Exists() {
			// fmt.Println("curr diff:", root.Node.T.Snapshot.Lense.CurrDiff)
			data := map[string]any{
				"CurrDiff": root.Node.T.Snapshot.Lense.CurrDiff,
			}
			G.CueValue = G.CueValue.FillPath(cue.ParsePath(p), data)
		}


		// fmt.Println(start, p)
		data := []map[string]any{}
		for _, h := range hist {
			d := map[string]any{
				"Timestamp": h.Timestamp,
				"Pos": h.Pos,
				"Data": h.Data,
			}
			// fmt.Println(h.Lense.CurrDiff)
			if h.Lense.CurrDiff.Exists() {
				d["CurrDiff"] = h.Lense.CurrDiff
			}
			data = append(data, d)
		}
		p += ".History"
		G.CueValue = G.CueValue.FillPath(cue.ParsePath(p), data)
		// fmt.Println(G.CueValue)
	} else {
		if root == nil {
			// The main problem here is knowing what datamodel a nested DM identirfier is part of.
			// It is trivial to construct an example where the same name appears in two different datamodels.

			fmt.Println("warning: root datamodel not found for child with history at, please pass the full datamodel to In values and select out later", hn.Hof.Path)
		} else {
			if G.Verbosity > 0 {
				fmt.Println("TODO: non-root datamodel history injection not implemented yet", hn.Hof.Path)
			}
		}
	}
}

func (G *Generator) injectOrdered(hn *hof.Node[any], dms []*datamodel.Datamodel, root *datamodel.Datamodel) error {
	if root == nil {
		return fmt.Errorf(noRootFmt, "@ordered", hn.Hof.Path)
	}
	// fmt.Println("found @ordered at: ", hn.Hof.Path)

	path := hn.Hof.Path
	path = strings.TrimPrefix(path, G.Name + ".")
	value := G.CueValue.LookupPath(cue.ParsePath(path))

	iter, err := value.Fields()
	if err != nil {
		return err
	}

	var ordered []cue.Value
	for iter.Next() {
		sel := iter.Selector().String()
		if sel == "$hof" {
			continue
		}

		val := iter.Value()
		ordered = append(ordered, val)
		// fmt.Printf("%s: %v\n", sel, iter.Value())
	}

	// fmt.Printf("%# v\n", val)

	// construct new ordered list for fields
	l := value.Context().NewList(ordered...)

	// fill into Gen value
	G.CueValue = G.CueValue.FillPath(cue.ParsePath(path + "Ordered"), l)
	return nil
}
