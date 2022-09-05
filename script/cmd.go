package script

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/script/runtime"
)

/*
var (
	MODE = flat.String("mode", "test", "runmode (experimental)")
	ENV = flag.String("env", ".env", "environment file to expose to scripts")
	DIR = flag.String("dir", "", "base directory to search for scripts")
	FILES = flag.String("files", "tests/*.hls", "filepath glob for matching script names from the base dir")
	WORK  = flag.String("work", ".hls/tests", "working directory for the scripts")
	LEAVE = flag.Bool("leave", false, "leave the script workdir in place for inspection")
)
*/

func RunRunFromArgs(args []string) error {

	defer func() {
		if r := recover(); r != nil {
			// fmt.Println("Recovered in f", r)
		}
	}()

	var cueFiles, hlsFiles, tbdFiles []string
	for _, a := range args {
		switch filepath.Ext(a) {

		case ".cue":
			cueFiles = append(cueFiles, a)
		case ".hls":
			hlsFiles = append(hlsFiles, a)

		default:
			tbdFiles = append(tbdFiles, a)
		}
	}

	// only dealing with hls files for now,
	// cue only should be pretty straight forward too
	// not really sure what a mix means yet
	//   for now just return error

	if len(cueFiles) > 0 && len(hlsFiles) > 0 {
		return fmt.Errorf("Cannot specify both cue and hls files at the same time yet. Pleases comment on what you think this should mean on github. Issue... ")
	}

	if len(cueFiles) > 0 {
		err := RunCUE(cueFiles)
		if err != nil {
			return err
		}
		return nil
	}

	if len(hlsFiles) > 0 {
		err := RunHLS(hlsFiles)
		if err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("Please specify args of filepath glob(s) of .cue XOR .hls globs")
	// return nil
}

func RunCUE(globs []string) error {
	fmt.Printf("lib/ops.RunCUE: %v %# v\n", globs, flags.RunFlags)

	return nil
}

// runs each glob element in order, globs are lexigraphically sorted
func RunHLS(globs []string) error {
	// fmt.Printf("lib/ops.RunHLS: %v %# v\n", globs, flags.RunFlags)

	for _, glob := range globs {
		err := runHLS(glob)
		if err != nil {
			return err
		}
	}

	return nil
}

func runHLS(glob string) error {
	r := runtime.Runner{
		// LogLevel: flags.RootVerbosePflag,
		LogLevel: "yes please",
	}

	p := runtime.Params{
		Mode:        "run",
		Setup:       envSetup,
		Dir:         ".",
		Glob:        glob,
		WorkdirRoot: ".",
		TestWork:    true,
	}

	runtime.RunT(r, p)

	// TODO check output / status?

	return nil
}

func envSetup(env *runtime.Env) error {

	// .env can contain lines of ENV=VAR
	content, err := ioutil.ReadFile(".env")
	if err != nil {
		// ignore errors, as the file likely doesn't exist
		return nil
	}

	for _, line := range strings.Split(string(content), "\n") {
		if strings.Contains(line, "=") {
			// todo, trim space here
			if line[0:1] == "#" {
				continue
			}
			env.Vars = append(env.Vars, line)
		}
	}

	return nil
}
