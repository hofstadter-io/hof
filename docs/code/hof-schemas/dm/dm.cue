package dm

import (
	"github.com/hofstadter-io/hof/schema"
)

// This is a complete Value tracked as one
// useful for schemas, config, and NoSQL
Object: {
	schema.Hof// needed for reFerences
	#hof: datamodel: root: true

	TrackHistory

	// all fields will be tracked
}

// This is like object, but supports cue values
// (todo, should support full lattice)
Value: {
	Object
	#hof: datamodel: cue: true
}

// This is a general datamodel useful in many applications
// It can be expanded and enriched to cover more
// Useful for SQL, APIs, forms, and similar
Datamodel: {
	schema.DHof// needed for reFerences
	#hof: datamodel: root: true
}

// Schema for a snapshot, can include anything else
Snapshot: {
	Timestamp: string | *""
}

// embedable history type
History: [...Snapshot]

TrackHistory: {
	#hof: datamodel: history: true // needed for CUE compat
	"Snapshot": Snapshot
	"History":  History
}
