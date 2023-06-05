Generator: gen.Generator & {
	// We make resources from the data model
	// and there is no new inputs for the user
	In: {
		Resources: (schema.DatamodelToResources & {"Datamodel": Datamodel}).Resources
	}

	// Add a new line in _All
	All: [
		for _, F in ResourceFiles {F},
	]

	// Define the resource Files
	ResourceFiles: [...gen.File] & [
			for _, R in In.Resources {
			In: {
				RESOURCE: R
			}
			TemplatePath: "resource.go"
			Filepath:     "\(Outdir)/resources/\(In.RESOURCE.Name).go"
		},
	]

}
