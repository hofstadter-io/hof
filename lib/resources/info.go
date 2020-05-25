package resources

import (
	"fmt"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
)

func RunInfoFromArgs(args []string, cmdflags flags.InfoFlagpole) (err error) {
	// fmt.Println("lib/resources.Info", args)

	rDir := ""
	dDir := ""

	// Check our flags
	if flags.RootResourcesDirPflag != "" {
		rDir = flags.RootResourcesDirPflag
	}
	if flags.RootDatamodelDirPflag != "" {
		dDir = flags.RootDatamodelDirPflag
	}

	if flags.RootGlobalPflag {
	} else if flags.RootLocalPflag {
	} else {
	}

	fmt.Println("Info", rDir, dDir)

	// print builtin from schema
	err = infoBuiltin(args)
	if err != nil {
		return err
	}

	err = infoCustom(args)
	if err != nil {
		return err
	}

	err = infoLocal(args)
	if err != nil {
		return err
	}

	err = infoRemote(args)
	if err != nil {
		return err
	}

	return nil
}

func infoBuiltin(args []string) error {
	fmt.Println("lib/resources.Info - builtin", args)

	return nil
}

func infoCustom(args []string) error {
	// fmt.Println("lib/resources.Info - custom", args)

	return nil
}

func infoLocal(args []string) error {
	fmt.Println("lib/resources.Info - local", args)

	return nil
}

func infoRemote(args []string) error {
	// fmt.Println("lib/resources.Info - remote", args)

	return nil
}
