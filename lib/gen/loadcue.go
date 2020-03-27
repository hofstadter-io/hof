package gen

import (
	"fmt"
	"time"

	"cuelang.org/go/cue"
)

func (G *Generator) LoadCue() (error) {
	// fmt.Println("Gen Load:", G.Name)

	var gen map[string]interface{}
	start := time.Now()

	// Decode the value into a temporary "generator" with timing
	err := G.CueValue.Decode(&gen)
	if err != nil {
		return err
	}

	// finalize load timing stats
	decodeTime := time.Now()
	G.Stats.CueLoadingTime = decodeTime.Sub(start)

	// Get the Generator Input (if it has one)
	In, ok := gen["In"].(map[string]interface{})
	if ok {
		G.In = In
	}

	// Get Out, or the files we want to render, required
	Out, ok := gen["Out"].([]interface{})
	if !ok {
		return fmt.Errorf("Generator: %q is missing 'Out' field.", G.Name)
	}

	// Turn G.Out elements into G.Files
	for i, O := range Out {
		file := O.(map[string]interface{})

		// Is this output missing a filename? then skip it
		if _, ok := file["Filename"]; !ok {
			mockname := fmt.Sprintf("noname-%d", i)
			F := &File {
				FileStats: FileStats{
					IsSkipped: 1,
				},
				FinalContent: []byte(mockname),
			}

			G.Files[mockname] = F
			continue
		}

		// Otherwise, we want to do something with this file
		fn := file["Filename"].(string)
		tp := file["Template"].(string)

		// Build up the files "In" value
		in, ok := file["In"].(map[string]interface{})
		if !ok {
			in = G.In
		} else {
			// Else, 'IN' has key and 'in' does not, add it
			for key, val := range G.In {
				if _, ok := in[key]; !ok {
					// fmt.Println("checking In, filling", key)
					in[key] = val
				}
			}
		}

		// Store the file in the generator
		F := &File {
			Filename: fn,
			Template: tp,
			In: in,
		}

		G.Files[F.Filename] = F

	}

	// TODO, should we erase the CueValue here so we release the memory?
	//       for now, yes we will
	G.CueValue = cue.Value{}

	return nil
}

