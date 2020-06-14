package resources

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"

	"cuelang.org/go/cue/load"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/cuetils"
)

func RunInfoFromArgs(args []string, cmdflags flags.InfoFlagpole) (err error) {
	// fmt.Println("lib/resources.Info", args)

	rDir := ""
	dDir := ""

	// Check our flags
	if flags.RootPflags.ResourcesDir != "" {
		rDir = flags.RootPflags.ResourcesDir
	} else {
		// TODO, look in context / config
		rDir = "resources"
	}
	if flags.RootPflags.DatamodelDir != "" {
		dDir = flags.RootPflags.DatamodelDir
	} else {
		// TODO, look in context / config
		dDir = "datamodel"
	}

	var bPrint, cPrint, wPrint bool

	// if any flags are set, lets be specific, otherwise, make all true
	if cmdflags.Builtin || cmdflags.Custom || cmdflags.Local {
		bPrint, cPrint, wPrint = cmdflags.Builtin, cmdflags.Custom, cmdflags.Local
	} else {
		bPrint, cPrint, wPrint = true, true, true
	}

	fmt.Println("Info", rDir, dDir, bPrint, cPrint, wPrint)

	// print builtin from schema
	err = infoBuiltin(rDir)
	if err != nil {
		return err
	}

	err = infoCustom(rDir)
	if err != nil {
		return err
	}

	err = infoWorkspace(rDir)
	if err != nil {
		return err
	}

	return nil
}

func infoBuiltin(rDir string) error {
	fmt.Println("Builtin Resources")
	fmt.Println("----------------------------")

	fmt.Println()
	return nil
}

func infoCustom(rDir string) error {
	fmt.Println("Custom Resources")
	fmt.Println("----------------------------")

	fmt.Println()
	return nil
}

func infoWorkspace(rDir string) error {
	fmt.Println("Workspace Resources")
	fmt.Println("----------------------------")

	var err error

	entrypoints := []string {}

	fis, err := ioutil.ReadDir(rDir)
	if err != nil {
		return err
	}
	for _, fi := range fis {
		if strings.HasSuffix(fi.Name(), ".cue") {
			fn := filepath.Join(rDir, fi.Name())
			entrypoints = append(entrypoints, fn)
		}
	}

	rCRT := &cuetils.CueRuntime{
		Entrypoints: entrypoints,
		CueConfig: &load.Config{
			ModuleRoot: "",
			Module: "",
			Package: "",
			Dir: "",
		},
	}

	err = rCRT.Load()
	if err != nil {
		return err
	}

	S, err := rCRT.CueValue.Struct()
	if err != nil {
		return err
	}

	rTypes := []string{}
	rElems := map[string][]string{}

	iter := S.Fields()
	for iter.Next() {

		rType := iter.Label()
		value := iter.Value()
		names := []string{}

		R, err := value.Struct()
		if err != nil {
			return err
		}

		rIter := R.Fields()
		for rIter.Next() {
			label := rIter.Label()
			names = append(names, label)
		}

		// sort our names
		sort.Strings(names)

		// add to our accumulators
		rTypes = append(rTypes, rType)
		rElems[rType] = names
	}

	// sort our types
	sort.Strings(rTypes)

	// print in a deterministic order
	for _, rT := range rTypes {
		rE := rElems[rT]
		fmt.Printf("  %-16s  %v\n", rT + ":", rE)
	}

	fmt.Println()
	return nil
}
