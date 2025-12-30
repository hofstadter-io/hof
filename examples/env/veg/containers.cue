@experiment(aliasv2)
package veg

import (
	"github.com/hofstadter-io/hof/lib/env/common/bases"
	"github.com/hofstadter-io/hof/lib/env/common/steps/lang"
	"github.com/hofstadter-io/hof/lib/env/common/steps/tool"
	"github.com/hofstadter-io/hof/lib/env/common/steps/utils"
	"github.com/hofstadter-io/hof/schemas/env"
)

ctr: {
	// apply these to all fields, except vegeta
	[string]~(key,_): {
		// todo, if we have more than containers in this struct, we can if $kind == "#container" { ... }
		name: string
		if key == "vegeta" {
			name: key
		}
		if key != "vegeta" {
			name: "veg-\(key)"
		}
		labels: env.DefaultLabels & {#name: name}
	}

	dev: env.#Container & {
		@env()
		#hof: metadata: description: "A development image with many tools"

		from: bases.debian

		steps: [
			// todo, put these with the tools that depend on them (if they are one)
			utils.apt.install & {#pkgs: [
				// deps for go/node/python -> c/c++ situations (like CGO)
				"g++",
				"gcc",
				"libc6-dev",
				"netbase",
				"pkg-config",
				"sq",
			]},

			// setup languages
			lang.go.default,
			lang.cue.default,
			lang.node.default,
			lang.python.default,

			// other binary tools

			// tools for agents
			tool.agents.lsp2mcp,
		]
	}

	ops: env.#Container & {
		@env()
		#hof: metadata: description: "Extension to veg-dev to add devops tooling"
		from: bases.debian

		steps: [
			// hof, tbd
			lang.cue.default,
			tool.k8s.kubectl,
			tool.k8s.helm,
			tool.k8s.crane,
			tool.hashicorp.terraform,
			tool.hashicorp.packer,
		]
	}

	incept: env.#Container & {
		@env()
		#hof: metadata: description: "Extension to veg-dev to add docker/dagger setup for inception"
		from: dev

		steps: [
			tool.dagger.cli,
			tool.docker.cli,
		]

		// whatever we import / user here, should also have mounts defined for easy reuse for runtime (run/up/asService)
	}

	vegeta: env.#Container & {
		@env()
		#hof: metadata: description: "all of the veggie images, it's over 9000"
		from: ops

		// TODO, realize the dagger way of diamond build pattern optimizations (once we branch more than 2 wide, and have _tools in a better place with #file/#dir)
		steps: incept.steps
	}

}
