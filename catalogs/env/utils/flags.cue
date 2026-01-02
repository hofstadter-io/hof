package utils

defaultFlags: {
	// source code origins
	local: string @tag(local)
	repo:  string @tag(repo)
	fork:  string @tag(fork)

	// control which origin
	use: *"local" | "repo" | "fork" @tag(use,short=local|repo|fork)

	// operation mode
	defaultModes: [...string] | *["lite" | "full" | "ci" | "canary" | "prod"]
	mode: string | or(defaultModes) @tag(mode)

	// git overrides
	branch: string | *"main" @tag(branch)
	target: string | *"main" @tag(target)
	gitRef: string | *"main" @tag(gitRef)
}
