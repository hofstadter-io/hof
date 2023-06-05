package dm

#Datamodel: {
	@datamodel()
	Name: string

	// Models in the data model, ordered
	Models: #Models

	...
}

#Models: [name=string]: #Model & {Name: name, ...}
#Model: {
	@dm_model()
	Name: string

	Fields: #Fields

	Relations?: #Relations
	...
}

#Fields: [name=string]: #Field & {Name: name, ...}
#Field: {
	@dm_field()
	Name: string

	Type: string

	...
}

#Relations: [name=string]: #Relation & {Name: name, ...}
#Relation: {
	@dm_relation()

	Name: string

	// Relation to another thing

	// This is the relation type
	Reln: "BelongsTo" | "HasOne" | "HasMany" | "ManyToMany"

	// This is the other type or side of the relation (a Model)
	Type: string

	...
}
