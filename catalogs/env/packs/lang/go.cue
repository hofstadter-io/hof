package lang

import (
	"github.com/hofstadter-io/hof/catalogs/env/bases"
	"github.com/hofstadter-io/hof/catalogs/env/utils"
	slang "github.com/hofstadter-io/hof/catalogs/env/steps/lang"
	"github.com/hofstadter-io/hof/schemas/env"
)

go: {
	// A base container with go tools, linters, and a shared cache.
	ctr: {
		base: env.#Container & {
			from: bases.debian13.minimal
			steps: [
				// deps for cgo and more, from official docs
				utils.apt.install & {#pkgs: [
					"g++",
					"gcc",
					"libc6-dev",
					"netbase",
					"pkg-config",
					"sq",
				]},
				slang.go.default,
			]
		}

		gopls: env.#Container & {
			from: base
			steps: [
				env.Expose & {port: 4000},
			]
		}

		// A dev container with gopls running and attached
		dev: env.#Container & {
			from: ctr.gopls
			steps: [
				env.BindService & {service: go.svc.lsp},
			]
		}
	}

	// gopls-as-a-service
	svc: {
		lsp: env.#Service & {
			// There is also a built in MCP server!
			#port: int | *0
			name: "gopls"
			ports: [{name: "lsp", port: 4000, frontend: #port}]
			args: ["gopls", "serve", "-port=4000"]
			// source: _
			source: go.ctr.gopls & {name: "gopls"}
		}
	}
}
