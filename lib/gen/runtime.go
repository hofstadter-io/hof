package gen

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"cuelang.org/go/cue"
	"github.com/fatih/color"
	"github.com/mattn/go-zglob"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/cuetils"
	"github.com/hofstadter-io/hof/lib/yagu"
)

type Runtime struct {
	// Setup options
	Entrypoints []string
	Flagpole    flags.GenFlagpole

	// TODO configuration
	mode    string
	verbose bool

	// Cue ralated
	CueRuntime      *cuetils.CueRuntime

	// Hof related
	Generators map[string]*Generator
	Shadow     map[string]*File
	Stats      *RuntimeStats
}

func NewRuntime(entrypoints []string, cmdflags flags.GenFlagpole) *Runtime {
	return &Runtime{
		Entrypoints: entrypoints,
		Flagpole:    cmdflags,
		Generators:  make(map[string]*Generator),
		Stats:       new(RuntimeStats),
	}
}

func (R *Runtime) ClearGenerators() {
	R.Generators = make(map[string]*Generator)
}

func (R *Runtime) LoadCue() (err error) {
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

		G := NewGenerator(label, value)
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

		// fmt.Println("Loading Generator:", G.Name)

		// Load the Generator!
		errsL := G.LoadCue()
		if len(errsL) != 0 {
			errs = append(errs, errsL...)
			continue
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

func (R *Runtime) RunGenerators() []error {
	start := time.Now()
	defer func() {
		end := time.Now()
		R.Stats.GenRunningTime = end.Sub(start)
	}()

	var errs []error

	// Load shadow, can this be done in parallel with the last step?
	// Don't do in parallel yet, Cue is slow and hungry for memory @ v0.0.16
	// CUE v0.4.0- is not concurrency safe, maybe v0.4.1 will introduce?
	for _, G := range R.Generators {
		gerrs := R.RunGenerator(G)
		if len(gerrs) > 0 {
			errs = append(errs, gerrs...)
		}
	}

	return errs
}

func (R *Runtime) RunGenerator(G *Generator) (errs []error) {
	if G.Disabled {
		return
	}

	if G.UseDiff3 {
		shadow, err := LoadShadow(G.Name, R.verbose)
		if err != nil {
			errs = append(errs, err)
			return errs
		}
		G.Shadow = shadow
	}

	// run this generator
	errsG := G.GenerateFiles()
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

func (R *Runtime) WriteOutput() []error {
	var errs []error

	for _, G := range R.Generators {
		gerrs := R.WriteGenerator(G)
		errs = append(errs, gerrs...)
	}

	return errs
}

func (R *Runtime) WriteGenerator(G *Generator) (errs []error) {
	if G.Disabled {
		return errs
	}

	writestart := time.Now()

	// Order is important here for implicit overriding of content

	// Start with static file globs
	for _, Static := range G.Statics {
		for _, Glob := range Static.Globs {
			bdir := ""
			if G.PackageName != "" {
				bdir = path.Join("cue.mod/pkg", G.PackageName)
			}
			matches, err := zglob.Glob(path.Join(bdir, Glob))
			if err != nil {
				err = fmt.Errorf("while globbing %s / %s\n%w\n", bdir, Glob, err)
				errs = append(errs, err)
				return errs
			}
			for _, match := range matches {
				mo := strings.TrimPrefix(match, Static.TrimPrefix)
				src := path.Join(bdir, match)
				dst := path.Join(G.Outdir, Static.OutPrefix, mo)

				// TODO?, make comparison and decide to write or not

				// normal location
				err := yagu.CopyFile(src, dst)
				if err != nil {
					err = fmt.Errorf("while copying static file %q\n%w\n", match, err)
					errs = append(errs, err)
					return errs
				}

				if G.UseDiff3 {
					// shadow location
					err = yagu.CopyFile(src, path.Join(SHADOW_DIR, G.Name, dst))
					if err != nil {
						err = fmt.Errorf("while copying static shadow file %q\n%w\n", match, err)
						errs = append(errs, err)
						return errs
					}
				}

				delete(R.Shadow, path.Join(G.Name, dst))
				delete(G.Shadow, path.Join(G.Name, dst))
				G.Stats.NumStatic += 1
				G.Stats.NumWritten += 1

			}

		}
	}

	// Then the static files in cue
	for p, content := range G.EmbeddedStatics {
		F := &File{
			Filepath:     path.Join(G.Outdir, p),
			FinalContent: []byte(content),
		}
		err := F.WriteOutput()
		if err != nil {
			errs = append(errs, err)
			return errs
		}
		if G.UseDiff3 {
			err = F.WriteShadow(path.Join(SHADOW_DIR, G.Name))
			if err != nil {
				errs = append(errs, err)
				return errs
			}
		}
		delete(R.Shadow, path.Join(G.Name, F.Filepath))
		delete(G.Shadow, path.Join(G.Name, F.Filepath))
		G.Stats.NumStatic += 1
		G.Stats.NumWritten += 1
	}

	// Finally write the generator files
	for _, F := range G.Files {
		// Write the actual output
		if F.DoWrite && len(F.Errors) == 0 {
			err := F.WriteOutput()
			if err != nil {
				errs = append(errs, err)
				return errs
			}
		}

		// Write the shadow too, or if it doesn't exist
		if G.UseDiff3 {
			if F.DoWrite || (F.IsSame > 0 && F.ShadowFile == nil) {
				err := F.WriteShadow(path.Join(SHADOW_DIR, G.Name))
				if err != nil {
					errs = append(errs, err)
					return errs
				}
			}
		}

		// remove from shadows map so we can cleanup what remains
		delete(R.Shadow, path.Join(G.Name, F.Filepath))
		delete(G.Shadow, path.Join(G.Name, F.Filepath))
	}

	// Cleanup File & Shadow
	// fmt.Println("Clean Shadow", G.Name)
	if G.UseDiff3 {
		for f, _ := range G.Shadow {
			genFilename := strings.TrimPrefix(f, G.Name+"/")
			shadowFilename := path.Join(SHADOW_DIR, f)
			fmt.Println("  -", G.Name, f, genFilename, shadowFilename)

			err := os.Remove(genFilename)
			if err != nil {
				if strings.Contains(err.Error(), "no such file or directory") {
					continue
				}
				errs = append(errs, err)
				return errs
			}

			err = os.Remove(shadowFilename)
			if err != nil {
				if strings.Contains(err.Error(), "no such file or directory") {
					continue
				}
				errs = append(errs, err)
				return errs
			}

			G.Stats.NumDeleted += 1
		}
	}

	for _, SG := range G.Generators {
		SG.UseDiff3 = G.UseDiff3
		sgerrs := R.WriteGenerator(SG)
		errs = append(errs, sgerrs...)
	}

	writeend := time.Now()
	G.Stats.WritingTime = writeend.Sub(writestart).Round(time.Millisecond)

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
				msg := fmt.Sprint("MERGE CONFLICT in:", F.Filepath)
				color.Red(msg)
			}
		}
	}
}
