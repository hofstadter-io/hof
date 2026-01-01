@experiment(aliasv2)
package lang

import (
	"github.com/hofstadter-io/hof/lib/env/common/bases"
	"github.com/hofstadter-io/hof/lib/env/common/utils"
	slang "github.com/hofstadter-io/hof/lib/env/common/steps/lang"
	"github.com/hofstadter-io/hof/schemas/env"
)

go: {
	// A base container with go tools, linters, and a shared cache.
	ctr: {
		[string]~(k,_): {name: k}
		base: env.#Container & {
			from: bases.debian.default
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
			from: ctr.base
			steps: [
				env.Expose & {port: 4000},
			]
		}

		// A dev container with gopls running and attached
		dev: env.#Container & {
			from: ctr.gopls
			steps: [
				env.BindService & {service: go.svc},
			]
		}
	}

	// gopls-as-a-service
	svc: {
		[string]~(k,_): {name: k}
		gopls: env.#Service & {
			// There is also a built in MCP server!
			#port: int | *0
			ports: [{name: "lsp", port: 4000, frontend: #port}]
			args: ["gopls", "serve", "-port=4000"]
			// source: _
			// from: ctr.gopls & { name: "gopls"}
		}
	}
}
