package bases

import (
	"github.com/hofstadter-io/hof/lib/env/common/steps/tool"
	"github.com/hofstadter-io/hof/lib/env/common/steps/utils"
	"github.com/hofstadter-io/hof/schemas/env"
)

debian: env.#Container & {
	#hof: metadata: description: "A minimal debian image with a few common tools"
	name: "trixie"
	from: "debian:13-slim"

	steps: [
		// default workdir (for wide default consistency)
		env.Workdir & {path: "/work"},

		// shared apt caches, for all derived images as well
		// ya'know, instead of cleaning and refetching all the time?
		utils.apt.mounts.varLib,
		utils.apt.mounts.varCache,
		// need to update once at the beginning
		utils.apt.update,

		// basics
		utils.apt.install & {#pkgs: [
			"ca-certificates",
			"curl",
			"git",
			"gnupg",
			"make",
			"unzip",
			"wget",
			"xz-utils",
			"zsh",
		]},

		// term customization
		tool.zsh.customize,
	]
}
