package cmd

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	dmcmd "github.com/hofstadter-io/hof/lib/datamodel/cmd"
	"github.com/hofstadter-io/hof/lib/gen"
	"github.com/hofstadter-io/hof/lib/runtime"
)

// gen.Runtime extends the common runtime.Runtime
type Runtime struct {
	*runtime.Runtime

	// Setup options
	GenFlags     flags.GenFlagpole
	Diff3FlagSet bool // this is so we can set it to true without and explicit "true"
}

func NewGenRuntime(RT *runtime.Runtime, gflags flags.GenFlagpole) (*Runtime) {
	return &Runtime{
		Runtime:  RT,
		GenFlags: gflags,
	}

}

func (R *Runtime) WriteOutput() []error {
	var errs []error
	if R.Flags.Verbosity > 0 {
		fmt.Println("Writing output")
	}

	for _, G := range R.Generators {
		gerrs := G.Write(R.OutputDir(R.GenFlags.Outdir), R.ShadowDir())
		errs = append(errs, gerrs...)
	}

	return errs
}

const SHADOW_DIR = ".hof/shadow/"

// ShadowDir returns the absolute path to shadow dir for this runtime.
// It accounts for module root and relative directories.
func (R *Runtime) ShadowDir() string {
	return filepath.Join(R.CueModuleRoot, SHADOW_DIR, R.RootToCwd, R.GenFlags.Outdir)
}

func (R *Runtime) localLoad() error {
	err := R.EnrichDatamodels(nil, dmcmd.EnrichDatamodel)
	if err != nil {
		return err
	}

	// the generators to load up
	gens := R.GenFlags.Generator
	// we want to skip any generators
	// if we are in adhoc mode and haven't set -G
	if len(R.GenFlags.Template) > 0 && len(gens) == 0 {
		// this value should not match user data
		// so we effectively omit all generators, besides adhoc
		gens = []string{"HOF_ADHOC_OMIT_GENERATORS"}
	}

	err = R.EnrichGenerators(gens, EnrichGeneratorBuilder(R))
	if err != nil {
		return err
	}

	err = R.Initialize()
	if err != nil {
		return err
	}

	return nil
}

func (R *Runtime) Initialize() error {
	if R.Flags.Verbosity > 1 {
		fmt.Printf("Runtime.Initialize()\n")
	}

	err := R.CreateAdhocGenerator()
	if err != nil {
		return err
	}

	/*
	for _, G := range R.Generators {
		errs := G.Initialize()
		if len(errs) != 0 {
			var emsg string
			for _, err := range errs {
				emsg += fmt.Sprintf("%s\n", err.Error())
			}
			return fmt.Errorf("while initializing %s:\n%s", G.Hof.Path, emsg)
		}
	}
	*/

	return nil
}

func (R *Runtime) RunGenerators() []error {
	start := time.Now()
	defer func() {
		end := time.Now()
		R.Stats.GenRunningTime = end.Sub(start)
	}()

	var errs []error

	// Load shadow, can this be done in parallel with the last step?
	// Don't do in parallel yet, Cue can be slow and hungry for memory
	// CUE is not concurrency safe yet, even if, this doesn't take that long anyway
	for _, G := range R.Generators {
		gerrs := R.RunGenerator(G)
		if len(gerrs) > 0 {
			errs = append(errs, gerrs...)
		}
	}

	return errs
}

func (R *Runtime) RunGenerator(G *gen.Generator) (errs []error) {
	if G.Disabled {
		return
	}

	outputDir := filepath.Join(R.OutputDir(R.GenFlags.Outdir), G.OutputPath())
	shadowDir := filepath.Join(R.ShadowDir(), G.ShadowPath())

	// late load shadow, only if we are going to generate
	err := G.LoadShadow(shadowDir)
	if err != nil {
		return []error{err}
	}

	// fmt.Println(G.CueValue)

	// run this generator
	errsG := G.GenerateFiles(outputDir)
	if len(errsG) > 0 {
		errs = append(errs, errsG...)
		return errs
	}

	// run any subgenerators
	for _, sg := range G.Generators {
		// make sure
		sg.UseDiff3 = G.UseDiff3
		sgerrs := R.RunGenerator(sg)
		if len(sgerrs) > 0 {
			errs = append(errs, sgerrs...)
		}
	}

	return errs
}

func (R *Runtime) PrintStats() {
	// find gens which ran
	gens := []string{}
	for _, G := range R.Generators {
		if !G.Disabled {
			gens = append(gens, G.Name)
		}
	}

	fmt.Printf("\nHof: %s\n==========================\n", "Runtime")
	fmt.Println("\nGens:", gens)
	fmt.Println(R.Stats)

	for _, G := range R.Generators {
		if G.Disabled {
			continue
		}

		G.Stats.CalcTotals(G)
		fmt.Printf("\nGen: %s\n==========================\n", G.Name)
		fmt.Println(G.Stats)
	}
}

