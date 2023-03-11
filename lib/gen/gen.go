package gen

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"strings"
	"time"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/cmd/hof/verinfo"
	"github.com/hofstadter-io/hof/lib/templates"
	"github.com/hofstadter-io/hof/lib/yagu"
	hfmt "github.com/hofstadter-io/hof/lib/fmt"
)

func Gen(args []string, rootflags flags.RootPflagpole, cmdflags flags.GenFlagpole) error {

	// shortcut when user wants to bootstrap a new generator module
	if cmdflags.InitModule != "" {
		return InitModule(args, rootflags, cmdflags)
	}

	// return GenLast(args, rootflags, cmdflags)
	verystart := time.Now()

	err := runGen(args, rootflags, cmdflags)
	if err != nil {
		return err
	}

	// final timing
	veryend := time.Now()
	elapsed := veryend.Sub(verystart).Round(time.Millisecond)

	if cmdflags.Stats {
		fmt.Printf("\nGrand Total Elapsed Time: %s\n\n", elapsed)
	}

	return nil
}

func runGen(args []string, rootflags flags.RootPflagpole, cmdflags flags.GenFlagpole) (err error) {
	// fix default Diff3 flag value when running hof gen
	// needs to be interwoven here, probably?
	// it's usage pattern is specific to our use cases right now, and want diff3 true for generators, but overridable if set to false
	hasDiff3Flag := false
	for _, arg := range os.Args {
		if arg == "--diff3" || arg == "-D" {
			hasDiff3Flag = true
			break
		}
	}

	// We need to set Diff3 default to true
	// when the user supplies generators and does not set flag
	if len(cmdflags.Template) == 0 {
		if !hasDiff3Flag {
			cmdflags.Diff3 = true
		}
	}

	// This is our runtime for code generation
	R, err := NewRuntime(args, rootflags, cmdflags)
	if err != nil {
		return err
	}
	if len(cmdflags.Template) == 0 {
		R.Diff3FlagSet = hasDiff3Flag
	}

	// log cue dirs
	if R.Verbosity > 0 {
		fmt.Println("CueDirs:", R.CueModuleRoot, R.WorkingDir, R.cwdToRoot)
	}

	// b/c shorter names
	LT := len(cmdflags.Template)
	LG := len(cmdflags.Generator)

	// generally there are more messages when in watch mode
	if shouldWatch(cmdflags) && cmdflags.AsModule == "" {
		fmt.Println("Loading CUE from", R.Entrypoints)
	}

	// this is the first time we create a runtime and load cue
	err = R.LoadCue()
	if err != nil {
		return err
	}

	/* We will run a generator if either
	   not adhoc or is adhoc with the -G flag
		So let's load them early, there is some helpful info in them
	*/
	if LT == 0 || LG > 0 {
		if R.Verbosity > 1 {
			fmt.Println("Loading Value Generator")
		}

		// load generators just so we can search for watch lists
		err := R.ExtractGenerators()
		if err != nil {
			return err
		}

		errsL := R.LoadGenerators()
		if len(errsL) > 0 {
			for _, e := range errsL {
				fmt.Println(e)
				// cuetils.PrintCueError(e)
			}
			return fmt.Errorf("\nErrors while loading generators\n")
		}

		if len(R.Generators) == 0 {
			return fmt.Errorf("no generators found")
		}

		if cmdflags.List {
			gens, err := R.ListGenerators()
			if err != nil {
				return err
			}
			if len(gens) == 0 {
				return fmt.Errorf("no generators found")
			}
			fmt.Printf("Available Generators\n  ")
			fmt.Println(strings.Join(gens, "\n  "))
			
			// print gens
			return nil
		}
	}

	if LT > 0 {
		if R.Verbosity > 1 {
			fmt.Println("Loading Adhoc Generator")
		}
		err = R.CreateAdhocGenerator()
		if err != nil {
			return err
		}
	}

	if cmdflags.AsModule != "" {
		return R.AsModule()
	}

	wfiles, xfiles, err := buildWatchLists(R, args, cmdflags)
	if err != nil {
		return err
	}

	// code gen func
	doGen := func(fast bool) (chan bool, error) {
		if R.Verbosity > 0 {
			fmt.Println("runGen.doGen: fast:", fast)
		}
		return R.genOnce(fast, shouldWatch(cmdflags), xfiles)
	}

	// no watch, we can now exit 0
	if !shouldWatch(cmdflags) {
		_, err := doGen(true)
		if err != nil {
			return err
		}
		return nil
	}

	fmt.Println("watching for changes...")

	// we are in watch mode, this loop does a complete reload
	var wg sync.WaitGroup

	// this is our main executor around full regen
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = DoWatch(doGen, false, true, wfiles, "full", make(chan bool, 2))
	}()

	// main process waits here for ctrl-c
	wg.Wait()

	if err != nil {
		return err
	}

	return nil
}

