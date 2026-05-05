package bases

import (
	"github.com/hofstadter-io/hof/catalogs/env/utils"
	"github.com/hofstadter-io/hof/schemas/env"
)

alpine3: {
	minimal: env.#Container & {
		@env(alpine3-min) // still need to do this, so that AtMade is true for viewing
		from: "alpine:3.23"
		steps: [
			// globally consistent for us
			env.Workdir & {path: "/work"},
			// stuff we don't want in the final image, but need around when building/running
			env.Mount & {path: "/var/log", source: env.#Cache & {name: "alpine3-var-log"}},
			env.Mount & {path: "/var/cache", source: env.#Cache & {name: "alpine3-var-cache"}},

			// update just once at the beginning
			utils.apk.update,
			utils.apk.upgrade,

			// the minimal essenitals
			utils.apk.install & {#pkgs: ["ca-certificates", "wget", "curl", "bash"]},
		]
	}
}
