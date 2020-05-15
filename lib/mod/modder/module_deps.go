package modder

import (
	"fmt"
	"strings"
)

func (mod *Module) PrintSelfDeps() error {
	for path, R := range mod.SelfDeps {
		fmt.Println("   ", path, "~", R.OldPath, R.OldVersion, "=>", R.NewPath, R.NewVersion)
	}

	return nil
}

func (mod *Module) LoadSelfDeps() error {
	for path, R := range mod.SelfDeps {
		fmt.Println("   ", path, "~", R.OldPath, R.OldVersion, "=>", R.NewPath, R.NewVersion)

		// create a module first

		// XXX is this the right place for this?
		// TODO Check if already good (i.e. ??? if in vendor and ok)
		// TODO Check mvs system cache in $HOME/.mvs/cache

		// We probably need to start module creating here

		// Handle local replaces
		/*
			if strings.HasPrefix(rep.NewPath, "./") || strings.HasPrefix(rep.NewPath, "../") {
				fmt.Println("Local replace:", rep)
				m := &Module{
					// TODO Think about Replace syntax options and the existence of git
					// XXX  How does go mod handle this question
					Module:         rep.OldPath,
					Version:        rep.OldVersion,
					ReplaceModule:  rep.NewPath,
					ReplaceVersion: rep.NewVersion,
				}

				// TODO ADD par.Work here - clone and ilook for sum..., then do checks and actions
				mdr.addDependency(m)
				continue
			}
		*/
		if strings.HasPrefix(R.NewPath, "./") || strings.HasPrefix(R.NewPath, "../") {
			fmt.Println("Local Replace:", R.OldPath, R.OldVersion, "=>", R.NewPath, R.NewVersion)
			// is it git or not?

			return nil
		}

		// OTHERWISE... it's a remote repository

		// is it git or a package repository? TBD

	}

	return nil
}
