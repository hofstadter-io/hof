package datamodel

import (
	"strings"

	"cuelang.org/go/cue"

	"github.com/hofstadter-io/hof/lib/hof"
)

func (dm *Datamodel) EnrichValue() error {
	// fmt.Println("dm.EnrichValue()", dm.Hof.Path)

	// val := dm.Value()

	// build a hof.Node tree from the gen.CueValue
	//nodes, err := hof.FindHofs(val)
	//if err != nil {
		//return err
	//}
	//nodes := dm.Node

	// fmt.Println(G.Hof.Path, len(dms), len(dms[0].Node.T.History()), len(gNs))

	// assert that there is only 1 gN
	//if len(nodes) != 1 {
	//  return fmt.Errorf("%s at %s created multiple $hof Nodes, you should not mix things like generators and datamodels withing the (exactly) same value, nesting is ok", dm.Hof.Label, dm.Hof.Path)
	//}

	//node := nodes[0]
	// gN.Print()

	err := dm.enrichR(dm.Node)
	if err != nil {
		return err
	}

	return nil
}

func (dm *Datamodel) enrichR(hn *hof.Node[Value]) error {

	// recursion into children
	for _, c := range hn.Children {
		err := dm.enrichR(c)
		if err != nil {
			return err
		}
	}

	// we do all the real work post-order recursion
	// so that parent enrichments include child enrichments

	// fmt.Println("DM.enrichR", dm.Hof.Path, hn.Hof.Path)

	// here, we need to inject any datamodel history(s)
	// this is going to need some fancy, recursive processing
	// so we will call out to a helper of some kind
	if hn.Hof.Datamodel.History {
		err := dm.enrichHistory(hn)
		if err != nil {
			return err
		}
	}

	// here we create an ordered version of the node at the same level
	//if hn.Hof.Datamodel.Ordered {
	//  err := dm.enrichOrdered(hn)
	//  if err != nil {
	//    return err
	//  }
	//}

	return nil
}

func (dm *Datamodel) enrichHistory(hn *hof.Node[Value]) error {
	// if G.Verbosity > 0 {
		 // fmt.Println("found @history at: ", hn.Hof.Path, hn.Hof.Metadata.ID, dm.Hof.Path, dm.Hof.Metadata.ID)
	// }

	// We want to walk the root node tree to find where it aligns with the current hn.
	// In this way, we can write code that is ignorant of where in the node tree it is.
	match := findHistoryMatchR(hn, dm.Node)

	// this should not happen because we already verified that we are in the root
	// so return an error
	if match == nil {
		// TODO return error
		// fmt.Println("  no match found")
		return nil
	}


	// get & check history
	hist := match.T.History()
	// if G.Verbosity > 0 {
		 // fmt.Println("injecting hist at: ", hn.Hof.Metadata.ID, match.Hof.Metadata.ID, len(hist), hist[0].Timestamp)
	// }

	// build up the label
	p := hn.Hof.Path

	// trim the datamodel label since we are already in there via G.Value
	start := dm.Hof.Label + "."
	p = strings.TrimPrefix(p, start)

	// This is the current snapshot, outside the history object
	// Inject the CurrDiff, if the model is dirty
	if match.T.Snapshot.Lense.CurrDiff.Exists() {
		s, err := snapshotToData(match.T.Snapshot)
		if err != nil {
			return err
		}
		dm.Value = dm.Value.FillPath(cue.ParsePath(p+".Snapshot"), s)
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
	// XXX TODO XXX inject refrence rather than value
	dm.Value = dm.Value.FillPath(cue.ParsePath(p), snaps)
	// fmt.Println(G.CueValue)

	return nil
}

func findHistoryMatchR(hn *hof.Node[Value], root *hof.Node[Value]) *hof.Node[Value] {
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

func snapshotToData(snap *Snapshot) (any, error) {
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
		s["PrevDiff"] = snap.Lense.PrevDiff
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
func (dm *Datamodel) enrichOrdered(hn *hof.Node[any]) error {
	// if G.Verbosity > 0 {
		// fmt.Println("found @ordered at: ", hn.Hof.Path)
	// }

	path := hn.Hof.Path
	path = strings.TrimPrefix(path, dm.Hof.Metadata.Name + ".")
	value := dm.Value.LookupPath(cue.ParsePath(path))

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
	dm.Value = dm.Value.FillPath(cue.ParsePath(path + "Ordered"), l)
	return nil
}
