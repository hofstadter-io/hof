MyGen: gen.#Generator & {

	// ...

	Templates: [{
		// Loaded templates with the default {{ }} delimiters
		Globs: ["templates/*"]
		TrimPrefix: "templates/"
	}, {
		// Loaded templates with the default alternate delimiters
		Globs: ["templates/alt/*"]
		TrimPrefix: "templates/"
		Delims: {
			LHS: "{%"
			RHS: "%}"
		}
	}]
}

