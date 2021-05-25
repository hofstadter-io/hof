package dm

import "strings"

#Datamodel: {
	Name: string

	Modelsets?: #Modelsets
	Models?: #Models
	Views?: #Views

	History?: #Modelsets
	...
}

#Modelsets: [name=string]: #Modelset & { Name: name, ... }
#Modelset: {
  Name: string
	modelsetName: strings.ToCamel(Name)
	ModelsetName: strings.ToTitle(Name)

	Models?: #Models
	Views?: #Views

	...
}

#Models: [name=string]: #Model & { Name: name, ... }
#Model: {
  Name: string
	modelName: strings.ToCamel(Name)
	ModelName: strings.ToTitle(Name)

	ORM: bool | *false
	SoftDelete: bool | *false
	Permissioned: bool | *true

	Views?: #Views

	Fields: #Fields

	Relations: #Relations

	...
}

#Fields: [name=string]: #Field & { Name: name, ... }
#Field: {
  Name: string
	fieldName: string | *strings.ToCamel(Name)
	FieldName: string | *strings.ToTitle(Name)

	type: string

	validation?: [string]: _

	private: bool | *false

  ...
}

#Relations: [name=string]: #Relation & { Name: name, ... }
#Relation: {
  Name: string
	relnName: string | *strings.ToCamel(Name)
	RelnName: string | *strings.ToTitle(Name)

	foreignKey?: string
	relation: "BelongsTo" | "HasOne" | "HasMany" | "Many2Many"
	type: string
	table?: string

  ...
}

#Views: [name=string]: #View & { Name: name, ... }
#View: {
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
