package mod

import (
	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/repos/cache"
)

func Clean(rflags flags.RootPflagpole) (error) {
	upgradeHofMods()

	return cache.CleanCache()
}
