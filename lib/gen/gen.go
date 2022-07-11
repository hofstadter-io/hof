package gen

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"strings"
	"time"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/templates"
	"github.com/hofstadter-io/hof/lib/yagu"
)

func Gen(args []string, rootflags flags.RootPflagpole, cmdflags flags.GenFlagpole) error {
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

func runGen(args []string, rootflags flags.RootPflagpole, cmdflags flags.GenFlagpole) error {

	// This is our runtime for codegeneration
	R := NewRuntime(args, cmdflags)
	R.Verbosity = rootflags.Verbosity

	// b/c shorter names
	LT := len(cmdflags.Template)
	LG := len(cmdflags.Generator)
	globs := cmdflags.WatchGlobs
	xcue := cmdflags.WatchXcue

	// determine watch mode
	//  excplicit: -w
	//  implicit:  -W/-X
	watch := cmdflags.Watch
	if len(globs) > 0 || len(xcue) > 0 {
		watch = true
	}

	// generally there are more messages when in watch mode
	if watch && cmdflags.AsModule == "" {
		fmt.Println("Loading CUE from", R.Entrypoints)
	}

	// this is the first time we create a runtime and load cue
	err := R.LoadCue()
	if err != nil {
		return err
	}

	/* We will run a generator if either
	   not adhoc or is adhoc with the -G flag
		So let's load them early, there is some helpful info in them
	*/
	if LT == 0 || LG > 0 {
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
		err = R.CreateAdhocGenerator()
		if err != nil {
			return err
		}
	}

	if cmdflags.AsModule != "" {
		return R.AsModule()
	}

	/* Build up watch list
		We need to buildup the watch list from flags
		and any generator we might run, which might have watch settings
	*/
	// todo, infer most entrypoints
	for _, arg := range args {
		info, err := os.Stat(arg)
		if err != nil {
			return err
		}
		if info.IsDir() {
			globs = append(globs, info.Name() + "/*")
		} else {
			globs = append(globs, info.Name())
		}
	}

	for _, G := range R.Generators {
		if G.Disabled {
			continue
		}
		globs = append(globs, G.WatchGlobs...)
		xcue = append(xcue, G.WatchXcue...)
		// add templates to full regen globs
		xcue = append(xcue, G.Templates[0].Globs...)
	}
	// add partial templates to xcue globs
	// can do outside loop since all gens have the same value
	xcue = append(xcue, R.Flagpole.Partial...)

	// this might be empty, we calc anyway for ease and sharing
	wfiles, err := yagu.FilesFromGlobs(globs)
	if err != nil {
		return err
	}
	xfiles, err := yagu.FilesFromGlobs(xcue)
	if err != nil {
		return err
	}

	// if we are in watch mode, let the user know what is being watched
	if watch {
		fmt.Printf("found %d glob files from %v\n", len(wfiles), globs)
		fmt.Printf("found %d xcue files from %v\n", len(xfiles), xcue)
	}

	// code gen func
	doGen := func(fast bool) (chan bool, error) {
		if R.Verbosity > 0 {
			fmt.Println("runGen.doGen: fast:", fast)
		}
		return R.genOnce(fast, watch, xfiles)
	}

	// no watch, we can now exit 0
	if !watch {
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

		if R.Flagpole.Stats {
			R.PrintStats()
			fmt.Printf("\nTotal Elapsed Time: %s\n\n", elapsed)
		}

		if len(errsG) > 0 {
			for _, e := range errsG {
				fmt.Println(e)
			}
			return nil, fmt.Errorf("\nErrors while generating output\n")
		}
		if len(errsW) > 0 {
			for _, e := range errsW {
				fmt.Println(e)
			}
			return nil, fmt.Errorf("\nErrors while writing output\n")
		}

		R.PrintMergeConflicts()

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

	go DoWatch(doGen, true, false, files, "xcue", quit)

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
	fmt.Printf("Initializing: %s/%s in pkg %s", module, name, pkg)

	// construct template input data
	data := map[string]interface{}{
		"Module": module,
		"Package": pkg,
		"Name": name,
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
	err = render("cue.mods", cuemodsTemplate)
	if err != nil {
		return err
	}
	err = render("cue.mod/module.cue", cuemodFileTemplate)
	if err != nil {
		return err
	}
	// todo, fetch deps
	msg, err := yagu.Bash("hof mod vendor cue")
	fmt.Println(msg)
	if err != nil {
		return err
	}
	// make some dirs
	dirs := []string{"templates", "partials", "statics", "examples", "gen", "schema"}
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

const newModuleTemplate =`
package {{ .Package }}

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
	WatchGlobs: [...string]
	// This is helpful when authoring generator modules
	WatchXcue:  [...string]

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

	// these are the default globs to load from disk
	Templates: [gen.#Templates & {Globs: ["./templates/**/*"], TrimPrefix: "./templates/"}]
	Partials:  [gen.#Templates & {Globs: ["./partials/**/*"], TrimPrefix: "./partials/"}]
	Statics:   [gen.#Statics & {Globs: ["./static/**/*"], TrimPrefix: "./static/"}]

	// The final list of files for hof to generate
	Out: [...gen.#File] & [
		// fill this with file values
	]

	// you can create any intermediate values you need internally

	// open, so your users can build on this
	...
}
`

