@experiment(aliasv2)
package veg

import (
	"github.com/hofstadter-io/hof/lib/env/common/bases"
	"github.com/hofstadter-io/hof/lib/env/common/steps/lang"
	"github.com/hofstadter-io/hof/lib/env/common/steps/tool"
	"github.com/hofstadter-io/hof/lib/env/common/steps/util"
	"github.com/hofstadter-io/hof/schemas/env"
)

let root = self

ctr: {
	min: env.#Container & {
		@env()
		#hof: metadata: {
			id:          "veg-min"
			name:        id
			description: "minimal veg, eat your veggies!"
		}
		name: #hof.metadata.name

		from: bases.debian

		steps: [
			hof.cli,
		]
	}
	dev: env.#Container & {
		@env()
		#hof: metadata: {
			id:          "veg-dev"
			name:        id
			description: "setup needed to work on veg"
		}
		name: #hof.metadata.name

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

			// binary tools
			hof.cli,
			tool.github.cli,

			// setup languages
			lang.go.default,
			lang.cue.default,
			lang.node.default,
			lang.python.default,
			lang.python.dev, // depends on node

			// tools for agents
			tool.agents.lsp2mcp,
		]
	}

	vegeta: env.#Container & {
		@env()
		#hof: metadata: {
			id:          "vegeta"
			name:        id
			description: "all of the veggie dev, it's over 9000"
		}
		name: #hof.metadata.name

		from: root.ctr.dev

		steps: [
			tool.k8s.kubectl,
			tool.k8s.helm,
			tool.k8s.crane,
			tool.github.cli,
			tool.cloud.gcloud,
			tool.dagger.cli,
			tool.docker.cli,
			tool.hashicorp.terraform,
			tool.hashicorp.packer,
			util.apt.install & {#pkgs: ["ansible"]},
		]
	}

	// set id for all ops-, used for caching in env, and default names based on that
	[=~"ops-"]~(k,_): { @env()
		#hof: metadata: { id: "veg-\(k)", name: string | *id }
		name: string | *#hof.metadata.name
	}
	// sugar image, override the name, keep id for caching
	"ops": { name: "veg-ops", root.ctr["ops-lite"] }

	// base ops container
	"ops-lite": env.#Container & {
		from: bases.debian
		steps: [
			hof.cli,
			lang.cue.default,
			tool.k8s.kubectl,
			tool.k8s.helm,
			tool.k8s.crane,
			tool.github.cli,
		]
	}

	// full ops container
	"ops-full": env.#Container & {
		from: root.ctr["ops-lite"]
		steps: [
			util.apt.install & {#pkgs: ["ansible"]},
			tool.hashicorp.terraform,
			tool.hashicorp.packer,
			tool.dagger.cli,
			tool.docker.cli,
		]
	}

	// create branches from the ops bases for each cloud cli
	for c, cli in _clis {
		"ops-lite-\(c)": env.#Container & {@env(), from: root.ctr["ops-lite"], steps: [cli]}
		"ops-full-\(c)": env.#Container & {@env(), from: root.ctr["ops-full"], steps: [cli]}
	}
	"ops-lite-all": env.#Container & {from: root.ctr["ops-lite"], steps: [for _, cli in _clis {cli}]}
	"ops-full-all": env.#Container & {from: root.ctr["ops-full"], steps: [for _, cli in _clis {cli}]}
	_clis: {
		gcp: tool.cloud.gcloud
		aws: tool.cloud.awscli
		az:  tool.cloud.azure
	}

}

fmtr: {
	[string]~(f,_): [string]~(k,_): {@env(), name: "fmtr-\(f)-\(k)"}
	black: {
		src: env.#HostDir & {path: "lib/fmt/tools/black"}
		img: env.#Container & {
			from: bases.debian
			steps: [
				lang.python.default,
				env.Dir & {path: "/work", source: src},
				env.Bash & {
					script: """
						pipenv --python /usr/bin/python3
						pipenv install
						"""
				},
				env.Entrypoint & {args: ["gunicorn", "app:app", "--bind", "0.0.0.0:3000", "--log-file", "-"]},
				env.Expose & {port: 3000},
			]
		}
	}
	prettier: {
		src: env.#HostDir & {path: "lib/fmt/tools/prettier"}
		img: env.#Container & {
			from: bases.debian
			steps: [
				util.apt.install & {#pkgs: [
					"gcc",
					"libc6-dev",
					"ruby-dev",
				]},
				env.Bash & {script: "gem install bundler haml prettier_print rbs syntax_tree syntax_tree-haml syntax_tree-rbs"},
				lang.node.install,
				env.Dir & {path: "/work", source: src},
				env.Exec & {args: ["yarn", "install", "--ignore-engines"]},
				env.Entrypoint & {args: ["node", "prettier.js"]},
				env.Expose & {port: 3000},
			]
		}
	}
}
