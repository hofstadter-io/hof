package cmd

import (
	"fmt"

	"cuelang.org/go/cue"

	"github.com/hofstadter-io/hof/lib/cuetils"
	"github.com/hofstadter-io/hof/lib/datamodel"
	"github.com/hofstadter-io/hof/lib/gen"
	"github.com/hofstadter-io/hof/lib/runtime"
)

func EnrichDatamodelBuilder(R *Runtime) func (R *runtime.Runtime, DM *datamodel.Datamodel) error {
	return func (rt *runtime.Runtime, dm *datamodel.Datamodel) error {
		err := dm.LoadHistory()
		if err != nil {
			return err
		}
		err = dm.CalcDiffs()
		if err != nil {
			return err
		}

		err = dm.EnrichValue()
		if err != nil {
			return err
		}

		// fill back value so available at root runtime value when decoding generators
		val := dm.Value.LookupPath(cue.ParsePath(dm.Hof.Label))
		R.Value = R.Value.FillPath(cue.ParsePath(dm.Hof.Path), val)
		
		return nil
	}
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
		// this is more of a decode from CUE, maybe too much and needs to be split up?
		// (all of it probably deserves to be within this Enrich function
		errs := G.DecodeFromCUE(R.Value)
		if len(errs) != 0 {
			var emsg string
			for _, err := range errs {
				err = cuetils.ExpandCueError(err)
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
			with ModuleName: ""

		Note, that when running hof from inside a generator module,
		it currently must be run from the root.

		See GitHub issue: https://github.com/hofstadter-io/hof/issues/103
		`

		if G.PackageName == "" {
			if R.Flags.Verbosity > 0 {
				fmt.Printf(warnModuleAuthorFmtStr, G.Hof.Metadata.Name, G.Hof.Path)
			}
		}


		// TODO, inject datamodel history into generator input, as needed
		// 1. discover any DM nodes inside our generator input
		// 2. if found, look up the DM in Runtime and merge with In at that point
		// 3. need to walk DM nodes for history, and merge at correct points
		// 4. do we need to remerge G.In into F.In, or should we delay this until render time
		//    what about needing to recurse to find where the value actually changed?
		//    we have an open issue about creating different diff formats and embedding them all
		//
		// also deal with Ordered nodes, this should be one (set) of functions to handle this
		//
		// can we avoid merging in CUE and instead merge in Go maps?
		//
		// we should write various functions for this and call where necessary
		// we may need History earlier, for outfile name interpolation, and may be able to skip here
		// maybe want to do late, so that we can avoid many steps on a file when we (eventually) check inputs for difference
		// [ FOR NOW, do everything in this Enrich function ]

		/*
		in := G.CueValue.LookupPath(cue.ParsePath("In"))
		if !in.Exists() {
			return fmt.Errorf("In gen:%s, missing In value", G.Name)
		}

		err = in.Decode(&something)
		if err != nil {
			return err
		}
		*/

		return nil
	}

}
