package datamodel

#Modelsets: [ModelsetName=string]: #Modelset & { Name: ModelsetName, ... }
#Modelset: {
	// The name is used to collect the models
  Name: string

	Modelsets: #Modelsets
	Models: #Models
	Views: #Models
}
