package dm

import "strings"

#Datamodel: {
	@datamodel()
	Name: string

	// Models in the data model, ordered
	Models: [...#Model]

	// Custom views not tied to a specific model
	Views: [...#View]

	// History of the data models
	// Should be populated by hof from .hof/dm
	History?: [version=string]: #Models
	...
}

#Models: [name=string]: #Model & { Name: name, ... }
#Model: {
	@dm_model()
  Name: string
	modelName: strings.ToCamel(Name)
	ModelName: strings.ToTitle(Name)

	SoftDelete: bool | *false
	Permissioned: bool | *true

	Views?: #Views

	Fields: #Fields

	Relations: #Relations

	...
}

#Fields: [name=string]: #Field & { Name: name, ... }
#Field: {
	@dm_field()
  Name: string
	fieldName: string | *strings.ToCamel(Name)
	FieldName: string | *strings.ToTitle(Name)

	Type: string

	Validation?: [string]: _

	Private: bool | *false

  ...
}

#Relations: [name=string]: #Relation & { Name: name, ... }
#Relation: {
	@dm_relation()
  Name: string
	relnName: string | *strings.ToCamel(Name)
	RelnName: string | *strings.ToTitle(Name)

	ForeignKey?: string
	Relation: "BelongsTo" | "HasOne" | "HasMany" | "Many2Many"
	Type: string
	Table?: string

  ...
}

#Views: [name=string]: #View & { Name: name, ... }
#View: {
	@dm_view()
  Name: string
	viewName: string | *strings.ToCamel(Name)
	ViewName: string | *strings.ToTitle(Name)

	Models: #Models

	Fields: #Fields
  ...
}

#CommonFields: {
	ID:        #UUID
	CID:       #CUID
	CreatedAt: #Datetime
	UpdatedAt: #Datetime
	DeletedAt: #Datetime
}
