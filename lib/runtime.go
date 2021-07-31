package lib

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/build"
	"cuelang.org/go/cue/errors"
	"cuelang.org/go/cue/load"
	"github.com/fatih/color"
	"github.com/mattn/go-zglob"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/cuetils"
	"github.com/hofstadter-io/hof/lib/gen"
	"github.com/hofstadter-io/hof/lib/yagu"
)

type Runtime struct {
	// Setup options
	Entrypoints []string
	Flagpole flags.GenFlagpole

	// TODO configuration
	mode string
	verbose bool

	// Cue ralated
	CueCTX          *cue.Context
	BuildInstances  []*build.Instance
	CueInstances    []*cue.Instance
	TopLevelValues  []cue.Value
	TopLevelStructs []*cue.Struct

	// Hof related
	Generators map[string]*gen.Generator
	Shadow map[string]*gen.File
}

func NewRuntime(entrypoints [] string, cmdflags flags.GenFlagpole) (*Runtime) {
	return &Runtime {
		Entrypoints: entrypoints,
		Flagpole: cmdflags,

		CueCTX: cuecontext.New(),

		Generators: make(map[string]*gen.Generator),
	}
}

func (R *Runtime) LoadCue() []error {

	var errs []error

	BIS := load.Instances(R.Entrypoints, nil)
	R.BuildInstances = BIS


	for _, bi := range BIS {
		if bi.Err != nil || bi.Incomplete {
			fmt.Println("LoadCue:", bi.Err, bi.Incomplete, bi.DepsErrors)
			// TODO add DepsErrors if needed
			es := errors.Errors(bi.Err)
			for _, e := range es {
				errs = append(errs, e.(error))
			}
			continue
		}

		// Build the Instance
		V := R.CueCTX.BuildInstance(bi)
		if V.Err() != nil {
		  es := errors.Errors(V.Err())
			// fmt.Println("BUILD ERR", es, I)
			for _, e := range es {
				errs = append(errs, e.(error))
			}
			continue
		}
		R.TopLevelValues = append(R.TopLevelValues, V)

		// Get top level struct from cuelang
		S, err := V.Struct()
		if err != nil {
			// fmt.Println("STRUCT ERR", err)
		  es := errors.Errors(err)
			for _, e := range es {
				errs = append(errs, e.(error))
			}
			continue
		}
		R.TopLevelStructs = append(R.TopLevelStructs, S)
	}

	if len(errs) > 0 {
		return errs
	}

	R.ExtractGenerators()

	return errs
}

func (R *Runtime) ExtractGenerators() {
	// loop ever all top level structs
	for _, S := range R.TopLevelStructs {

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

					// are there flags to match?
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
					}

					// passed, we should generate
					hasgen = true
					break
				}
			}

			if !hasgen {
				continue
			}

			G := gen.NewGenerator(label, value)
			R.Generators[label] = G
		}
	}
}

func (R *Runtime) LoadGenerators() []error {
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
			fmt.Println("  Load Error:", errsL)
			errs = append(errs, errsL...)
			continue
		}

		// TODO, flatten any nested generators?
		// this would eleminiate all the recursion in other functions
		// would still need it here (in a new func)
	}

	return errs

}

