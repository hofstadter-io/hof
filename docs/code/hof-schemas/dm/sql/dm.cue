package sql

import (
	"github.com/hofstadter-io/hof/schema"
	"github.com/hofstadter-io/hof/schema/dm/fields"
)

CommonFields: {
	ID:        fields.UUID
	CreatedAt: fields.Datetime
	UpdatedAt: fields.Datetime
}

SoftDelete: {
	DeletedAt: fields.Datetime
}

Datamodel: {
	@datamodel()

	// these are the models for the application
	// they can map onto database tables and apis
	Models: {
		// create point where hof can list, info, etc..
		$hof: datamodel: node: true
		@node()

		$hof: datamodel: ordered: true
		@ordered() // for stability, see below

		// each struct field is a Model
		[N= !="$hof"]: Model & {$hof: metadata: name: N}
	}

	// OrderedModels: [...Model] will be
	// inject here for order stability
}

Model: M={
	schema.DHof// needed for reFerences

	$hof: datamodel: history: true // needed for CUE compat
	@history() // hof only shorthand

	// Lineage fields will be filled by hof
	// $hof: Lense: ...
	// $hof: History: ...

	// for easy access
	Name: M.$hof.metadata.name

	// These are the fields of a model
	// they can map onto database columnts and form fields
	Fields: {
		// create point where hof can list, info, etc..
		$hof: datamodel: node: true
		@node()

		// for stability, see below
		$hof: datamodel: ordered: true
		@ordered() // shorthand

		[N= !="$hof"]: Field & {$hof: metadata: name: N}
	}

	// OrderedFields: [...Fields] will be
	// inject here for order stability

	// if we want Relations as a separate value
	// we can process the fields to extract them
}

Field: F={
	schema.DHof// needed for references
	$hof: datamodel: history: true
	@history() // shorthand

	// Lineage fields will be filled by hof
	// $hof: Lense: ...
	// $hof: History: ...

	// for easy access
	Name: F.$hof.metadata.name
	Type: string

	// relation type, open to be flexible
	Reln?: string

	// what about {val, *val, []val, []*val}
	// we probably don't care about pointer here
	//   that is a language detail (code gen target)

	// we can enrich this for various types
	// in our app or other reusable datamodels
}
