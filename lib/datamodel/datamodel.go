package datamodel

import (
	"github.com/hofstadter-io/hof/lib/hof"
)

// this is basically the same as a Value
// except that it reporesents a conceptual root
// and we want specific functions on it
// that are different from general Nodes'
// handling and recursion
// type Datamodel *hof.Node[Value]
type Datamodel struct {
	*hof.Node[Value]
}

func DatamodelType(DM *Datamodel) string {
	// if explicitly set to CUE value
	//   todo, can we look for incomplete values?
	//   seems problematic, can't separate bad config
	//   so probably not, but leaving this comment here
	if DM.Hof.Datamodel.Cue {
		return "value"
	}

	// if history at root & no children ... a bit hacky, but will do
	//   is Root correct here? what about just hist & no children? (all leaf nodes are then objects?
	// if DM.Hof.Datamodel.Root && DM.Hof.Datamodel.History && len(DM.Children) == 0 {
	if DM.Hof.Datamodel.History && len(DM.Children) == 0 {
		return "object"
	}

	// otherwise generic datamodel
	return "datamodel"
}

func (dm *Datamodel) Status() string {
	if has, _ := dm.HasHistory(); !has {
		return "no-history"
	}

	// TODO... dirty or version
	if dm.HasDiff() {
		return "dirty"
	}

	return "ok"
}
