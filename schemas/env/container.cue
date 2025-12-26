package env

import (
	"github.com/hofstadter-io/hof/schemas"
)

// Definition for a container env
Container: {
	schemas.Hof
	#hof: env: {
		root: true
		kind: "container"
	}

	name: string

	// need some kind of from for host / git / oci
	from: string | Container

	steps: [...]
	labels: [string]: string

	labels: {
		"org.opencontainers.image.title":   string | *name

    // TODO, support some well-known keywords for git
		"org.opencontainers.image.version": string | *"latest"
		"org.opencontainers.image.commit":  string | *"dirty"
	}
}
