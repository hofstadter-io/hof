@experiment(aliasv2)
package devex

import (
	"github.com/hofstadter-io/hof/lib/env/devex/steps/lang"
	"github.com/hofstadter-io/hof/lib/env/devex/steps/tool"
	"github.com/hofstadter-io/hof/lib/env/devex/steps/utils"
	"github.com/hofstadter-io/hof/schemas/env"
)

veg: {
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
    labels: env.DefaultLabels & { #name: name }
  }

	base: env.#Container & {
		@env()
		#hof: metadata: description: "A minimal debian image with a few common tools"
		from: "debian:13-slim"

		steps: [
			// default workdir (for wide default consistency)
			env.Workdir & {path: "/work"},

			// shared apt caches, for all derived images as well
			// ya'know, instead of cleaning and refetching all the time?
			utils.apt.mounts.varLib,
			utils.apt.mounts.varCache,
			// need to update once at the beginning
			utils.apt.update,

			// basics
			utils.apt.install & {#pkgs: [
				"ca-certificates",
				"curl",
				"git",
				"gnupg",
				"make",
				"unzip",
				"wget",
				"xz-utils",
				"zsh",
			]},

			// term customization
			tool.zsh.customize,
		]
	}

	dev: env.#Container & {
		@env()
		#hof: metadata: description: "A development image with many tools"

		from: base

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
		from: dev

		steps: [
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
