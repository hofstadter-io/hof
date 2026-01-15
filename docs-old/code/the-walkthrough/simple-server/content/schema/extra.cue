Server: {
	// ...

	// list of file globs to be embedded into the server when built
	EmbedGlobs: [...string]

	// enable prometheus metrics
	Prometheus: bool

	// auth settings (optional)
	Auth?: {
		apikey: bool | *false
		basic:  bool | *false
	}
}
