package mod

import (
	"fmt"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
)

func Link(rflags flags.RootPflagpole) (error) {
	upgradeHofMods()

	cm, err := loadRootMod()
	if err != nil {
		return err
	}

	if rflags.Verbosity > 0 {
		fmt.Println("linking deps for:", cm.Module)
	}

	err = cm.ensureCached()
	if err != nil {
		return err
	}

	return cm.Vendor("link", rflags.Verbosity)
}
