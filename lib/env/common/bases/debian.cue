package bases

import (
	"github.com/hofstadter-io/hof/lib/env/common/utils"
	"github.com/hofstadter-io/hof/schemas/env"
)

debian13: {

	minimal: env.#Container & {
		#hof: {
			id: "debian13-min"
			metadata: {	
				name: id
				description: "A minimal debian13 image with updates and certs"
			}
		}

		name: string | *"debian13-min"
		from: "debian:13-slim"

		steps: [
			// default workdir (for wide default consistency)
			env.Workdir & {path: "/work"},
			// env.Mount & { path: "/tmp", source: env.#Cache & { name: "debian-13-tmp"}},
			env.Mount & {path: "/var/log", source: env.#Cache & {name: "debian-13-var-log"}},
			env.Mount & {path: "/var/cache", source: env.#Cache & {name: "debian-13-var-cache"}},

			// shared apt caches, for all derived images as well
			// ya'know, instead of cleaning and refetching all the time?
			utils.apt.mounts.varLib,
			// need to update once at the beginning
			utils.apt.update,

			// just certs
			utils.apt.install & {#pkgs: ["ca-certificates", "wget", "curl"]}, // shouldn't need wget/curl, we can do that at this level
		]
	}

	default: env.#Container & {
		#hof: {
			id: "debian13"
			metadata: {	
				name: id
				description: "A default debian13 image with common packages and tools"
			}
		}

		from: "debian:13-slim"

		steps: [
			// default workdir (for wide default consistency)
			env.Workdir & {path: "/work"},
			// env.Mount & { path: "/tmp", source: env.#Cache & { name: "debian-13-tmp"}},
			env.Mount & {path: "/var/log", source: env.#Cache & {name: "debian-13-var-log"}},
			env.Mount & {path: "/var/cache", source: env.#Cache & {name: "debian-13-var-cache"}},

			// shared apt caches, for all derived images as well
			// ya'know, instead of cleaning and refetching all the time?
			utils.apt.mounts.varLib,
			// need to update once at the beginning
			utils.apt.update,

			// basics
			utils.apt.install & {#pkgs: [
				"apt-transport-https",
				"ca-certificates",
				"curl",
				"git",
				"git-absorb",
				"git-lfs",
				"gnupg",
				"jq",
				"lsb-release",
				"make",
				"snap",
				"tree",
				"unzip",
				"wget",
				"xz-utils",
				"zsh",
			]},
		]
	}

}
