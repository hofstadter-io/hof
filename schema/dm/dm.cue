package dm

#Datamodel: {
	@datamodel()
	Name: string

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
	@dm_model()
	Name: string

	Fields: #Fields
	// (turn Ordered* into a default calculation, so user can always write their own)
	OrderedFields: [...#Field] | *[ for F in Fields {F}]

	Relations?: #Relations
	Views?:     #Views

	...
}

#Fields: [name=string]: #Field & {Name: name, ...}
#Field: {
	@dm_field()
	Name: string

	// this should be a string you can use within your templates
	Type: string

	...
}

#Views: [name=string]: #View & {Name: name, ...}
#View: {
	@dm_view()
	Name: string

	Models: #Models
	// (turn Ordered* into a default calculation, so user can always write their own)
	Fields: #Fields
	OrderedFields: [...#Field] | *[ for F in Fields {F}]

	...
}

#Relations: [name=string]: #Relation & {Name: name, ...}
#Relation: {
	@dm_relation()

	Name: string

	// Relation to another thing

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
