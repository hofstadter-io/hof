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
	cueDecodeTime := time.Now()
	G.Stats.CueLoadingTime = cueDecodeTime.Sub(start)

	return G.decodeGenerator(gen)
}

func (G *Generator) decodeGenerator(gen map[string]interface{}) (error) {

	// Get Out, or the files we want to render, required
	Out, ok := gen["Out"].([]interface{})
	if !ok {
		return fmt.Errorf("Generator: %q is missing 'Out' field.", G.Name)
	}

	// Get the Generator Input (if it has one)
	In, ok := gen["In"].(map[string]interface{})
	if ok {
		G.In = In
	}

	G.PackageName, _  = gen["PackageName"].(string)

	// In cue code
	G.NamedTemplates, _ = gen["NamedTemplates"].(map[string]string)
	G.NamedPartials,  _ = gen["NamedPartials"].(map[string]string)
	G.StaticFiles,    _ = gen["StaticFiles"].(map[string]string)

	// Eventually loaded from disk
	G.TemplatesDir, _ = gen["TemplatesDir"].(string)
	G.PartialsDir, _  = gen["PartialsDir"].(string)
	G.StaticGlobs, _  = gen["StaticGlobs"].([]string)

	// TODO, load subgenerators

	// Decode generator files
	// Turn G.Out elements into G.Files
	for i, O := range Out {
		file := O.(map[string]interface{})

		F := G.decodeFile(i, file)

		G.Files[F.Filepath] = F

	}

	// TODO, should we erase the CueValue here so we release the memory?
	//       for now, yes we will
	G.CueValue = cue.Value{}

	return nil
}

func (G *Generator) decodeFile(i int, file map[string]interface{}) *File {

	// Is this output missing a filename? then skip it
	if _, ok := file["Filepath"]; !ok {
		mockname := fmt.Sprintf("noname-%d", i)
		F := &File {
			FileStats: FileStats{
				IsSkipped: 1,
			},
			FinalContent: []byte(mockname),
		}
		return F

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

	F := &File {
		In: in,
	}

	// Meta information
	F.Filepath = file["Filepath"].(string)
	F.Template = file["Template"].(string)
	F.TemplateName = file["TemplateName"].(string)
	F.TemplateSystem = file["TemplateSystem"].(string)

	// deleimters
	F.AltDelims  = file["AltDelims"].(bool)
	F.SwapDelims = file["SwapDelims"].(bool)

	F.LHS2_D = file["LHS2_D"].(string)
	F.RHS2_D = file["RHS2_D"].(string)
	F.LHS3_D = file["LHS3_D"].(string)
	F.RHS3_D = file["RHS3_D"].(string)

	F.LHS2_S = file["LHS2_S"].(string)
	F.RHS2_S = file["RHS2_S"].(string)
	F.LHS3_S = file["LHS3_S"].(string)
	F.RHS3_S = file["RHS3_S"].(string)

	F.LHS2_T = file["LHS2_T"].(string)
	F.RHS2_T = file["RHS2_T"].(string)
	F.LHS3_T = file["LHS3_T"].(string)
	F.RHS3_T = file["RHS3_T"].(string)

	return F
}