func (R *Runtime) RunGenerators() []error {
	var errs []error
	// var err error

	/*
	R.Shadow, err = gen.LoadShadow("", R.verbose)
	if err != nil {
		errs = append(errs, err)
		return errs
	}
	*/

	// Load shadow, can this be done in parallel with the last step?
	// Don't do in parallel yet, Cue is slow and hungry for memory @ v0.0.16
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

	shadow, err := gen.LoadShadow(G.Name, R.verbose)
	if err != nil {
		errs = append(errs, err)
		return errs
	}

	G.Shadow = shadow

	errsG := G.GenerateFiles()
	if len(errsG) > 0 {
		errs = append(errs, errsG...)
		return errs
	}

	for _, sg := range G.Generators {
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

	// TODO, remove this? do we even have global shadow files? The section to load them is commented out
	// Clean global shadow, incase any generators were removed
	for f, _ := range R.Shadow {
		// deal with leading shadow dir name?
		idx := strings.Index(f,"/")
		if idx < 0 {
			idx = 0
		} else {
			idx += 1
		}
		// fmt.Println("  +", f, idx)
		// fmt.Println("  -", f, f[idx:])
		err := os.Remove(f[idx:])
		if err != nil {
			if strings.Contains(err.Error(), "no such file or directory") {
				continue
			}
			errs = append(errs, err)
			continue
		}

		err = os.Remove(path.Join(gen.SHADOW_DIR, f))
		if err != nil {
			if strings.Contains(err.Error(), "no such file or directory") {
				continue
			}

			errs = append(errs, err)
			continue
		}
	}

	return errs
}

func (R *Runtime) WriteGenerator(G *gen.Generator) (errs []error) {
	if G.Disabled {
		return errs
	}

	writestart := time.Now()

	// Order is important here for implicit overriding of content

	// Start with static file globs
	for _, Glob := range G.StaticGlobs {
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
			// trim first level directory
			clean := Glob[:strings.Index(Glob, "/")]
			mo := strings.TrimPrefix(match, clean)
			src := path.Join(bdir, match)
			dst := path.Join(G.Outdir, mo)

			// TODO, make comparison and decide to write or not

			// normal location
			err := yagu.CopyFile(src, dst)
			if err != nil {
				err = fmt.Errorf("while copying static real file %q\n%w\n", match, err)
				errs = append(errs, err)
				return errs
			}

			// shadow location
			err = yagu.CopyFile(src, path.Join(".hof", G.Name, dst))
			if err != nil {
				err = fmt.Errorf("while copying static shadow file %q\n%w\n", match, err)
				errs = append(errs, err)
				return errs
			}

			delete(R.Shadow, path.Join(G.Name, dst))
			delete(G.Shadow, path.Join(G.Name, dst))
			G.Stats.NumStatic += 1
			G.Stats.NumWritten += 1

		}

	}

	// Then the static files in cue
	for p, content := range G.StaticFiles {
		F := &gen.File {
			Filepath: path.Join(G.Outdir, p),
			FinalContent: []byte(content),
		}
		err := F.WriteOutput()
		if err != nil {
			errs = append(errs, err)
			return errs
		}
		err = F.WriteShadow(path.Join(gen.SHADOW_DIR, G.Name))
		if err != nil {
			errs = append(errs, err)
			return errs
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
		if F.DoWrite || (F.IsSame > 0 && F.ShadowFile == nil) {
			err := F.WriteShadow(path.Join(gen.SHADOW_DIR, G.Name))
			if err != nil {
				errs = append(errs, err)
				return errs
			}
		}

		// remove from shadows map so we can cleanup what remains
		delete(R.Shadow, path.Join(G.Name, F.Filepath))
		delete(G.Shadow, path.Join(G.Name, F.Filepath))
	}

	// Cleanup File & Shadow
	// fmt.Println("Clean Shadow", G.Name)
	for f, _ := range G.Shadow {
		genFilename := strings.TrimPrefix(f, G.Name + "/")
		shadowFilename := path.Join(gen.SHADOW_DIR, f)
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

	for _, SG := range G.Generators {
		sgerrs := R.WriteGenerator(SG)
		errs = append(errs, sgerrs...)
	}

	writeend := time.Now()
	G.Stats.WritingTime = writeend.Sub(writestart).Round(time.Millisecond)

	return errs
}

func (R *Runtime) PrintStats() {
	for _, G := range R.Generators {
		if G.Disabled {
			continue
		}

		G.Stats.CalcTotals(G)
		fmt.Printf("\n%s\n==========================\n", G.Name)
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
