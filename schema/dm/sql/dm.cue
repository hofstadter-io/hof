package sql

import (
	"github.com/hofstadter-io/hof/schema"
)

Datamodel: {
	$hof: datamodel: root: true

	// these are the models for the application
	// they can map onto database tables and apis
	Models: {
		$hof: datamodel: node:    true
		$hof: datamodel: ordered: true
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
		$hof: datamodel: node:    true
		$hof: datamodel: ordered: true
		[N= !="$hof"]: Field & {Name: N}
	}

	// OrderedFields: [...Fields] will be
	// inject here for order stability

	// if we want Relations as a separate value
	// we can process the fields to extract them
}

Field: {
	// TODO, decide if these should be the default
	// schema.DHof// needed for reFerences
	// $hof: datamodel: history: true // needed for CUE compat

	Name: string
	Type: string

	// relation type, open to be flexible
	Relation?: {
		Name:  string
		Type:  "has-one" | "has-many" | "belongs-to" | "many-to-many"
		Other: string // technically a cue path, but as a string
	}

	// what about {val, *val, []val, []*val}
	// we probably don't care about pointer here
	//   that is a language detail (code gen target)

	// we can enrich this for various types
	// in our app or other reusable datamodels
}
