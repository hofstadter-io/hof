package example

// here we are applying a schema to our input data
// note how the label is the same in both files
Input: #Input

// This is your input schema
#Input: {
	// this is a CUE pattern to apply #Type to every key
	[key=string]: #Type & { 
		// here we are enriching the input, mapping key -> Name
		Name: key
	}
}

// Schema for a Type
#Type: {
	// Name to use in target languages
	Name: string

	// This is a CUE pattern for a struct of structs
	// you can set nested fields based on the key name in [key=string]
	Fields: [field=string]: #Field & { Name: field }

	// Enum of relation types with the key being the Name of the other side
	Relations: [other=string]: "BelongsTo" | "HasOne" | "HasMany" | "ManyToMany"
}

// Schema for a Field
#Field: {
	// Name to use in target languages
	Name: string
	// the type as a string, for flexibility
	Type: string
}

// note, if you remove "Input" from all of the files
// this will have the same effect, but less control at the top-level
// Top-level schema are helpful when you don't control the input data format
// [key=string]: #Type & { Name: key }
