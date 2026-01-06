@experiment(aliasv2)
package veg

import (
	"github.com/hofstadter-io/hof/catalogs/env/bases"
	"github.com/hofstadter-io/hof/catalogs/env/packs"
	"github.com/hofstadter-io/hof/catalogs/env/steps"
	"github.com/hofstadter-io/hof/catalogs/env/utils"
	"github.com/hofstadter-io/hof/schemas/env"
)

_packs: packs
_steps: steps

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
			_steps.tool.zsh.customize,

			// deps for go/node/python -> c/c++ situations (like CGO)
			utils.apt.install & {#pkgs: ["gcc", "libc6-dev"]},

			// setup languages
			_steps.lang.go.defaultSteps,
			_steps.lang.cue.default,
			_steps.lang.node.default,
			_steps.lang.python.default,
			_steps.lang.python.dev, // depends on node

			// tools for agents
			_steps.tool.github.cli,
			_steps.tool.agents.lsp2mcp,

      // add a bunch of tools (from packs)
      _packs.containers.docker.cli.install,
      _packs.containers.dagger.cli.install,
      _packs.containers.cosign.cli.install,
      _packs.containers.buildah.cli.install,
      _packs.containers.dive.cli.install,

			// still to be moved to packs
			_steps.tool.hashicorp.packer,
			_steps.tool.hashicorp.terraform,
			_steps.tool.k8s.kubectl,
			_steps.tool.k8s.crane,
			_steps.tool.k8s.helm,
			_steps.tool.k8s.kind.binary,
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
			// config / env stuff
			_steps.tool.k8s.kind.config,

      // add the socket for inception
      env.UnixSocket & { path: "/var/run/docker.sock", source: host.docker.socket },

			// // bind lsp servers, started on demand
			// env.BindService & {service: lang.go.lsp},
			// env.BindService & {service: lang.cue.lsp},
			// env.BindService & {service: lang.node.lsp},
			// env.BindService & {service: lang.python.lsp},

			// add hof late, because it changes frequently
			hof.File.linux,
			env.Dir & {path: "/work", source: src.code},
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
			_steps.tool.hashicorp.terraform,
			_steps.tool.hashicorp.packer,
			_steps.tool.k8s.kubectl,
			_steps.tool.k8s.helm,
			_steps.tool.k8s.crane,
			_steps.tool.github.cli,
		]
	}

	// create branches from the ops bases for each cloud cli
	for c, cli in _clis {
		"ops-\(c)": env.#Container & {@env(), from: root.ctr.ops, steps: [cli]}
	}
	"ops-all": env.#Container & {from: root.ctr["ops"], steps: [for _, cli in _clis {cli}]}
	_clis: {
		gcp: _steps.tool.cloud.gcloud
		aws: _steps.tool.cloud.awscli
		az:  _steps.tool.cloud.azure
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
				_steps.lang.python.default,
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
			@env(fmt-pretteir-img)
			from: bases.debian13.default
			steps: [
				utils.apt.install & {#pkgs: [
					"gcc",
					"libc6-dev",
					"ruby-dev",
				]},
				env.Bash & {script: "gem install bundler haml prettier_print rbs syntax_tree syntax_tree-haml syntax_tree-rbs"},
				_steps.lang.node.install,
				env.Dir & {path: "/work", source: src},
				env.Exec & {args: ["yarn", "install", "--ignore-engines"]},
				env.Entrypoint & {args: ["node", "prettier.js"]},
				env.Expose & {port: 3000},
			]
		}
	}
}
