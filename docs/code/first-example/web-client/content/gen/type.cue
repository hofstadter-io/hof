#Generator: gen.#Generator & {

	// Exposed to the user
	Datamodel: schema.#Datamodel

	// Added to the template input
	In: {
		DM: Datamodel
	}

	// include the type files for rendering
	All: [
		for _, F in TypeFiles {F},
	]

	// Define the files generated from our models
	TypeFiles: [...gen.#File] & [
			for _, M in Datamodel.Models {
			In: {
				TYPE: {
					// embed the model fields
					M

					// Extend to include the fields in CUE order with a list
					// This is needed because we want a deterministic order
					// For example, when defining database columns or caluclating a diff
					// We don't want to sort, rather we want to maintain the order the user specifies
					// While CUE has consistent ordering, the internal Go maps do not
					// the mapping from CUE -> Go -> template rendering can misorder
					OrderedFields: [ for _, F in M.Fields {F}]
				}
			}
			TemplatePath: "type.go"
			// We name each file the same as the Model name and 
			Filepath: "\(Outdir)/types/\(In.TYPE.Name).go"
		},
	]
}
