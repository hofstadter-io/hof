package lang

import (
	"github.com/hofstadter-io/hof/lib/env/common/steps/util"
	"github.com/hofstadter-io/hof/schemas/env"
)

cue: {
	#ver: string | *"0.15.1"

	caches: {
		cueMods: env.Volume & {
			name: "cue-mods-\(#ver)"
			type: "cache"
		}
	}

	default: [
		// any env vars?
		// with cache
		install,
	]

	install: [
		util.githubBin & {#repo: "cue-lang/cue", #ver: cue.#ver},
	]
}
