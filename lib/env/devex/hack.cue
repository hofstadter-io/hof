package devex

import (
	"github.com/hofstadter-io/hof/schemas/env"
)

_reg: "host.docker.internal:5000"

flags: {
	src: string | *"/Users/tony/adk/go" @tag(src)
}

hack: dev: env.#Container & {

	@env()
	name: "hack"
	#hof: metadata: description: "hack dev container"

	from: "\(_reg)/veg-dev:local"

	// shared volume from current directory

	steps: [
		// the code
		hack.src,

		// the lsps
		// env.BindService & {service: hack.gopls},
		// env.BindService & {service: hack.cuepls},
	]

}

hack: {
	src: env.Mount & {
		path: "/work"
		source: env.#HostDir & {
			@env()
			name: "hack-src"
			path: flags.src
			#hof: metadata: description: "hack dev source"
		}
	}
}

// this version produces two entries for hack-src in `hof env list`
// perhaps we can consolidate them and list both paths (all generally)
// hack: {
//   src: env.#HostDir & {
//     @env()
//     name: "hack-src"
//     path: string | *"/Users/tony/adk/go" @tag(src)
//     #hof: metadata: description: "hack dev source"
//   }
//   srcMount: [
//     env.Mount & { path: "/work", source: hack.src.$out },
//     env.Workdir & { path: "/work" },
//   ]
// }

hack: gopls: env.#Service & {
	@env()

	// There is also a built in MCP server!

	name: "gopls"
	#hof: metadata: description: "hack.gopls service"

	_port: 4000
	ports: [{
		name: "lsp"
		port: _port
	}]

	args: ["gopls", "serve", "-port=\(_port)"]

	source: env.#Container & {
    name: "veg-dev"
		from: "\(_reg)/\(name):local"

		steps: [
			env.Expose & {port: _port},
			hack.src,
		]
	}

}

hack: cuepls: env.#Service & {
	@env()
	name: "cuepls"
	#hof: metadata: description: "hack.cuepls service"

	_port: 4001
	ports: [{
		name: "lsp"
		port: _port
	}]

	args: ["cue", "lsp", "serve", "-port=\(_port)"]

	source: env.#Container & {
    name: "veg-dev"
		from: "\(_reg)/\(name):local"
		steps: [
			env.Expose & {port: _port},
			hack.src,
		]
	}
}
