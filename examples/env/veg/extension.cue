package veg

import (
	"github.com/hofstadter-io/hof/schemas/env"
)

extn: {
	vscode: {
		src: env.#HostDir & {
			name: "vscode-src"
			path: "."
			include: [
				"package.json",
				"pnpm-lock.yaml",
				"pnpm-workspace.yaml",
				"extensions/vscode",
			]
		}
		build: env.#Container & {
			@env(), @id(vscode-build)
			name: "vscode-build"
			from: "\(flags.registry)/veg-dev:local"
			steps: [
				env.Dir & {path: "/work", source: src},
				env.Bash & {script: "pnpm install"},
				env.Bash & {script: "pnpm build:extn:vscode"},
			]
		}
	}
}

// or split values over files
dist: {
	// vscode: env.#ExportDir & {
	// 	path: "dist/vscode"
	// 	wipe: true
	// 	sources: []
	// }
}
