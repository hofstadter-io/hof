package gen

import (
	"fmt"
	"time"

	"cuelang.org/go/cue"
)

func (G *Generator) LoadCue() (errs []error) {
	// fmt.Println("Gen Load:", G.Name)
	start := time.Now()

	if err := G.loadOutdir(); err != nil {
		errs = append(errs, err)
	}

	if err := G.loadIn(); err != nil {
		errs = append(errs, err)
	}

	if err := G.loadTemplates(); err != nil {
		errs = append(errs, err)
	}

	if err := G.loadPartials(); err != nil {
		errs = append(errs, err)
	}

	if err := G.loadStatics(); err != nil {
		errs = append(errs, err)
	}

	if err := G.loadEmbeddedTemplates(); err != nil {
		errs = append(errs, err)
	}

	if err := G.loadEmbeddedPartials(); err != nil {
		errs = append(errs, err)
	}

	if err := G.loadEmbeddedStatics(); err != nil {
		errs = append(errs, err)
	}

	if err := G.loadOut(); err != nil {
		errs = append(errs, err)
	}

	if err := G.loadPackageName(); err != nil {
		errs = append(errs, err)
	}

	// Load Subgens

	// finalize load timing stats
	cueDecodeTime := time.Now()
	G.Stats.CueLoadingTime = cueDecodeTime.Sub(start)

	// Initialize Generator
	errsI := G.Initialize()
	if len(errsI) != 0 {
		fmt.Println("  Init Error:", errsI)
		errs = append(errs, errsI...)
	}

	G.debugLoad()

	return errs
}

func (G *Generator) debugLoad() {
	if !G.Debug {
		return
	}
	fmt.Println(G.Name, G.Outdir)
	fmt.Println("Out:    ", len(G.Out))
	fmt.Println("Tmpl:   ", len(G.Templates))
	fmt.Println("Prtl:   ", len(G.Partials))
	fmt.Println("Stcs:   ", len(G.Statics))
	fmt.Println("ETmpl:  ", len(G.EmbeddedTemplates))
	fmt.Println("EPrtl:  ", len(G.EmbeddedPartials))
	fmt.Println("EStcs:  ", len(G.EmbeddedStatics))
	fmt.Println()
	fmt.Println(G.PackageName, G.Disabled)
	fmt.Println("TMap:   ", len(G.TemplateMap))
	fmt.Println("PMap:   ", len(G.PartialsMap))
	fmt.Println("Files:  ", len(G.Files))
	fmt.Println("Shdw:   ", len(G.Shadow))
}

func (G *Generator) loadOutdir() error {
	val := G.CueValue.LookupPath(cue.ParsePath("Outdir"))
	if val.Err() != nil {
		return val.Err()
	}

	return val.Decode(&G.Outdir)
}

func (G *Generator) loadIn() error {
	val := G.CueValue.LookupPath(cue.ParsePath("In"))
	if val.Err() != nil {
		return val.Err()
	}

	G.In = make(map[string]interface{})
	return val.Decode(&G.In)
}

func (G *Generator) loadTemplates() error {
	val := G.CueValue.LookupPath(cue.ParsePath("Templates"))
	if val.Err() != nil {
		return val.Err()
	}

	G.Templates = make([]*TemplateGlobs, 0)
	return val.Decode(&G.Templates)
}

func (G *Generator) loadPartials() error {
	val := G.CueValue.LookupPath(cue.ParsePath("Partials"))
	if val.Err() != nil {
		return val.Err()
	}

	G.Partials = make([]*TemplateGlobs, 0)
	return val.Decode(&G.Partials)
}

func (G *Generator) loadStatics() error {
	val := G.CueValue.LookupPath(cue.ParsePath("Statics"))
	if val.Err() != nil {
		return val.Err()
	}

	G.Statics = make([]*StaticGlobs, 0)
	return val.Decode(&G.Statics)
}

func (G *Generator) loadEmbeddedTemplates() error {
	val := G.CueValue.LookupPath(cue.ParsePath("EmbeddedTemplates"))
	if val.Err() != nil {
		return val.Err()
	}

	G.EmbeddedTemplates = make(map[string]*TemplateContent)
	return val.Decode(&G.EmbeddedTemplates)
}

func (G *Generator) loadEmbeddedPartials() error {
	val := G.CueValue.LookupPath(cue.ParsePath("EmbeddedPartials"))
	if val.Err() != nil {
		return val.Err()
	}

	G.EmbeddedPartials = make(map[string]*TemplateContent)
	return val.Decode(&G.EmbeddedPartials)
}

func (G *Generator) loadEmbeddedStatics() error {
	val := G.CueValue.LookupPath(cue.ParsePath("EmbeddedStatics"))
	if val.Err() != nil {
		return val.Err()
	}

	G.EmbeddedStatics = make(map[string]string)
	return val.Decode(&G.EmbeddedStatics)
}

func (G *Generator) loadVal() error {

	return nil
}

func (G *Generator) loadOut() error {
	val := G.CueValue.LookupPath(cue.ParsePath("Out"))
	if val.Err() != nil {
		return val.Err()
	}

	Out := make([]*File, 0)
	err := val.Decode(&Out)
	if err != nil {
		return err
	}

	// need this extra work to load In into a cue.Value
	L, err := val.List()
	if err != nil {
		return err
	}

	G.Out = make([]*File, 0)
	i := 0
	for L.Next() {
		v := L.Value()
		in := v.LookupPath(cue.ParsePath("In"))

		// Only keep valid elements
		// Invalid include conditional elements in CUE Gen which are not "included"
		elem := Out[i]
		if elem != nil {
			// If In exists
			if in.Err() == nil {
				// merge with G.In
				for k, v := range G.In {
					// only copy in top-level elements which do not exist already
					if _, ok := elem.In[k]; !ok {
						elem.In[k] = v
					}
				}
			} else {
				// else, just use G.In
				elem.In = G.In
			}

			G.Out = append(G.Out, elem)
		}
		i++
	}

	return nil
}

func (G *Generator) loadPackageName() error {
	val := G.CueValue.LookupPath(cue.ParsePath("PackageName"))
	if val.Err() != nil {
		return val.Err()
	}

	return val.Decode(&G.PackageName)
}

func (G *Generator) decodeStuffOld() (errs []error) {

	// TODO, load subgenerators
	// Get the Generator Input (if it has one)
	/*
		Subgens, ok := gen["Generators"].(map[string]interface{})
		if ok && len(Subgens) > 0 {
			for sname, subgen := range Subgens {
				sg := NewGenerator(sname, cue.Value{})
				sgmap := subgen.(map[string]interface{})
				sgerrs := sg.decodeGenerator(sgmap)
				if len(sgerrs) > 0 {
					errs = append(errs, sgerrs...)
				}

				G.Generators[sname] = sg
			}
		}
	*/
	return errs
}
