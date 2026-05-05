// Generator definition
#Generator: gen.#Generator & {
	// Define the resource Files
	ResourceFiles: [...gen.#File] & [
			// REST handlers, as before
			// HTML content
			for _, R in In.Resources {
			In: {
				RESOURCE: R
			}
			TemplatePath: "resource.html"
			Filepath:     "\(Outdir)/client/\(strings.ToLower(In.RESOURCE.Name))"
		},
		// HTML content
		for _, R in In.Resources {
			In: {
				RESOURCE: R
			}
			TemplatePath: "resource.js"
			Filepath:     "\(Outdir)/client/\(strings.ToLower(In.RESOURCE.Name)).js"
		},
	]
}
