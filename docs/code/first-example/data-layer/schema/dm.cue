package schema

import (
	"github.com/hofstadter-io/hof/schema/dm/sql"
)

Datamodel: sql.Datamodel & {
	Models: [N= !="$hof"]: Model
}

Model: {
	// field used for indexing
	Index: string

	// Adds GoType
	Relations: [string]: R={
		GoType: string

		// Switch pattern
		GoType: [
			if R.Reln == "BelongsTo" {"*\(R.Type)"},
			if R.Reln == "HasOne" {"*\(R.Type)"},
			if R.Reln == "HasMany" {"[]*\(R.Type)"},
			if R.Reln == "ManyToMany" {"[]*\(R.Type)"},
			"unknown relation type: \(R.Reln)",
		][0]
	}
}
