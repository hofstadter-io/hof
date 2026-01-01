package bases

import (
	"github.com/hofstadter-io/hof/lib/env/common/utils"
	"github.com/hofstadter-io/hof/schemas/env"
)

rocky8: {
	minimal: env.#Container & {
		@id(rocky-8-minimal)
		#hof: metadata: description: "A minimal rocky8 image with updates and certs"

		from: "rockylinux:8:minimal"

		steps: [
			// default workdir (for wide default consistency)
			env.Workdir & {path: "/work"},
			env.Mount & {path: "/var/log", source: env.#Cache & {name: "rocky-8-var-log"}},
			env.Mount & {path: "/var/cache", source: env.#Cache & {name: "rocky-8-var-cache"}},

			// shared apt caches, for all derived images as well
			// ya'know, instead of cleaning and refetching all the time?
			utils.dnf.mounts.varLib,
			// need to update once at the beginning
			utils.dnf.update,

			// just certs
			utils.apt.install & {#pkgs: ["ca-certificates", "wget", "curl"]}, // shouldn't need wget/curl, we can do that at this level
		]
	}
}

rocky9: {
	minimal: env.#Container & {
		@id(rocky-9-minimal)
		#hof: metadata: description: "A minimal rocky9 image with updates and certs"

		from: "rockylinux:9:minimal"

		steps: [
			// default workdir (for wide default consistency)
			env.Workdir & {path: "/work"},
			env.Mount & {path: "/var/log", source: env.#Cache & {name: "rocky-9-var-log"}},
			env.Mount & {path: "/var/cache", source: env.#Cache & {name: "rocky-9-var-cache"}},

			// shared apt caches, for all derived images as well
			// ya'know, instead of cleaning and refetching all the time?
			utils.dnf.mounts.varLib,
			// need to update once at the beginning
			utils.dnf.update,

			// just certs
			utils.apt.install & {#pkgs: ["ca-certificates", "wget", "curl"]}, // shouldn't need wget/curl, we can do that at this level
		]
	}

}