package schema

import (
	"github.com/hofstadter-io/hof/schema/dm/sql"
)

Datamodel: sql.Datamodel & {
	// Add our enrichment to the models
	// Models: [M=string]: Model

	Models: [M=string]: Model
}

// We will enrich our data model with this
Model: sql.Model & {
	// field used for indexing
	Index?: string

	// Adds GoType
	Fields: [F=string]: Field & {Relation?: Name: string | *F}
}

Field: sql.Field & EnrichRelation

EnrichRelation: {
	Relation?: {
		Type:   string
		Other:  string
		GoType: string

		// Switch pattern
		GoType: [
			if Type == "belongs-to" {"*\(Other)"},
			if Type == "has-one" {"*\(Other)"},
			if Type == "has-many" {"[]*\(Other)"},
			if Type == "many-to-many" {"[]*\(Other)"},
			"unknown relation type: \(Type)",
		][0]
	}
}
