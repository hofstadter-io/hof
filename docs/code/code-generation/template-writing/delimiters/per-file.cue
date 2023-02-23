MyGen: gen.#Generator & {

	// ...

	Out: [...gen.#File] & [
		// A file with the default {{ }} delimiters
		{
			TemplateContent: "Val.a = '{{ .Val.a }}'\n"
			Filepath:        "default.txt"
		},

		// Alternate delims
		{
			TemplateContent: "Val.a = '{% .Val.a %}'\n"
			Filepath:        "altdelim.txt"

			// Set alternate delimiters to anything you like
			TemplateDelims: {
				LHS: "{%"
				RHS: "%}"
			}
		},
	]
}
