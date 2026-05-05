package bases

import (
	"github.com/hofstadter-io/hof/catalogs/env/utils"
	"github.com/hofstadter-io/hof/schemas/env"
)

debian13: {

	minimal: env.#Container & {
		@id(debian13-min)
		#hof: {
			metadata: {
				name:        "debian13-min"
				description: "A minimal debian13 image with updates and certs"
			}
		}
		from: "debian:13-slim"

		steps: [
			// default workdir (for wide default consistency)
			env.Workdir & {path: "/work"},
			env.Mount & {path: "/var/log", source: env.#Cache & {name: "debian-13-var-log"}},
			// env.Mount & {path: "/var/cache", source: env.#Cache & {name: "debian-13-var-cache"}},

			// shared apt caches, for all derived images as well
			// ya'know, instead of cleaning and refetching all the time?
			env.Mount & {path: "/var/lib/apt/lists", source: env.#Cache & {name: "debian-13-var-lib-apt-lists"}},

			// need to update once at the beginning
			utils.apt.update,
			utils.apt.upgrade, // upgrade should really happen in the base image from SCRATCH, perhaps we'll make some of those

			// just certs
			utils.apt.install & {#pkgs: ["ca-certificates", "wget", "curl"]}, // shouldn't need wget/curl, we can do that at this level
		]
	}

	default: env.#Container & {
		@id(debian13-default)
		#hof: {
			metadata: {
				name:        "debian13-default"
				description: "A default debian13 image with common packages and tools"
			}
		}
		from: "debian:13-slim"

		// cmd.test.tasks.go.steps.0.0.from.from.steps... 5.0.0.args.2
		steps: [
			// default workdir (for wide default consistency)
			env.Workdir & {path: "/work"},
			env.Mount & {path: "/var/log", source: env.#Cache & {name: "debian-13-var-log"}},
			// env.Mount & {path: "/var/cache", source: env.#Cache & {name: "debian-13-var-cache"}},

			// shared apt caches, for all derived images as well
			// ya'know, instead of cleaning and refetching all the time?
			env.Mount & {path: "/var/lib/apt/lists", source: env.#Cache & {name: "debian-13-var-lib-apt-lists"}},

			// need to update once at the beginning
			utils.apt.update,
			// utils.apt.upgrade, // upgrade should really happen in the base image from SCRATCH, perhaps we'll make some of those

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
				"ssh-client",
				"tree",
				"unzip",
				"vim",
				"wget",
				"xz-utils",
				"zsh",
			]},
		]
	}

}
