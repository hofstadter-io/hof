package tool

import (
	"github.com/hofstadter-io/hof/lib/env/common/steps/util"
)

dagger: {
	#ver: string | *"0.19.8"
	cli: [
		util.githubBin & {#repo: "dagger/dagger", #ver: dagger.#ver},
	]
}
