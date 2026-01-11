@experiment(aliasv2)
package veg

import (
	"github.com/hofstadter-io/hof/schemas/env"
)

let root = self

extn: {
	vscode: {
		#ver: "0.0.1"
		src: env.#Dir & {
			@env(vscode-src)
			name: "vscode-src"
			sources: [root.src.code]
			include: [
				"package.json",
				"pnpm-lock.yaml",
				"pnpm-workspace.yaml",
				"extensions/vscode",
			]
		}
		build: env.#Container & {
			@env(vscode-build)
			name: "vscode-build"
			from: "\(root.flags.registry)/veg-dev:local"
			steps: [
				// add source
				env.Dir & {path: "/work", source: src},

				// linting 
				// shouldi to make sure package.json is up to date with what we see

				// actual build steps
				env.Sh & {script: "pnpm install"},
				env.Sh & {script: "pnpm vscode:build:prod"},
				env.Sh & {script: "pnpm vscode:package"},
			]
		}
		vsix: env.#File & {
			@env(vscode-vsix)
			trimPrefix: "/work/extensions/vscode/extension/"
			path:       "\(trimPrefix)veg-\(#ver).vsix"
			source:     build
		}
	}
}
