@experiment(aliasv2)
package veg

import (
	"github.com/hofstadter-io/hof/lib/env/common/bases"
	"github.com/hofstadter-io/hof/lib/env/common/steps/lang"
	"github.com/hofstadter-io/hof/lib/env/common/steps/tool"
	"github.com/hofstadter-io/hof/lib/env/common/steps/util"
	"github.com/hofstadter-io/hof/schemas/env"
)

ctr~C: {
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
			util.apt.install & {#pkgs: [
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
			hof.cli,
			tool.github.cli,

			// tools for agents
			tool.agents.lsp2mcp,
		]
	}

	ops: env.#Container & {
		@env()
		from: bases.debian

		steps: [
			hof.cli,
			lang.cue.default,
			tool.k8s.kubectl,
			tool.k8s.helm,
			tool.k8s.crane,
			util.apt.install & {#pkgs: ["ansible"]},
			tool.hashicorp.terraform,
			tool.hashicorp.packer,
		]
	}
	_clis: {
		gcp: tool.cloud.gcloud
		aws: tool.cloud.awscli
		az:  tool.cloud.azure
	}
	for c, cli in _clis {
		"ops-\(c)": env.#Container & {@env(), from: ops, steps: [cli]}
	}
	"ops-all": env.#Container & {@env(), from: ops, steps: [for _, cli in _clis {cli}]}

	incept: env.#Container & {
		@env()
		from: C["ops-all"]

		steps: [
			tool.dagger.cli,
			tool.docker.cli,
		]
	}

	vegeta: env.#Container & {
		@env()
		#hof: metadata: description: "all of the veggie images, it's over 9000"
		from: ops

		steps: incept.steps
	}
}