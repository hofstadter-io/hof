package schema

import "strings"

#Datamodel: {
	Name: string

	Modelsets?: #Modelsets
	Models?: #Models
	Views?: #Views
	...
}

#Modelsets: [name=string]: #Modelset & { Name: name, ... }
#Modelset: {
  Name: string
	modelsetName: strings.ToCamel(Name)
	ModelsetName: strings.ToTitle(Name)

	Models?: #Models
	Views?: #Views
}

#Models: [name=string]: #Model & { Name: name, ... }
#Model: {
  Name: string
	modelName: strings.ToCamel(Name)
	ModelName: strings.ToTitle(Name)

	Views?: #Views

	Fields: #Fields
  ...
}

#Fields: [name=string]: #Field & { Name: name, ... }
#Field: {
  Name: string
	fieldName: string | *strings.ToCamel(Name)
	FieldName: string | *strings.ToTitle(Name)

	private: bool | *false

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
