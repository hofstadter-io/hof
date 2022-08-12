package dm

#Common: {
	Name: string
	Labels: [string]: string
}

#Datamodel: {
	@hof(datamodel)
	#Common

	// Models in the data model, ordered
	Models: #Models
	// (turn Ordered* into a default calculation, so user can always write their own)
	OrderedModels: [...#Model] | *[ for M in Models {M}]

	// Custom views not tied to a specific model
	Views?: #Views

	// History of the data models
	// Will be populated by hof
	History?: [version=string]: #Datamodel
	...
}

#Models: [name=string]: #Model & {Name: name, ...}
#Model: {
	@hof(model)
	#Common

	Fields: #Fields

	// this might all best be filtered to prevent cycles
	Relations: #Relations
	Views?:     #Views

	// TODO, can we calc this in hof and maintain output order stability?
	// (turn Ordered* into a default calculation, so user can always write their own)
	OrderedFields: [...#Field] | *[ for F in Fields {F}]
	OrderedRelations: [...#Relation] | *[ for R in Relations {R}]

	...
}

#Fields: [name=string]: #Field & {Name: name, ...}
#Field: {
	@hof(field)
	#Common

	// this should be a string you can use within your templates
	Type: string

	...
}

#Views: [name=string]: #View & {Name: name, ...}
#View: {
	@hof(view)
	#Common

	// this might all best be filtered to prevent cycles
	Models: #Models
	Fields: #Fields
	Relations: #Relations

	// Todo, make these calculated
	// (turn Ordered* into a default calculation, so user can always write their own)
	OrderedModels:    [...#Model] | *[ for M in Models {M}]
	OrderedFields:    [...#Field] | *[ for F in Fields {F}]
	OrderedRelations: [...#Relation] | *[ for R in Relations {R}]

	...
}

#Relations: [name=string]: #Relation & {Name: name, ...}
#Relation: {
	@hof(reln)
	#Common

	// The relation to another thing

	// This is the relation type, open for debate on what this could or should be
	Reln: "BelongsTo" | "HasOne" | "HasMany" | "ManyToMany" | string

	// This is the other type or side of the relation
	// It is left open so it can be a string or CUE reference
	// 
	// An issue here is cycle detection
	// so we can't use CUE references if we want
	// both sides of the relation to point at eachother for bookkeeping reasons
	//
	// for now, hofmods and examples will assume this is a CUE path to lookup
	Type: string | _

	...
}

// TODO, helper to extract basic fields into a new value
// so we can have better / more info to the other side
// while still preventing cycle errors
#MakeReln: { ... }

