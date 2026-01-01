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
		from: bases.debian.minimal
		steps: [hof.cli]
	}
	dev: env.#Container & {
		@env()
		#hof: metadata: {
			id:          "veg-dev"
			name:        id
			description: "setup needed to work on veg"
		}
		name: #hof.metadata.name

		from: bases.debian.default

		steps: [
			// customization
			tool.zsh.customize,

			// deps for go/node/python -> c/c++ situations (like CGO)
			util.apt.install & {#pkgs: [ "gcc", "libc6-dev" ]},

			// binary tools
			hof.cli,
			tool.github.cli,

			// setup languages
			lang.go.default,
			// lang.cue.default,
			lang.node.default,
			lang.python.default,
			lang.python.dev, // depends on node

			// devops stuff

			// tools for agents
			tool.agents.lsp2mcp,
		]
	}

	// set id for all ops-, used for caching in env, and default names based on that
	[=~"ops-"]~(k,_): { @env()
		#hof: metadata: { id: "veg-\(k)", name: string | *id }
		name: string | *#hof.metadata.name
	}

	// base ops container
	"ops": env.#Container & {
		from: bases.debian.minimal
		steps: [
			hof.cli,
			tool.hashicorp.terraform,
			tool.hashicorp.packer,
			tool.k8s.kubectl,
			tool.k8s.helm,
			tool.k8s.crane,
			tool.github.cli,
		]
	}

	// create branches from the ops bases for each cloud cli
	for c, cli in _clis {
		"ops-\(c)": env.#Container & {@env(), from: root.ctr.ops, steps: [cli]}
	}
	"ops-all": env.#Container & {from: root.ctr["ops"], steps: [for _, cli in _clis {cli}]}
	_clis: {
		gcp: tool.cloud.gcloud
		aws: tool.cloud.awscli
		az:  tool.cloud.azure
		ansible: util.apt.install & {#pkgs: ["ansible"]}
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
