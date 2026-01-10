@experiment(aliasv2)
package hack

import (
	"github.com/hofstadter-io/hof/catalogs/env/packs"
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
		from: packs.lang.go.ctr.base
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

}
