package bases

import (
	"github.com/hofstadter-io/hof/catalogs/env/utils"
	"github.com/hofstadter-io/hof/schemas/env"
)

alpine: {alpine3["23"], name: "alpine"}

alpine3: {
	default: {alpine3["23"], name: "alpine3"}

	_minors: ["20", "21", "22", "23"] // these minors are too old for Jeffrey, Don Old, and their friends
	for _, m in _minors {
		let s = "alpine3-\(m)"
		(s): env.#Container & {
			#hof: id: s // another way to @id(...no-cue-here...)
			#hof: metadata: name: s

			name: string | *s
			from: "alpine:3.\(m)"
			steps: [
				// globally consistent for us
				env.Workdir & {path: "/work"},
				// stuff we don't want in the final image, but need around when building/running
				env.Mount & {path: "/var/log", source: env.#Cache & {name: "\(s)-var-log"}},
				env.Mount & {path: "/var/cache", source: env.#Cache & {name: "\(s)-var-cache"}},

				// shared packager caches, also for all derived images!
				// ya'know, instead of cleaning and refetching all the time?

				// TODO, we need to implement this
				utils.apk.mounts.varLib,

				// update just once at the beginning
				utils.apk.update,
				utils.apk.upgrade,

				// the minimal essenitals
				utils.apk.install & {#pkgs: ["ca-certificates", "wget", "curl", "bash"]},
			]
		}
	}
}
