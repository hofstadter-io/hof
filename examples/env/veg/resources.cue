@experiment(aliasv2)
package veg

import (
	"github.com/hofstadter-io/hof/schemas/env"
)

src: {
	repo: env.#GitRepo & {
		@env()
		name: "repo"
		url:  _flags.repo
	}
	code: env.#HostDir & {
		@env()
		path: _flags.code
	}
}

incept: {
	[string]~(K,_): {name: K}
	docker: env.#HostSocket & {@env(), path: "unix:///Users/tony/.colima/default/docker.sock"}
	dagger: env.#HostSocket & {@env(), path: ""}
}

out: {
	// this should filter from either repo or host
	cuemod: env.#Dir & {
		@env()
		path:   "."
		source: src.repo
		include: [
			"cue.mod",
			"schemas",
			"examples",
			"lib/env/common",
			// way more to come here
		]
	}

	// go builds, cross-arch/os
	cli: {}

	vscode: {}

	docs: {}

	fmtrs: {}
}
