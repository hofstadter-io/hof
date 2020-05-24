package workspace

import (
	"fmt"
	"strings"

	"github.com/hofstadter-io/hof/lib/yagu"
)

func RunCloneFromArgs(module, name string) error {
	fmt.Println("lib/workspace.Clone", module, name)

	// look for a version
	version := ""
	if strings.Contains(module, "@") {
		flds := strings.Split(module, "@")
		module, version = flds[0], flds[1]
	}

	parts, err := CheckSplitModuleName(module)
	if err != nil {
		return err
	}

	// if no name supplied, default to last part of module format
	if name == "" {
		name = parts[2]
	}

	// check directory exists
	exists, err := yagu.CheckPathExists(name)
	if err != nil {
		return err
	}
	if exists {
		fmt.Println("Directory with name already exists: ", name)
		return nil
	}

	err = cloneWorkspace(module, version, name)
	if err != nil {
		return err
	}

	return nil
}

func cloneWorkspace(module, version, name string) error {
	// prefix with https
	module = "https://" + module

	err := yagu.CloneRepoIntoDir(module, version, name)
	if err != nil {
		return err
	}

	return nil
}
