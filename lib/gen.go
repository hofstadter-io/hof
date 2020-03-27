package lib

import (
	"fmt"
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
		G.Stats.CalcTotals(G)
		fmt.Printf("\n%s\n==========================\n", G.Name)
		fmt.Println(G.Stats)

		for _, F := range G.Files {
			if F.IsConflicted > 0 {
				msg := fmt.Sprint("MERGE CONFLICT in:", F.Filename)
				color.Red(msg)
			}
		}
	}
	veryend := time.Now()

	elapsed := veryend.Sub(verystart).Round(time.Millisecond)
	fmt.Printf("\nTotal Elapsed Time: %s\n\n", elapsed)

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

			if strings.HasPrefix(label, "Gen") {
				G := gen.NewGenerator(label, value)
				GS[G.Name] = G
			}
		}
	}

	return GS, nil
}
