package lib

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/build"
	"cuelang.org/go/cue/load"
	"github.com/fatih/color"
	"github.com/mattn/go-zglob"

	"github.com/hofstadter-io/hof/lib/gen"
	"github.com/hofstadter-io/hof/lib/util"
)

type Runtime struct {
	// Setup options
	Entrypoints []string
	Expressions []string

	// TODO configuration
	mode string
	verbose bool

	// Cue ralated
	CueRT           *cue.Runtime
	BuildInstances  []*build.Instance
	CueInstances    []*cue.Instance
	TopLevelValues  []cue.Value
	TopLevelStructs []*cue.Struct

	// Hof related
	Generators map[string]*gen.Generator
}

func NewRuntime(entrypoints, expressions [] string) (*Runtime) {
	return &Runtime {
		Entrypoints: entrypoints,
		Expressions: expressions,

		CueRT: &cue.Runtime{},

		Generators: make(map[string]*gen.Generator),
	}
}

func (R *Runtime) LoadCue() []error {

	var errs []error

	BIS := load.Instances(R.Entrypoints, nil)
	R.BuildInstances = BIS


	for _, bi := range BIS {
		if bi.Err != nil {
			errs = append(errs, bi.Err)
			continue
		}

		// Build the Instance
		I, err := R.CueRT.Build(bi)
		if err != nil {
			errs = append(errs, bi.Err)
			continue
		}
		R.CueInstances = append(R.CueInstances, I)

		// Get top level value from cuelang
		V := I.Value()
		R.TopLevelValues = append(R.TopLevelValues, V)

		// Get top level struct from cuelang
		S, err := V.Struct()
		if err != nil {
			errs = append(errs, err)
			continue
		}
		R.TopLevelStructs = append(R.TopLevelStructs, S)
	}

	R.ExtractHofItems()

	return errs
}

func (R *Runtime) ExtractHofItems() {
	// TODO, what about other things in top level values? or instances?
	// loop ever all top level structs
	for _, S := range R.TopLevelStructs {

		// Loop through all top level fields
		iter := S.Fields()
		for iter.Next() {

			label := iter.Label()
			value := iter.Value()

			// is generator?
			if strings.HasPrefix(label, "HofGen") {
				short := strings.TrimPrefix(label, "HofGen")
				G := gen.NewGenerator(short, value)
				R.Generators[short] = G

				// Disbale if not in expressions
				if len(R.Expressions) > 0 {
					found := false
					for _, expr := range R.Expressions {
						if short == expr {
							found = true
							break
						}
					}

					if !found {
						G.Disabled = true
					}
				}
			}
			// end generator
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

		// Load the Generator!
		errsL := G.LoadCue()
		if len(errsL) != 0 {
			errs = append(errs, errsL...)
			continue
		}

		errsI := G.Initialize()
		if len(errsI) != 0 {
			errs = append(errs, errsI...)
			continue
		}

		// TODO, flatten any nested generators?
	}

	return errs

}

func (R *Runtime) RunGenerators() []error {
	var errs []error

	// Load shadow, can this be done in parallel with the last step?
	shadow, err := gen.LoadShadow(R.verbose)
	if err != nil {
		errs = append(errs, err)
		return errs
	}

	// Don't do in parallel yet, Cue is slow and hungry for memory @ v0.0.16
	for _, G := range R.Generators {
		if G.Disabled {
			continue
		}

		G.Shadow = shadow

		errsG := G.GenerateFiles()
		if len(errsG) > 0 {
			errs = append(errs, errsG...)
			continue
		}

	}

	return errs
}

func (R *Runtime) WriteOutput() []error {
	var errs []error


	for _, G := range R.Generators {
		if G.Disabled {
			continue
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
				continue
			}
			for _, match := range matches {
				// trim first level directory
				clean := Glob[:strings.Index(Glob, "/")]
				mo := strings.TrimPrefix(match, clean)
				src := path.Join(bdir, match)
				dst := path.Join(G.Outdir, mo)

				// normal location
				err := util.CopyFile(src, dst)
				if err != nil {
					err = fmt.Errorf("while copying static file %q\n%w\n", match, err)
					errs = append(errs, err)
					continue
				}

				// shadow location
				err = util.CopyFile(src, path.Join(".hof", dst))
				if err != nil {
					err = fmt.Errorf("while copying static file %q\n%w\n", match, err)
					errs = append(errs, err)
					continue
				}

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
				continue
			}
			err = F.WriteShadow()
			if err != nil {
				errs = append(errs, err)
				continue
			}
		}

		// Finally write the generator files
		for _, F := range G.Files {
			// Write the actual output
			if F.DoWrite {
				err := F.WriteOutput()
				if err != nil {
					errs = append(errs, err)
					continue
				}
			}

			// Write the shadow too, or if it doesn't exist
			if F.DoWrite || (F.IsSame > 0 && F.ShadowFile == nil) {
				err := F.WriteShadow()
				if err != nil {
					errs = append(errs, err)
					continue
				}
			}

			// remove from shadows map so we can cleanup what remains
			delete(G.Shadow, F.Filepath)
		}

		// Cleanup File & Shadow
		for f, _ := range G.Shadow {
			err := os.Remove(f)
			if err != nil {
				errs = append(errs, err)
				continue
			}
			err = os.Remove(path.Join(gen.SHADOW_DIR, f))
			if err != nil {
				errs = append(errs, err)
				continue
			}
			G.Stats.NumDeleted += 1
		}

		writeend := time.Now()
		G.Stats.WritingTime = writeend.Sub(writestart).Round(time.Millisecond)

	}

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
