package gen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"cuelang.org/go/cue"
	"github.com/fatih/color"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/cuetils"
)

type Runtime struct {
	sync.Mutex

	// Setup options
	Entrypoints  []string
	Flagpole     flags.GenFlagpole
	Diff3FlagSet bool

	// TODO configuration
	mode      string
	Verbosity int
	NoFormat  bool

	// Cue ralated
	CueRuntime    *cuetils.CueRuntime
	CueModuleRoot string
	WorkingDir    string
	rootToCwd     string  // module root -> working dir (foo/bar)
	cwdToRoot     string  // module root <- working dir (../..)

	// Create related
	OriginalWkdir string

	// Hof related
	Generators map[string]*Generator
	Stats      *RuntimeStats
}

func NewRuntime(entrypoints []string, rootflags flags.RootPflagpole, cmdflags flags.GenFlagpole) (*Runtime, error) {
	R := &Runtime{
		Entrypoints: entrypoints,
		Flagpole:    cmdflags,
		Generators:  make(map[string]*Generator),
		Stats:       new(RuntimeStats),
		Verbosity:  rootflags.Verbosity,
		NoFormat:   cmdflags.NoFormat,
	}

	var err error

	// calc cue dirs
	R.CueModuleRoot, err = cuetils.FindModuleAbsPath("")
	if err != nil {
		return R, err
	}
	// TODO: we could make this configurable
	R.WorkingDir, _ = os.Getwd()
	if R.CueModuleRoot != "" {
		R.cwdToRoot, err = filepath.Rel(R.WorkingDir, R.CueModuleRoot)
		if err != nil {
			return R, err
		}
		R.rootToCwd, err = filepath.Rel(R.CueModuleRoot, R.WorkingDir)
		if err != nil {
			return R, err
		}
	}

	return R, nil
}

// OutputDir returns the absolute path to output dir for this runtime.
// Generators will make subdir contributions at read/write time
func (R *Runtime) OutputDir() string {
	if strings.HasPrefix(R.Flagpole.Outdir, "/") {
		return R.Flagpole.Outdir
	}
	return filepath.Join(R.CueModuleRoot, R.rootToCwd, R.Flagpole.Outdir)
}

// ShadowDir returns the absolute path to shadow dir for this runtime.
// Generators will make subdir contributions at read/write time
func (R *Runtime) ShadowDir() string {
	return filepath.Join(R.CueModuleRoot, SHADOW_DIR, R.rootToCwd, R.Flagpole.Outdir)
}

// Clears and reloads a runtime, rereading inputs and reprocessing everything
// fast determines if the CUE code is reloaded and evaluated or not (fast is not).
// These modes correspond to the -W (full) and -X (fast) watch flags
func (R *Runtime) Reload(fast bool) error {
	R.Lock()
	defer R.Unlock()

	if !fast {
		err := R.LoadCue()
		if err != nil {
			return err
		}
	}

	R.ClearGenerators()

	err := R.ExtractGenerators()
	if err != nil {
		return err
	}

	errsL := R.LoadGenerators()
	if len(errsL) > 0 {
		for _, e := range errsL {
			fmt.Println(e)
		}
		return fmt.Errorf("\nErrors while loading generators\n")
	}

	if len(R.Flagpole.Template) > 0 {
		err = R.CreateAdhocGenerator()
		if err != nil {
			return err
		}
	}

	return nil
}

func (R *Runtime) ClearGenerators() {
	R.Generators = make(map[string]*Generator)
}

func (R *Runtime) LoadCue() (err error) {
	if R.Verbosity > 0 {
		fmt.Println("Loading CUE from:", R.Entrypoints)
	}
	start := time.Now()
	defer func() {
		end := time.Now()
		R.Stats.CueLoadingTime = end.Sub(start)
	}()

	R.CueRuntime, err = cuetils.CueRuntimeFromEntrypointsAndFlags(R.Entrypoints)
	if err != nil {
		return err
	}

	return nil
}

func (R *Runtime) ListGenerators() (gens []string, err error) {
	// conditions which mean we should list all
	anyGen := len(R.Flagpole.Generator) == 1 && R.Flagpole.Generator[0] == "*"
	notGen := len(R.Flagpole.Generator) == 0
	allGen := anyGen || notGen

	// loop ever all top level structs
	S, err := R.CueRuntime.CueValue.Struct()
	if err != nil {
		return gens, err
	}

	// Loop through all top level fields
	iter := S.Fields()
	for iter.Next() {

		// label := iter.Label()
		value := iter.Value()
		attrs := value.Attributes(cue.ValueAttr)

		// find top-level with gen attr
		for _, A := range attrs {
			// does it have "@gen()"
			if A.Name() == "gen" {

				// if -G '*', then we skip the following checks
				if allGen {
					gens = append(gens, A.Contents())
				} else {
					// some -G was set, but was not '*'
					if len(R.Flagpole.Generator) > 0 {
						vals := cuetils.AttrToMap(A)
						for _, g := range R.Flagpole.Generator {
							// matched attribute contents
							// todo, use regex or double**
							// or has prefix?
							if _, ok := vals[g]; ok {
								gens = append(gens, A.Contents())
								break
							}
						}

					} else {
						fmt.Println("can we even get here?")
						// gens = append(gens, A.Contents())
					}
				}
			}
		}

	}

	return gens, nil
}

