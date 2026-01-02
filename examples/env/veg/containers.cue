@experiment(aliasv2)
package veg

import (
	"github.com/hofstadter-io/hof/catalogs/env/bases"
	// "github.com/hofstadter-io/hof/catalogs/env/packs"
	isteps "github.com/hofstadter-io/hof/catalogs/env/steps"
	"github.com/hofstadter-io/hof/catalogs/env/utils"
	"github.com/hofstadter-io/hof/schemas/env"
)

let root = self

ctr: {
	min: env.#Container & {
		@env()
		#hof: {
			id: "veg-min"
			metadata: {
				name:        id
				description: "minimal veg, eat your veggies!"
			}
		}
		name: #hof.metadata.name
		from: bases.debian.minimal
		steps: [hof.File.linux]
	}
	dev: env.#Container & {
		@env()
		#hof: {
			id: "veg-dev"
			metadata: {
				name:        id
				description: "setup needed to work on veg"
			}
		}
		name: #hof.metadata.name

		from: bases.debian.default

		steps: [
			// customization
			isteps.tool.zsh.customize,

			// deps for go/node/python -> c/c++ situations (like CGO)
			utils.apt.install & {#pkgs: ["gcc", "libc6-dev"]},

			// binary tools
			hof.File.linux,
			isteps.tool.github.cli,

			// setup languages
			isteps.lang.go.defaultSteps,
			// lang.cue.default,
			isteps.lang.node.default,
			isteps.lang.python.default,
			isteps.lang.python.dev, // depends on node

			// devops stuff
			// tool.hashicorp.terraform,
			// tool.hashicorp.packer,
			// tool.k8s.kubectl,
			// tool.k8s.helm,
			// tool.k8s.crane,

			// // bind lsp servers, started on demand
			// env.BindService & {service: lang.go.lsp},
			// env.BindService & {service: lang.cue.lsp},
			// env.BindService & {service: lang.node.lsp},
			// env.BindService & {service: lang.python.lsp},

			// tools for agents
			isteps.tool.agents.lsp2mcp,
		]
	}

	// set id for all ops-, used for caching in env, and default names based on that
	[=~"ops"]~(k,_): {@env()
		#hof: {id: "veg-\(k)", metadata: {name: string | *id}}
		name: string | *#hof.metadata.name
	}

	// base ops container
	"ops": env.#Container & {
		from: bases.debian.default
		steps: [
			hof.File.linux,
			isteps.tool.hashicorp.terraform,
			isteps.tool.hashicorp.packer,
			isteps.tool.k8s.kubectl,
			isteps.tool.k8s.helm,
			isteps.tool.k8s.crane,
			isteps.tool.github.cli,
		]
	}

	// create branches from the ops bases for each cloud cli
	for c, cli in _clis {
		"ops-\(c)": env.#Container & {@env(), from: root.ctr.ops, steps: [cli]}
	}
	"ops-all": env.#Container & {from: root.ctr["ops"], steps: [for _, cli in _clis {cli}]}
	_clis: {
		gcp: isteps.tool.cloud.gcloud
		aws: isteps.tool.cloud.awscli
		az:  isteps.tool.cloud.azure
		ansible: utils.apt.install & {#pkgs: ["ansible"]}
	}

}

fmtr: {
	[string]~(f,_): [string]~(k,_): {@env(), name: "fmtr-\(f)-\(k)"}
	black: {
		src: env.#HostDir & {path: "lib/fmt/tools/black"}
		img: env.#Container & {
			from: bases.debian.default
			steps: [
				isteps.lang.python.default,
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
			from: bases.debian.default
			steps: [
				utils.apt.install & {#pkgs: [
					"gcc",
					"libc6-dev",
					"ruby-dev",
				]},
				env.Bash & {script: "gem install bundler haml prettier_print rbs syntax_tree syntax_tree-haml syntax_tree-rbs"},
				isteps.lang.node.install,
				env.Dir & {path: "/work", source: src},
				env.Exec & {args: ["yarn", "install", "--ignore-engines"]},
				env.Entrypoint & {args: ["node", "prettier.js"]},
				env.Expose & {port: 3000},
			]
		}
	}
}
