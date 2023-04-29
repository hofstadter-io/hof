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

	err = G.upgradeDMsR(gN, dms, nil)
	if err != nil {
		return err
	}

	return nil
}

func (G *Generator) upgradeDMsR(hn *hof.Node[any], dms []*datamodel.Datamodel, root *datamodel.Datamodel) error {
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
		return nil
	}

	// here, we need to inject any datamodel history(s)
	// this is going to need some fancy, recursive processing
	// so we will call out to a helper of some kind
	if hn.Hof.Datamodel.History {
		err := G.injectHistory(hn, dms, root)
		if err != nil {
			return err
		}
	}

	// here we create an ordered version of the node at the same level
	if hn.Hof.Datamodel.Ordered {
		err := G.injectOrdered(hn, dms, root)
		if err != nil {
			return err
		}
	}

	return nil
}

func (G *Generator) injectHistory(hn *hof.Node[any], dms []*datamodel.Datamodel, root *datamodel.Datamodel) error {
	if root == nil {
		return fmt.Errorf(noRootFmt, "@history", hn.Hof.Path)
	}

	if G.Verbosity > 0 {
		fmt.Println("found @history at: ", hn.Hof.Path, root.Hof.Path)
	}

	// We want to walk the root node tree to find where it aligns with the current hn.
	// In this way, we can write code that is ignorant of where in the node tree it is.
	match := findHistoryMatchR(hn, root.Node)

	// this should not happen because we already verified that we are in the root
	// so return an error
	if match == nil {
		// TODO return error
		// fmt.Println("  no macth found")
		return nil
	}


	// get & check history
	hist := match.T.History()
	if G.Verbosity > 0 {
		fmt.Println("injecting hist at: ", hn.Hof.Metadata.ID, match.Hof.Metadata.ID, len(hist), hist[0].Timestamp)
	}

	// build up the label
	p := hn.Hof.Path

	// trim the datamodel label since we are already in there via G.Value
	start := G.Hof.Label + "."
	p = strings.TrimPrefix(p, start)

	// This is the current snapshot, outside the history object
	// Inject the CurrDiff, if the model is dirty
	if match.T.Snapshot.Lense.CurrDiff.Exists() {
		s, err := snapshotToData(match.T.Snapshot)
		if err != nil {
			return err
		}
		G.CueValue = G.CueValue.FillPath(cue.ParsePath(p+".Snapshot"), s)
	}


	// fmt.Println(start, p)
	// Datafy each snapshot
	snaps := []any{}
	for _, h := range hist {
		s, err := snapshotToData(h)
		if err != nil {
			return err
		}
		snaps = append(snaps, s)
	}

	// Inject the value at the current path as "History" list
	p += ".History"
	G.CueValue = G.CueValue.FillPath(cue.ParsePath(p), snaps)
	// fmt.Println(G.CueValue)

	return nil
}

func findHistoryMatchR(hn *hof.Node[any], root *hof.Node[datamodel.Value]) *hof.Node[datamodel.Value] {
	// are we currently there?
	// TODO: make this check better
	// current limitation is no shared names in the same datamodel tree
	//   (ID really, but when not set, then ID = name)
	//   so this could suffice for a while if we tell users to set the ID in this case
	//   maybe we can just force this by having a check somewhere during loading
	if root.Hof.Metadata.ID == hn.Hof.Metadata.ID {
		return root
	}

	// recurse if not there yet, we fully unwind on first match
	// (which is where the naming issue above comes from)
	for _, c := range root.Children {
		m := findHistoryMatchR(hn, c)
		if m != nil {
			return m
		}
	}

	return nil
}

func snapshotToData(snap *datamodel.Snapshot) (any, error) {
	s := make(map[string]any)
	// true for history snapshot entries
	// false for the snapshot we use on the current value to hold the current diff (dirty datamodel)
	if snap.Timestamp != "" {
		s["Timestamp"] = snap.Timestamp
		s["Pos"] = snap.Pos
		s["Data"] = snap.Data
	}

	// check to see if this snapshot has a diff
	// (true for all but the "first" (in time)
	if snap.Lense.CurrDiff.Exists() {
		// TODO, add more diff types & formats here
		s["CurrDiff"] = snap.Lense.CurrDiff
		// s["DownDiff"] = snap.Lense.DownDiff
	}

	return s, nil
}

// the point of this is to have a stable order for cue values specified as a struct
// since they get turned into Go maps, which have random order during iteration
// generated code can shift around while being the "same"
// This is where we auto-fill from @ordered(), but users can also do this manually
// Note | XXX, CUE's order may change between versions, they are working towards defining a stable order
//   at which point we will use the same for consistency. We should be backwards compatible at this point
//   but there is risk until then
func (G *Generator) injectOrdered(hn *hof.Node[any], dms []*datamodel.Datamodel, root *datamodel.Datamodel) error {
	if root == nil {
		return fmt.Errorf(noRootFmt, "@ordered", hn.Hof.Path)
	}

	if G.Verbosity > 0 {
		fmt.Println("found @ordered at: ", hn.Hof.Path)
	}

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
