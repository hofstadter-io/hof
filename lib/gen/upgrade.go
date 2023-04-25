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

func (G *Generator) upgradeDMs(dms []*datamodel.Datamodel) error {

	val := G.CueValue

	gNs, err := hof.FindHofs(val)
	if err != nil {
		return err
	}

	fmt.Println(G.Hof.Path, len(dms), len(dms[0].Node.T.History()), len(gNs))

	// assert that there is only 1 gN
	if len(gNs) != 1 {
		return fmt.Errorf("%s at %s created multple $hof Nodes, you should not mix things like generators and datamodels withing the (exactly) same value, nesting is ok", G.Hof.Label, G.Hof.Path)
	}

	gN := gNs[0]
	gN.Print()

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

	// here, we need to inject any datamodel history(s)
	// this is going to need some fancy, recursive processing
	// so we will call out to a helper of some kind
	if hn.Hof.Datamodel.History {
		G.injectHistory(hn, dms, root)
	}

	// here we create an ordered version of the node at the same level
	if hn.Hof.Datamodel.Ordered {
		G.injectOrdered(hn, dms, root)
	}

	//
	// what about the diffs? (or more generally the lens)
	//

	// recursion
	for _, c := range hn.Children {
		G.upgradeDMsR(c, dms, root)
	}
}

func (G *Generator) injectHistory(hn *hof.Node[any], dms []*datamodel.Datamodel, root *datamodel.Datamodel) {
	if root == nil {
		fmt.Println(noRootFmt, "@history", hn.Hof.Path)
		return
	}

	fmt.Println("found @history at: ", hn.Hof.Path)
	if hn.Hof.Datamodel.Root {
		hist := root.Node.T.History()
		fmt.Println("injecting hist at: ", hn.Hof.Path, len(hist), hist[0].Timestamp)
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


		fmt.Println(start, p)
		data := []map[string]any{}
		for _, h := range hist {
			d := map[string]any{
				"Timestamp": h.Timestamp,
				"Pos": h.Pos,
				"Data": h.Data,
			}
			fmt.Println(h.Lense.CurrDiff)
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
			fmt.Println("warning: root datamodel not found for child with history at, please pass the full datamodel to In values and select out later", hn.Hof.Path)
		} else {
			fmt.Println("TODO: non-root datamodel history injection not implemented yet", hn.Hof.Path)
		}
	}
}

func (G *Generator) injectOrdered(hn *hof.Node[any], dms []*datamodel.Datamodel, root *datamodel.Datamodel) error {
	if root == nil {
		return fmt.Errorf(noRootFmt, "@ordered", hn.Hof.Path)
	}
	fmt.Println("found @ordered at: ", hn.Hof.Path)

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
