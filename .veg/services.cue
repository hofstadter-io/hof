package veg

import (
	"github.com/hofstadter-io/hof/schemas/env"
)

svc: gopls: env.#Service & {
	@env()

	name: "gopls"
	#hof: metadata: description: "veg.gopls service"

	// There is also a built in MCP server!

	let _port = flags.ports.gopls
	ports: [{name: "lsp", port: _port}]
	args: ["gopls", "serve", "-port=\(_port)"]

	source: env.#Container & {
		name: "gopls"
		from: ctr.dev
		steps: [
			env.Expose & {port: _port},
			env.Mount & {path: "/work", source: src.code},
		]
	}

}

svc: cuepls: env.#Service & {
	@env()
	name: "cuepls"
	#hof: metadata: description: "veg.cuepls service"

	let _port = flags.ports.cuepls
	ports: [{name: "lsp", port: _port}]
	args: ["cue", "lsp", "serve", "-port=\(_port)"]

	source: env.#Container & {
		name: "cuepls"
		from: ctr.dev
		steps: [
			env.Expose & {port: _port},
			env.Mount & {path: "/work", source: src.code},
		]
	}
}
