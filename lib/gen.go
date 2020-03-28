package lib

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/load"
	"github.com/fatih/color"

	"github.com/hofstadter-io/hof/lib/gen"
)

func Gen(entrypoints, expressions []string, mode string) (string, error) {
	verystart := time.Now()

	verbose := false

	GS, err := extractGenerators(entrypoints)
	if err != nil {
		return "", err
	}

	var errs []error

	// Don't do in parallel yet, Cue is slow and hungry for memory @ v0.0.16
	for _, G := range GS {
		// TODO compare against expressions

		err := G.LoadCue()
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return "", fmt.Errorf("Errors while loading Cue files:\n%w\n", errs)
	}

	// TODO, the rest, in parallel? Templates and all the file work
	// could probably go parallel across the generators

	// Load shadow, can this be done in parallel with the last step?
	shadow, err := gen.LoadShadow(verbose)
	if err != nil {
		return "", err
	}

	// Yes, we are sharing this here
	// TODO add a lock to the files eventually
	for _, G := range GS {
		G.Shadow = shadow

		err = G.GenerateFiles()
		if err != nil {
			return "", err
		}

	}

	// Finally, cleanup anything that remains in shadow


	for _, G := range GS {

		writestart := time.Now()

		for _, F := range G.Files {
			// Write the actual output
			if F.DoWrite {
				err = F.WriteOutput()
				if err != nil {
					errs = append(errs)
					continue
				}
			}

			// Write the shadow too, or if it doesn't exist
			if F.DoWrite || (F.IsSame > 0 && F.ShadowFile == nil) {
				err = F.WriteShadow()
				if err != nil {
					errs = append(errs)
					continue
				}
			}

			// remove from shadows map so we can cleanup what remains
			delete(G.Shadow, F.Filepath)
		}

		// Cleanup File & Shadow
		for f, _ := range G.Shadow {
			fmt.Println("Removing:", f)
			err := os.Remove(f)
			if err != nil {
				errs = append(errs)
				continue
			}
			err = os.Remove(path.Join(gen.SHADOW_DIR, f))
			if err != nil {
				errs = append(errs)
				continue
			}
			G.Stats.NumDeleted += 1
		}

		writeend := time.Now()
		G.Stats.WritingTime = writeend.Sub(writestart).Round(time.Millisecond)

		// Calc and print stats
		G.Stats.CalcTotals(G)
		fmt.Printf("\n%s\n==========================\n", G.Name)
		fmt.Println(G.Stats)

		// Print mrege issues
		for _, F := range G.Files {
			if F.IsConflicted > 0 {
				msg := fmt.Sprint("MERGE CONFLICT in:", F.Filepath)
				color.Red(msg)
			}
		}
	}

	// Print final timings
	veryend := time.Now()
	elapsed := veryend.Sub(verystart).Round(time.Millisecond)
	fmt.Printf("\nTotal Elapsed Time: %s\n\n", elapsed)

	if len(errs) > 0 {
		return "", fmt.Errorf("Errors while loading Cue files:\n%w\n", errs)
	}

	return "", nil
}


func extractGenerators(entrypoints []string) (gen.Generators, error) {
	GS := gen.Generators{}
	var RT cue.Runtime

	// TODO, config the second "config" arg here based on flags
	BIS := load.Instances(entrypoints, nil)
	for _, bi := range BIS {
		if bi.Err != nil {
			return GS, bi.Err
		}
		i, err := RT.Build(bi)
		if err != nil {
			return GS, err
		}

		// Get top level struct from cuelang
		toplevel, err := i.Value().Struct()
		if err != nil {
			return GS, err
		}

		// Loop through all top level fields
		iter := toplevel.Fields()
		for iter.Next() {

			label := iter.Label()
			value := iter.Value()

			if strings.HasPrefix(label, "HofGen") {
				G := gen.NewGenerator(label, value)
				GS[G.Name] = G
			}
		}
	}

	return GS, nil
}
