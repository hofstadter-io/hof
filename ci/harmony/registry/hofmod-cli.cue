package registry

Registry: "hofmod-cli": {
	// set some defaults for all cases, git ref from the label
	[r=string]: {ref: r, url: "https://github.com/hofstadter-io/hofmod-cli"}

	// run the same commands for all cases
	[string]: scripts: [
		"""
			hof mod tidy
			cue eval -a
			""",
	]

	"v0.8.6": {
		type: "cli"
	}

	"branch/master": {
		type: "cli"
	}

}
