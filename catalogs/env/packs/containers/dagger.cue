package containers

import (
	"github.com/hofstadter-io/hof/catalogs/env/utils"
)

dagger: {
	#ver: string | *"0.19.8"

	// hack for now to maintain consistency in pack.<tool>.cli.install
	cli: install: [
		utils.githubBin & {#repo: "dagger/dagger", #ver: dagger.#ver},
	]
}
