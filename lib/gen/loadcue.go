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
		if _, ok := file["Filepath"]; !ok {
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

		// TODO, better checking and/or decode directly into golang structs
		// But... they do all have defaults in the schema, so we will probably be OK
		// EXCEPT, spelling errors...

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

		// Meta information
		fn := file["Filepath"].(string)
		tc := file["Template"].(string)
		tn := file["TemplateName"].(string)
		ts := file["TemplateSystem"].(string)

		// deleimters
		ALT  := file["AltDelims"].(bool)
		SWAP := file["SwapDelims"].(bool)

		LHS2_D := file["LHS2_D"].(string)
		RHS2_D := file["RHS2_D"].(string)
		LHS3_D := file["LHS3_D"].(string)
		RHS3_D := file["RHS3_D"].(string)

		LHS2_S := file["LHS2_S"].(string)
		RHS2_S := file["RHS2_S"].(string)
		LHS3_S := file["LHS3_S"].(string)
		RHS3_S := file["RHS3_S"].(string)

		LHS2_T := file["LHS2_T"].(string)
		RHS2_T := file["RHS2_T"].(string)
		LHS3_T := file["LHS3_T"].(string)
		RHS3_T := file["RHS3_T"].(string)

		// Store the file in the generator
		F := &File {
			In: in,
			Filepath: fn,
			Template: tc,
			TemplateName: tn,
			TemplateSystem: ts,

			AltDelims: ALT,
			SwapDelims: SWAP,

			LHS2_D: LHS2_D,
			RHS2_D: RHS2_D,
			LHS3_D: LHS3_D,
			RHS3_D: RHS3_D,

			LHS2_S: LHS2_S,
			RHS2_S: RHS2_S,
			LHS3_S: LHS3_S,
			RHS3_S: RHS3_S,

			LHS2_T: LHS2_T,
			RHS2_T: RHS2_T,
			LHS3_T: LHS3_T,
			RHS3_T: RHS3_T,

		}

		G.Files[F.Filepath] = F

	}

	// TODO, should we erase the CueValue here so we release the memory?
	//       for now, yes we will
	G.CueValue = cue.Value{}

	return nil
}

