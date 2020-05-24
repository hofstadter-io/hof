package workspace

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/hofstadter-io/hof/lib/mod"
	"github.com/hofstadter-io/hof/lib/yagu"
)

func RunInitFromArgs(module, name string) error {
	fmt.Println("lib/workspace.Init", module, name)

	parts, err := CheckSplitModuleName(module)
	if err != nil {
		return err
	}

	// if no name supplied, default to last part of module format
	if name == "" {
		name = parts[2]
	}

	// current directory
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	wdps := strings.Split(wd, "/")
	cwd := wdps[len(wdps)-1]

	// check module name against dirctory so we know what to do
	// if name is not the current directory, we want to create and cd into it
	if name != cwd {
		// check directory exists
		exists, err := yagu.CheckPathExists(name)
		if err != nil {
			return err
		}
		if exists {
			fmt.Println("Directory with name already exists: ", name)
			return nil
		}

		err = yagu.Mkdir(name)
		if err != nil {
			return err
		}

		err = os.Chdir(name)
		if err != nil {
			return err
		}
	}

	initd, err := yagu.CheckPathExists(".hofcfg.cue")
	if err != nil {
		return err
	}
	if initd {
		fmt.Println("Workspace already initialized")
		return nil
	}

	err = initWorkspaceDirs()
	if err != nil {
		return err
	}

	err = initWorkspaceFiles(module, name)
	if err != nil {
		return err
	}

	err = initWorkspaceMods(module, name)
	if err != nil {
		return err
	}

	// get latest CWD
	nwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// add workspace to global contexts
	err = addWorkspaceToGlobalContext(name, nwd)
	if err != nil {
		return err
	}

	return nil
}

func initWorkspaceDirs() error {

	dirs := []string{
		"design",
		"models",
		"resources",
	}

	for _, dir := range dirs {
		err := yagu.Mkdir(dir)
		if err != nil {
			return err
		}
	}

	return nil
}

func initWorkspaceFiles(module, name string) error {

	cfg := fmt.Sprintf(initialHofcfg, module, name)
	err := ioutil.WriteFile(".hofcfg.cue", []byte(cfg), 0644)
	if err != nil {
		return err
	}

	// Check the other files are not created before writing
	initd, err := yagu.CheckPathExists(".hofshh.cue")
	if err != nil {
		return err
	}
	if !initd {
		shh := fmt.Sprintf(initialHofshh, module, name)
		err := ioutil.WriteFile(".hofshh.cue", []byte(shh), 0644)
		if err != nil {
			return err
		}
	}

	// Check the other files are not created before writing
	initd, err = yagu.CheckPathExists(".hofctx.cue")
	if err != nil {
		return err
	}
	if !initd {
		ctx := fmt.Sprintf(initialHofctx, module, name)
		err := ioutil.WriteFile(".hofctx.cue", []byte(ctx), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

const initialHofctx = `
current: first

first: {
  Module: "%s"
  Workspace: "%s"
}
`

const initialHofcfg = `
Module: "%s"
Name: "%s"
Dir:  "./"
`

const initialHofshh = `
// put secrets in here

workspace: {
  %s: {
    secret: "change me now that you see me"
  }
}
`

func initWorkspaceMods(module, name string) error {
	// make sure we have loaded modder info
	mod.InitLangs()

	// initialize a Cue module
	err := mod.Init("cue", module)
	if err != nil {
		return err
	}

	return nil
}

func addWorkspaceToGlobalContext(name, dir string) error {

	return nil
}

