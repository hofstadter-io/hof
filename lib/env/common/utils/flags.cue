package utils

defaultFlags: {
	// source code origins
	local: string @tag(local)
	fork:  string @tag(fork)
	repo:  string @tag(repo)

	// control which origin
	use: *"local" | "fork" | "repo" @tag(use,short=local|fork|repo)

	// operation mode
	defaultModes: "lite" | "full" | "ci" | "canary" | "prod"
	mode:         string | *defaultModes @tag(mode)

	// git overrides
	branch: string | *"main" @tag(branch)
	target: string | *"main" @tag(target)
	gitref: string | *"main" @tag(gitref)
}
