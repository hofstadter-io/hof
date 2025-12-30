package bases

import (
	"github.com/hofstadter-io/hof/lib/env/common/steps/tool"
	"github.com/hofstadter-io/hof/lib/env/common/steps/util"
	"github.com/hofstadter-io/hof/schemas/env"
)

debian: env.#Container & {
	#hof: metadata: description: "A minimal debian image with a few common tools"
	from: "debian:13-slim"

	steps: [
		// default workdir (for wide default consistency)
		env.Workdir & {path: "/work"},

		// shared apt caches, for all derived images as well
		// ya'know, instead of cleaning and refetching all the time?
		util.apt.mounts.varLib,
		util.apt.mounts.varCache,
		// need to update once at the beginning
		util.apt.update,

		// basics
		util.apt.install & {#pkgs: [
			"apt-transport-https",
			"ca-certificates",
			"curl",
			"git",
			"git-absorb",
			"git-lfs",
			"gnupg",
			"lsb-release",
			"make",
			"snap",
			"unzip",
			"wget",
			"xz-utils",
			"zsh",
		]},

		// term customization
		tool.zsh.customize,
	]
}
