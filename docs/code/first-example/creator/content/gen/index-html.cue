// Generator definition
#Generator: gen.#Generator & {
	OnceFiles: [...gen.#File] & [
			{
			TemplatePath: "index.html"
			Filepath:     "\(Outdir)/client/index.html"
		},
	]
}
