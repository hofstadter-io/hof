package env

import (
	"github.com/hofstadter-io/hof/schemas"
)

// todo, registry auth

// Definition for a container env
#Container: {
	// always configure a veg node around containers
	schemas.Hof
	#hof: env: {
		root: true // need to figure out what this really means, how it interacts with discovery & cli vs walking a CUE value to construct a giant dagger dag
		kind: "container"
	}

	// duplicative, but needed for the decoding switch statement simplicity (it uses $kind, while `veg node` uses #hof.kind)
	$kind: "#container"

	// the name of the container or environment
	name: string

	// need some kind of from for host / git / oci
	from: string | #Container | #HostImage

	// you can do this in steps, but it might be nice to have it
	// 1. extracted / separate for easy usage in k8s (i.e.)
	// 2. auto add them at the beginning before any steps
	// this seems a reasonable DX
	envs: [string]: string

	// steps to build an image or environment
	steps: [...]

	// labels are applied at the end
	labels: [string]: string
	// defaults (standard to have these three)
	labels: {
		"org.opencontainers.image.title": string | *name
		"org.opencontainers.image.version": string | *"latest"
		"org.opencontainers.image.commit":  string | *"dirty"
		// TODO, support some well-known keywords for git
	}
}