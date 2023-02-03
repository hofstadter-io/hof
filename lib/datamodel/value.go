package datamodel

import (
	"cuelang.org/go/cue"

	"github.com/hofstadter-io/hof/lib/hof"
)

type Value struct {
	*hof.Node[Value]

	// evaluated & concrete (for code gen, later)
	cval cue.Value
	data map[string]any

	// If configured, order fields via CUE
	//   for stable code generation, but better ux
	// $hof: datamodel: ordered: "path.to.struct"
	//   "." is this struct
	// ... Do we even need to save this, or create on demand?
	// ... users can always create it themselves in CUE
	// ... this seems like the better way to go
	// ... make DM concrete, then make Ordered*
	// OrderedValues []cue.Value
	// maybe we expose an iterator?

	// curr & lineage
	Snapshot  *Snapshot

	// history is only on the most current Value
	history   []*Snapshot
}

func (V *Value) Status() string {
	if len(V.history) == 0 {
		return "no-history"
	}

	// TODO... dirty or version
	if V.hasDiffR() {
		return "dirty"
	}

	return "ok"
}

// History returns the full history for a Value
func (V *Value) History() []*Snapshot {
	return V.history
}

