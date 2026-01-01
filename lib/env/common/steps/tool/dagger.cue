package tool

import (
	"github.com/hofstadter-io/hof/lib/env/common/utils"
)

dagger: {
	#ver: string | *"0.19.8"
	cli: [
		utils.githubBin & {#repo: "dagger/dagger", #ver: dagger.#ver},
	]
}