func (R *Runtime) ExtractGenerators() error {
	allGen := len(R.Flagpole.Generator) == 1 && R.Flagpole.Generator[0] == "*"
	hasT := len(R.Flagpole.Template) > 0

	// loop ever all top level structs
	S, err := R.CueRuntime.CueValue.Struct()
	if err != nil {
		return err
	}

	// Loop through all top level fields
	iter := S.Fields()
	for iter.Next() {

		label := iter.Label()
		value := iter.Value()
		attrs := value.Attributes(cue.ValueAttr)

		// find top-level with gen attr
		hasgen := false
		for _, A := range attrs {
			// does it have "@gen()"
			if A.Name() == "gen" {

				// if -G '*', then we skip the following checks
				if !allGen {
					// some -G was set, but was not '*'
					if len(R.Flagpole.Generator) > 0 {
						vals := cuetils.AttrToMap(A)
						match := false
						for _, g := range R.Flagpole.Generator {
							if _, ok := vals[g]; ok {
								match = true
								break
							}
						}

						if !match {
							continue
						}
					} else {
						// not -G was set, if a -T was set...
						// we are in adhoc mode and skip all gens
						// (hmmm) will we even get here?
						//   an earlier shortcircuit may prevent this
						//   this is defensive anyhow
						if hasT {
							continue
						}
					}
				}
				// passed, we should generate
				hasgen = true
				break
			}
		}

		if !hasgen {
			continue
		}

		G := NewGenerator(label, value, R)
		R.Generators[label] = G
	}

	return nil
}

func (R *Runtime) LoadGenerators() []error {
	start := time.Now()
	defer func() {
		end := time.Now()
		R.Stats.GenLoadingTime = end.Sub(start)
	}()

	var errs []error

	// Don't do in parallel yet, Cue is slow and hungry for memory @ v0.0.16
	for _, G := range R.Generators {
		if G.Disabled {
			continue
		}
		// some values to copy from runtime to generator
		G.verbosity = R.Verbosity
		G.Diff3FlagSet = R.Diff3FlagSet
		G.UseDiff3     = R.Flagpole.Diff3

		if R.Verbosity > 1 {
			fmt.Println("Loading Generator:", G.Name)
		}

		// Load the Generator! (from in memory CUE)
		// this is more of a decode from CUE
		errsL := G.DecodeFromCUE()
		if len(errsL) != 0 {
			errs = append(errs, errsL...)
			continue
		}

		// this should only happen when
		// 1. module author creating example in own module
		// 2. user misconfiguration, so we should inform
		// 3. you are a user doing this in a subdir completely?
		const warnModuleAuthorFmtStr = `
		You are running the '%s' generator
			with PackageName: ""

		Note, that when running hof from inside a generator module,
		it currently must be run from the root.

		See GitHub issue: https://github.com/hofstadter-io/hof/issues/103
		`

		if G.PackageName == "" {
			if R.Verbosity > 0 {
				fmt.Printf(warnModuleAuthorFmtStr, G.Name)
			}
		}

		// TODO, flatten any nested generators?
		// this would eleminiate all the recursion in other functions
		// would still need it here (in a new func)
	}

	/* from previous file */
	// NOTE3: maybe this goes here, and we make R "AdhocGen" aware
	// if LT > 0 {  R.CreateAdhocGenerator(rootflags, cmdflags) }

	// TODO, NOTE2: we should override gen2subgen withing this call
	// we might need NOTE3 to pass adhoc partials into gens and subgens
	/* from previous file */

	err := R.CreateAdhocGenerator()
	if err != nil {
		errs = append(errs, err)
	}
	// TODO, consider merging adhoc templates / partials
	// into generators, so we might override or fill in
	// this could enable more powerful reuse by allowing
	// a generator to use anoather "generic" generator module,
	// which itself, would capture a pattern or algorithm?
	// this "generic" module would be usable across targets
	// NOTE, this might just be the location where adhoc
	// can fill things in, see NOTE2 above for gen2subgen

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

func (R *Runtime) PrintMergeConflicts() {
	for _, G := range R.Generators {
		if G.Disabled {
			continue
		}

		for _, F := range G.Files {
			if F.IsConflicted > 0 {
				msg := fmt.Sprint("MERGE CONFLICT in: ", F.Filepath)
				color.Red(msg)
			}
		}
	}
}
