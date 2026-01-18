@experiment(aliasv2)
package veg

import (

	"github.com/hofstadter-io/hof/schemas/env"
)

let root = self

docs: {

	src: env.#Dir & {
		@env(docs-src)
		sources: [root.src.code]
		include: [
			"package.json",
			"pnpm-lock.yaml",
			"pnpm-workspace.yaml",
			"docs/",
		]
	}
	watched: env.#HostDir & {
		path:      flags.disk
		gitignore: true
		include: ["docs"]
		trimPrefix: "docs"
	}

	_mounts: [
		env.Mount & {path: "/work/node_modules", source: env.#Cache & { name: "veg-docs-root-node-modules"}},
		env.Mount & {path: "/work/docs", source: env.#Cache & {
			name: "veg-docs-docs-work-dir"
			watch: true
			source: watched
			// source: env.#Dir & { sources: [src], include: ["docs"], trimPrefix: "docs"}
		}},
		env.Mount & {path: "/work/docs/node_modules", source: env.#Cache & { name: "veg-docs-docs-node-modules"}},
	]

	base: env.#Container & {
		@env(docs-base)
		from: root.ctr.dev
		steps: [
			env.Mount & {path: "/work/node_modules", source: env.#Cache & { name: "veg-docs-root-node-modules"}},
			env.Mount & {path: "/work/docs/node_modules", source: env.#Cache & { name: "veg-docs-docs-node-modules"}},

			// dumb JS ecosystem...
			env.EnvVars & {CI: "true"},

			// install dependencies
			env.Dir & {source: src, include: ["package.json","pnpm-lock.yaml","pnpm-workspace.yaml","docs/package.json"]},
			env.Sh & {script: "cd docs && pnpm install --frozen-lockfile"},

			// add code and prep for what comes next
			env.Dir & {source: src},
			env.Workdir & { path: "/work/docs"},
		]
	}

	dev: env.#Service & {
		@env(docs-dev)
		hostname: "docs"
		ports: [{port: 3000}]
		source: env.#Container & {
			from: base
			steps: [
				env.Entrypoint & {args: ["pnpm", "run", "dev"]},
				_mounts,
			]
		}
	}

	hack: env.#Container & {
		@env(docs-hack)
		from: base
		steps: [
			env.DefaultTerm & {args: ["zsh"]},
			env.BindService & { service: dev },
			_mounts,
		]
	}

}
