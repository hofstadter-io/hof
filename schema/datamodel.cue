package schema

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

	Models?: #Models
	Views?: #Models
}

#Models: [name=string]: #Model & { Name: name, ... }
#Model: {
  Name: string

	Views?: #Models

	Fields: [string]: {...}

  ...
}

#Views: [name=string]: #View & { Name: name, ... }
#View: {
  Name: string
	Models: #Models

	Fields: [string]: _
  ...
}

