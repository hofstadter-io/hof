package ops

import (
	"fmt"
	"io/ioutil"
	// "os"
	"path/filepath"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/format"
	"cuelang.org/go/cue/load"

	// "github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/cuetils"
	"github.com/hofstadter-io/hof/lib/yagu"
)

const JUMP_FILE_NAME = "jumps.cue"

func RunJumpFromArgs(args []string) error {
	// fmt.Println("lib/ops.Jump", args)

	jfn := JUMP_FILE_NAME

	// Check our flags
	//if flags.RootPflags.ResourcesDir != "" {
		//rDir := flags.RootPflags.ResourcesDir
		//jfn = filepath.Join(rDir, JUMP_FILE_NAME)
	//} else if flags.RootPflags.Local {

	//if flags.RootPflags.Local {
		//rDir := "resources"
		//jfn = filepath.Join(rDir, JUMP_FILE_NAME)
	//} else {
		//bDir, err := os.UserConfigDir()
		//if err != nil {
			//return err
		//}
		//rDir := "resources"
		//jfn = filepath.Join(bDir, "hof", rDir, JUMP_FILE_NAME)
	//}
	exists, _ := yagu.CheckPathExists(jfn)
	if !exists {
		content := `
		package resources

		jumps: {}
		`
		yagu.Mkdir(filepath.Dir(jfn))
		ioutil.WriteFile(jfn, []byte(content), 0644)
	}

	var err error

	entrypoints := []string {jfn}

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

	if len(args) == 0 {
		bytes, err := format.Node(rCRT.CueValue.Syntax())
		if err != nil {
			return err
		}
		fmt.Printf("%s\n\n", string(bytes))
		return nil
	}

	rest := []string{}
	for i, arg := range args {
		if arg == "__" {
			rest = args[i+1:]
			args = args[:i]
		}
	}

	// fmt.Println("jump", args, rest)

	args = append([]string{"jumps"}, args...)

	V := rCRT.CueValue.Lookup(args...)

	switch K := V.Kind(); K {

	// we found a jump, now leap!
	case cue.StringKind:
		str, err := V.String()
		if err != nil {
			return err
		}
		return runJumpCommand(str, rest)

	// there's more so print it
	case cue.StructKind:
		bytes, err := format.Node(V.Syntax())
		if err != nil {
			return err
		}

		fmt.Printf("%s\n\n", string(bytes))
		return nil

	// TODO, case for list of strings for multiple commands to run?

	// some sort of error or non-existant path
	case cue.BottomKind:
		return fmt.Errorf("Path not found: %q    (or it has errors)", args)

	// unhandled kind
	default:
		return fmt.Errorf("Unsupported kind '%q' at requested path %q", K, args)
	}

	return nil
}

func runJumpCommand(cmd string, args []string) error {
	// TODO, handle or pass in flags
	// TODO, how to parse? is that why runtime is a thing? ('bash -c'), except we are assuming these are all bash
	// fmt.Printf("running: %q\n", cmd)

	out, err := yagu.BashTmpScriptWithArgs(cmd, args)
	fmt.Print(out)
	return err
}
