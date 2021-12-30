package dm

#Datamodel: {
	@datamodel()
	Name: string

	// Models in the data model, ordered
	Models: #Models

	// Custom views not tied to a specific model
	Views?: #Views

	// History of the data models
	// Will be populated by hof
	History?: [version=string]: #Datamodel
	...
}

#Models: [name=string]: #Model & { Name: name, ... }
#Model: {
	@dm_model()
  Name: string

	Fields: #Fields

	Relations?: #Relations
	Views?:     #Views

	...
}

#Fields: [name=string]: #Field & { Name: name, ... }
#Field: {
	@dm_field()
  Name: string

	Type: string

  ...
}

#Relations: [name=string]: #Relation & { Name: name, ... }
#Relation: {
	@dm_relation()
  Name: string

	// Relation to another thing
		
	// This is the relation type, we probably want this more open
	// Relation: "BelongsTo" | "HasOne" | "HasMany" | "Many2Many"
	Relation: _

	// This is the other type or side of the relation
	// It is left open so it can be a string or CUE reference
	// 
	// An issue here is cycle detection
	// so we can't use CUE references if we want
	// both sides of the relation to point at eachother for bookkeeping reasons
	Type: _

  ...
}

#Views: [name=string]: #View & { Name: name, ... }
#View: {
	@dm_view()
  Name: string

	Models: #Models
	Fields: #Fields

  ...
}