func (R *Runtime) genOnce(fast, watch bool, files []string) (chan bool, error) {
	verystart := time.Now()

	doGen := func(fast bool) (chan bool, error) {
		if R.Verbosity > 0 {
			fmt.Println("genOnce.doGen: fast:", fast)
		}
		
		err := R.Reload(fast)
		if err != nil {
			return nil, err
		}

		// issue #20 - Don't print and exit on error here, wait until after we have written, so we can still write good files
		errsG := R.RunGenerators()
		errsW := R.WriteOutput()

		// final timing
		veryend := time.Now()
		elapsed := veryend.Sub(verystart).Round(time.Millisecond)

		// TODO (correctness)
		// ordering for the remainder of this function is unclear
		hasErr := false

		if len(errsG) > 0 {
			hasErr = true
			for _, e := range errsG {
				fmt.Println(e)
			}
		}
		if len(errsW) > 0 {
			hasErr = true
			for _, e := range errsW {
				fmt.Println(e)
			}
		}

		// TODO (shadow) not sure if we want to clean up gens without error?
		// right now, if any error, then no clean
		if !hasErr {
			errsS := R.CleanupRemainingShadow()
			if len(errsS) > 0 {
				hasErr = true
				for _, e := range errsS {
					fmt.Println(e)
				}
			}
		}

		R.PrintMergeConflicts()

		if R.Flagpole.Stats {
			R.PrintStats()
			fmt.Printf("\nTotal Elapsed Time: %s\n\n", elapsed)
		}

		if hasErr {
			return nil, fmt.Errorf("ERROR: while running geneators")
		}


		return nil, nil
	} // end doGen

	// run code gen
	_, err :=  doGen(fast)
	if err != nil {
		return nil, err
	}

	// return if not watching
	if !watch {
		return nil, nil
	}

	quit := make(chan bool, 2)

	go DoWatch(doGen, true, false, files, "fastWG", quit)

	return quit, err
}

func InitModule(args []string, rootflags flags.RootPflagpole, cmdflags flags.GenFlagpole) error {
	name := cmdflags.InitModule
	module := "hof.io"
	if strings.Contains(name,"/") {
		i := strings.LastIndex(name,"/")
		module, name = name[:i], name[i+1:]
	}
	// possibly extract explicit package
	pkg := name
	if strings.Contains(name,":") {
		i := strings.LastIndex(name,":")
		name, pkg = name[:i], name[i+1:]
	}
	fmt.Printf("Initializing: %s/%s in pkg %s\n", module, name, pkg)

	ver := verinfo.HofVersion
	if !strings.HasPrefix(ver, "v") {
		ver = "v" + ver
	}

	// construct template input data
	data := map[string]interface{}{
		"CueVer": verinfo.CueVersion,
		"HofVer": ver,
		"Module": module,
		"Name": name,
		"Package": pkg,
	}

	// local helper to render and write embedded templates
	render := func(outpath, content string) error {
		if rootflags.Verbosity > 0 {
			fmt.Println("rendering:", outpath)
		}
		ft, err := templates.CreateFromString(outpath, content, nil)
		if err != nil {
			return err
		}
		bs, err := ft.Render(data)
		if err != nil {
			return err
		}
		if outpath == "-" {
			fmt.Println(string(bs))
			return nil
		} else {
			bs, err = hfmt.FormatSource(outpath, bs, "", nil, true)
			if err != nil {
				return err
			}
			if strings.Contains(outpath, "/") {
				dir, _ := filepath.Split(outpath)
				err := os.MkdirAll(dir, 0755)
				if err != nil {
					return err
				}
			}
			return os.WriteFile(outpath, bs, 0644)
		}
	}

	err := render(name + ".cue", newModuleTemplate)
	if err != nil {
		return err
	}
	err = render("cue.mod/module.cue", cuemodFileTemplate)
	if err != nil {
		return err
	}
	// todo, fetch deps
	msg, err := yagu.Shell("hof mod tidy", "")
	fmt.Println(msg)
	if err != nil {
		return err
	}
	// make some dirs
	dirs := []string{"templates", "partials", "statics", "examples", "creators", "gen", "schema"}
	for _, dir := range dirs {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	err = render("-", finalMsg)
	if err != nil {
		return err
	}

	return nil
}

const initMsg = `To run the '{{.Name}}' generator...
  $ hof gen        ... or ...
  $ hof gen{{range .Entrypoints}} {{.}}{{ end }} {{ .Name }}.cue -G {{ .Name }}
`
const newModuleTemplate = `
package {{ snake .Package }}

import (
	"github.com/hofstadter-io/hof/schema/gen"
)

// This is example usage of your generator
{{ camelT .Name }}Example: #{{ camelT .Name }}Generator & {
	@gen({{ .Name }})

	// inputs to the generator
	Data: { ... }
	Outdir: "./out/"
	
	// File globs to watch and trigger regen when changed
	// Normally, a user would set this to their designs / datamodel
	WatchFull: [...string]
	// This is helpful when authoring generator modules
	WatchFast:  [...string]

	// required by examples inside the same module
	// your users do not set or see this field
	PackageName: ""
}


// This is your reusable generator module
#{{ camelT .Name }}Generator: gen.#Generator & {

	//
	// user input fields
	//

	// this is the interface for this generator module
	// typically you enforce schema(s) here
	// Data: _
	// Input: #Input

	//
	// Internal Fields
	//

	// This is the global input data the templates will see
	// You can reshape and transform the user inputs
	// While we put it under internal, you can expose In
	// or you can omit In and skip having a global context
	In: {
		// fill as needed
		...
	}

	// required for hof CUE modules to work
	// your users do not set or see this field
	PackageName: string | *"{{ .Module }}/{{ .Name }}"

	// The final list of files for hof to generate
	// fill this with file values
	Out: [...gen.#File] & [
	]

	// you can create any intermediate values you need internally

	// open, so your users can build on this
	...
}
`
