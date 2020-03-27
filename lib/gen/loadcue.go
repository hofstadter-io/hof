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
	for _, O := range Out {
		file := O.(map[string]interface{})
		fn := file["Filename"].(string)
		tp := file["Template"].(string)
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

