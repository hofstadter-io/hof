package script

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/yagu"
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

func Run(args []string, rflags flags.RootPflagpole, cflags flags.RunFlagpole) error {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	var hlsFiles []string
	for _, a := range args {
		switch filepath.Ext(a) {

		case ".hls", ".txt", "txtar":
			hlsFiles = append(hlsFiles, a)
		}
	}

	if len(hlsFiles) > 0 {
		err := RunHLS(hlsFiles, rflags, cflags)
		if err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("Please specify args of filepath glob(s) to .hls|.txt testscript files")
}

// runs each glob element in order, globs are lexigraphically sorted
func RunHLS(globs []string, rflags flags.RootPflagpole, cflags flags.RunFlagpole) error {
	// fmt.Printf("lib/ops.RunHLS: %v %# v\n", globs, flags.RunFlags)

	for _, glob := range globs {
		err := runHLS(glob, rflags, cflags)
		if err != nil {
			return err
		}
	}

	return nil
}

func runHLS(glob string, rflags flags.RootPflagpole, cflags flags.RunFlagpole) error {
	r := &runtime.Runner{
		// LogLevel: flags.RootVerbosePflag,
		// LogLevel: "",
	}

	keep := true
	if cflags.Mode == "test" {
		keep = cflags.KeepTestdir
	}

	p := runtime.Params{
		Mode:        cflags.Mode,
		Setup:       envSetup,
		Dir:         ".",
		Glob:        glob,
		WorkdirRoot: cflags.Workdir,
		TestWork:    keep,
	}

	if cflags.Workdir != "" {
		if err := yagu.Mkdir(cflags.Workdir); err != nil {
			return err
		}
	}

	runtime.RunT(r, p)

	if r.Failed {
		return fmt.Errorf("failed in %s", glob)
	}

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
