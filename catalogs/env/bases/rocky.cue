package bases

import (
	"github.com/hofstadter-io/hof/catalogs/env/utils"
	"github.com/hofstadter-io/hof/schemas/env"
)

rocky8: {
	minimal: env.#Container & {
		@id(rocky8-min)
		#hof: metadata: {
			name:        "rocky8-min"
			description: "A minimal rocky8 image with updates and certs"
		}

		from: "rockylinux:8.9"

		steps: [
			// default workdir (for wide default consistency)
			env.Workdir & {path: "/work"},
			env.Mount & {path: "/var/log", source: env.#Cache & {name: "rocky-8-var-log"}},
			env.Mount & {path: "/var/cache", source: env.#Cache & {name: "rocky-8-var-cache"}},

			// shared apt caches, for all derived images as well
			// ya'know, instead of cleaning and refetching all the time?
			// utils.dnf.mounts.varLib,
			// need to update once at the beginning
			utils.dnf.update,
			// utils.dnf.upgrade,

			// just certs
			utils.dnf.install & {#pkgs: ["ca-certificates", "wget", "curl"]}, // shouldn't need wget/curl, we can do that at this level
		]
	}
}

rocky9: {
	minimal: env.#Container & {
		@id(rocky9-min)
		#hof: metadata: {
			name:        "rocky9-min"
			description: "A minimal rocky9 image with updates and certs"
		}
		name: "rocky9-min"
		from: "rockylinux:9.3"

		steps: [
			// default workdir (for wide default consistency)
			env.Workdir & {path: "/work"},
			env.Mount & {path: "/var/log", source: env.#Cache & {name: "rocky-9-var-log"}},
			env.Mount & {path: "/var/cache", source: env.#Cache & {name: "rocky-9-var-cache"}},

			// shared apt caches, for all derived images as well
			// ya'know, instead of cleaning and refetching all the time?
			// utils.dnf.mounts.varLib,
			// need to update once at the beginning
			utils.dnf.update,
			// utils.dnf.upgrade,

			// bare essentials (ca-certs & curl already installed)
			utils.dnf.install & {#pkgs: ["wget"]}, // shouldn't need wget/curl, we can do that at this level
		]
	}

}
