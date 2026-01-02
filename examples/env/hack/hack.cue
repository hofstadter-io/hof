@experiment(aliasv2)
package hack

import (
	"github.com/hofstadter-io/hof/examples/env/veg"
	"github.com/hofstadter-io/hof/catalogs/env/packs/lang"
	"github.com/hofstadter-io/hof/catalogs/env/utils"
	"github.com/hofstadter-io/hof/schemas/env"
)

_flags: {
	reg: string | *"localhost:5000"     @tag(reg)
	src: string | *"/Users/tony/adk/go" @tag(src)
}

hack: {
	[string]~(k,_): {@env(), name: k}
	dev: env.#Container & {
		from: lang.go.ctr.base
		steps: [
			// the code
			env.Mount & {path: "/work", source: hack.src},
			utils.apt.install & {#pkgs: ["netcat-openbsd"]},

			// the lsps
			env.BindService & {service: hack.gopls},
			env.BindService & {service: hack.cuepls},
		]
	}

	src: env.#HostDir & {path: _flags.src}

	// gopls: lang.go.svc.gopls & {
	//   @env()
	// }

	gopls: env.#Service & {
		_port: 4000
		ports: [{name: "lsp", port: _port}]
		args: ["gopls", "serve", "-port=\(_port)"]
		source: env.#Container & {
			name: "veg-dev"
			from: veg.ctr.dev

			steps: [
				env.Expose & {port: _port},
				env.Mount & {path: "/work", source: hack.src},
			]
		}
	}

	cuepls: env.#Service & {
		_port: 4001
		ports: [{name: "lsp", port: _port}]
		args: ["cue", "lsp", "serve", "-port=\(_port)"]
		source: env.#Container & {
			name: "veg-dev"
			from: veg.ctr.dev
			steps: [
				env.Expose & {port: _port},
				env.Mount & {path: "/work", source: hack.src},
			]
		}
	}

}
