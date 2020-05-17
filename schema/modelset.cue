package schema

#Modelsets: [ModelsetName=string]: #Modelset & { Name: ModelsetName, ... }
#Modelset: {
	// The name is used to collect the models
  Name: string
	Entrypoint: string
	Workbase: string | *""

	Stores: [Key=string]: string
}
