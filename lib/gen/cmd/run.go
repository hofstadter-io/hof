package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/datamodel"
	"github.com/hofstadter-io/hof/lib/gen"
	"github.com/hofstadter-io/hof/lib/runtime"
)

func Run(args []string, rflags flags.RootPflagpole, gflags flags.GenFlagpole) error {
	// return GenLast(args, rootflags, cmdflags)
	verystart := time.Now()

	err := run(args, rflags, gflags)
	if err != nil {
		return err
	}

	// final timing
	veryend := time.Now()
	elapsed := veryend.Sub(verystart).Round(time.Millisecond)

	if gflags.Stats {
		fmt.Printf("\nGrand Total Elapsed Time: %s\n\n", elapsed)
	}

	return nil
}

func run(args []string, rflags flags.RootPflagpole, gflags flags.GenFlagpole) error {
	// shortcut when user wants to bootstrap a new generator module
	if gflags.InitModule != "" {
		return InitModule(args, rflags, gflags)
	}

	// create our core runtime
	r, err := runtime.New(args, rflags)
	if err != nil {
		return err
	}
	// upgrade to a generator runtime
	R := NewGenRuntime(r, gflags)

	// log cue dirs
	if R.Flags.Verbosity > 1 {
		fmt.Println("CueDirs:", R.CueModuleRoot, R.WorkingDir, R.CwdToRoot)
	}

	// First time load (not-fast)
	err = R.Reload(false)
	if err != nil {
		return err
	}

	if len(R.Generators) == 0 {
		return fmt.Errorf("no generators found")
	}

	if R.GenFlags.List {
		// TODO...
		// 1. use table printer
		// 2. move this command up, large blocks of this ought
		gens := make([]string, 0, len(R.Generators))
		for _, G := range R.Generators {
			gens = append(gens, G.Hof.Metadata.Name)
		}
		if len(gens) == 0 {
			return fmt.Errorf("no generators found")
		}
		fmt.Printf("Available Generators\n  ")
		fmt.Println(strings.Join(gens, "\n  "))
		
		// print gens
		return nil
	}

	// we need generators loaded at this point
	if R.GenFlags.AsModule != "" {
		return R.adhocAsModule()
	}

	return R.runGen()
}

// Clears and reloads a runtime, rereading inputs and reprocessing everything
// fast determines if the CUE code is reloaded and evaluated or not (fast will skip CUE).
func (R *Runtime) Reload(fast bool) (err error) {
	R.Lock()
	defer R.Unlock()

	if R.Flags.Verbosity > 1 {
		fmt.Printf("Runtime.Reload(%b)\n", fast)
	}

	R.Clear()

	if !fast {
		err = R.Load()
		if err != nil {
			return err
		}
	}

	err = R.localLoad()
	if err != nil {
		return err
	}



	return nil
}

func (R *Runtime) Clear() {
	R.Datamodels = make([]*datamodel.Datamodel, 0, len(R.Datamodels))
	R.Generators = make([]*gen.Generator, 0, len(R.Generators))
}

func EnrichGeneratorBuilder(R *Runtime) func (R *runtime.Runtime, G *gen.Generator) error {

	return func (rt *runtime.Runtime, G *gen.Generator) error {

		if G.Disabled {
			return nil
		}
		// some values to copy from runtime to generator
		G.Verbosity     = R.Flags.Verbosity
		G.Diff3FlagSet  = R.Diff3FlagSet
		G.UseDiff3      = R.GenFlags.Diff3
		G.NoFormat      = R.GenFlags.NoFormat

		// todo, we would like to get rid of these if possible
		G.CueModuleRoot = R.CueModuleRoot
		G.WorkingDir    = R.WorkingDir
		G.CwdToRoot     = R.CwdToRoot

		if R.Flags.Verbosity > 1 {
			fmt.Println("Loading Generator:", G.Hof.Metadata.Name)
		}

		// Load the Generator! (from in memory CUE)
		// this is more of a decode from CUE
		errs := G.DecodeFromCUE()
		if len(errs) != 0 {
			var emsg string
			for _, err := range errs {
				emsg += fmt.Sprintf("%s\n", err.Error())
			}
			return fmt.Errorf("while decoding %s:\n%s", G.Hof.Path, emsg)
		}

		// this should only happen when
		// 1. module author creating example in own module
		// 2. user misconfiguration, so we should inform
		// 3. you are a user doing this in a subdir completely?
		const warnModuleAuthorFmtStr = `
		You are running the %q generator at %q
			with PackageName: ""

		Note, that when running hof from inside a generator module,
		it currently must be run from the root.

		See GitHub issue: https://github.com/hofstadter-io/hof/issues/103
		`

		if G.PackageName == "" {
			if R.Flags.Verbosity > 0 {
				fmt.Printf(warnModuleAuthorFmtStr, G.Hof.Metadata.Name, G.Hof.Path)
			}
		}

		return nil
	}

}
