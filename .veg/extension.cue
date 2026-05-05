@experiment(aliasv2)
package veg

import (
	"github.com/hofstadter-io/hof/schemas/env"
)

let root = self

extn: {
	#ver: string | *"v0.7.0-alpha.2"
	vscode: {
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
			from: root.ctr.dev
			steps: [
				// add mounts

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
