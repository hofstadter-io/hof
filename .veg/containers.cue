@experiment(aliasv2)
package veg

import (
	"github.com/hofstadter-io/hof/catalogs/env/bases"
	"github.com/hofstadter-io/hof/catalogs/env/packs"
	"github.com/hofstadter-io/hof/catalogs/env/utils"
	"github.com/hofstadter-io/hof/schemas/env"
)

_packs: packs & {
	flags: {
		goos: "linux"
		arch: flags.arch
	}
}

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
		from: bases.debian13.minimal
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

		// cmd.test.tasks.go.steps.0.0.from.from...  steps.5.0.0.args.2
		from: bases.debian13.default

		steps: [
			// customization
			_packs.tool.zsh.customize,

			// deps for go/node/python -> c/c++ situations (like CGO)
			utils.apt.install & {#pkgs: ["gcc", "libc6-dev"]},

			// setup languages
			_packs.lang.go.defaultSteps,
			_packs.lang.cue.defaultSteps,
			_packs.lang.python.astralSteps,
			_packs.lang.node.defaultSteps,

			// spell checker
			_packs.lang.node.cspell.install,
			// TODO, sync or use custom dictionaries with vscode and others
			// or more likely have a way to load them and adding them here should be a one liner

			// tools for agents
			_packs.tool.github.cli,
			_packs.tool.agents.lsp2mcp,

			// add a bunch of tools (from packs)
			_packs.containers.docker.cli.install,
			_packs.containers.dagger.cli.install,
			_packs.containers.cosign.cli.install,
			_packs.containers.buildah.cli.install,
			_packs.containers.dive.cli.install,

			// still to be moved to packs
			_packs.tool.hashicorp.packer,
			_packs.tool.hashicorp.terraform,
			_packs.tool.k8s.kubectl,
			_packs.tool.k8s.crane,
			_packs.tool.k8s.helm,
			_packs.tool.k8s.kind.binary,
		]
	}
	run: env.#Container & {
		@env()
		#hof: {
			id: "veg-run"
			metadata: {
				name:        id
				description: "runtime veg-dev, with socket, secrets, and such"
			}
		}
		name: #hof.metadata.name

		from: dev
		steps: [
			// add hof late (not in dev), because it changes frequently
			hof.File.linux,

			// config / env stuff
			// _packs.tool.k8s.kind.config,
			// going to switch to k3d / k3s

			// add the socket for inception
			env.UnixSocket & {path: "/var/run/docker.sock", source: host.docker.socket},

			// // bind lsp servers, started on demand
			// env.BindService & {service: lang.go.lsp},
			// env.BindService & {service: lang.cue.lsp},
			// env.BindService & {service: lang.node.lsp},
			// env.BindService & {service: lang.python.lsp},

			env.Dir & {path: "/root/.ssh", source: secrets.dotssh},
			env.Dir & {path: "/root/.kube", source: secrets.kubecfg},
			env.SecretVars & {
				GOOGLE_API_KEY: secrets.google
			},
		]
	}

	// set id for all ops-, used for caching in env, and default names based on that
	[=~"ops"]~(k,_): {@env()
		#hof: {id: "veg-\(k)", metadata: {name: string | *id}}
		name: string | *#hof.metadata.name
	}

	// base ops container
	"ops": env.#Container & {
		from: bases.debian13.default
		steps: [
			hof.File.linux,
			_packs.tool.hashicorp.terraform,
			_packs.tool.hashicorp.packer,
			_packs.tool.k8s.kubectl,
			_packs.tool.k8s.helm,
			_packs.tool.k8s.crane,
			_packs.tool.github.cli,
		]
	}

	// create branches from the ops bases for each cloud cli
	for c, cli in _clis {
		"ops-\(c)": env.#Container & {@env(), from: root.ctr.ops, steps: [cli]}
	}
	"ops-all": env.#Container & {from: root.ctr["ops"], steps: [for _, cli in _clis {cli}]}
	_clis: {
		gcp: _packs.tool.cloud.gcloud
		aws: _packs.tool.cloud.awscli
		az:  _packs.tool.cloud.azure
		ansible: utils.apt.install & {#pkgs: ["ansible"]}
	}

}

fmtr: {
	[string]~(f,_): [string]~(k,_): {@env(), name: "fmtr-\(f)-\(k)"}
	black: {
		src: env.#HostDir & {@env(fmt-black-src), path: "lib/fmt/tools/black"}
		img: env.#Container & {
			@env(fmt-black-img)
			from: bases.debian13.default
			steps: [
				_packs.lang.python.defaultSteps,
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
		src: env.#HostDir & {@env(fmt-prettier-src), path: "lib/fmt/tools/prettier"}
		img: env.#Container & {
			@env(fmt-prettier-img)
			from: bases.debian13.default
			steps: [
				utils.apt.install & {#pkgs: [
					"gcc",
					"libc6-dev",
					"ruby-dev",
				]},
				env.Bash & {script: "gem install bundler haml prettier_print rbs syntax_tree syntax_tree-haml syntax_tree-rbs"},
				_packs.lang.node.defaultSteps,
				env.Dir & {path: "/work", source: src},
				env.Exec & {args: ["yarn", "install", "--ignore-engines"]},
				env.Entrypoint & {args: ["node", "prettier.js"]},
				env.Expose & {port: 3000},
			]
		}
	}
}
