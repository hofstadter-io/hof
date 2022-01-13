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

	if err := G.loadVal(); err != nil {
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

	if serr := G.loadOut(); serr != nil {
		errs = append(errs, serr...)
	}

	if err := G.loadPackageName(); err != nil {
		errs = append(errs, err)
	}

	// Initialize Generator
	errsI := G.Initialize()
	if len(errsI) != 0 {
		errs = append(errs, errsI...)
	}

	// finalize load timing stats
	end := time.Now()
	G.Stats.LoadingTime = end.Sub(start)

	// Load Subgens
	if serr := G.loadSubgens(); serr != nil {
		errs = append(errs, serr...)
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
	val := G.CueValue.LookupPath(cue.ParsePath("Val"))
	if val.Err() != nil {
		return val.Err()
	}

	G.Val = val

	return nil
}

func (G *Generator) loadOut() []error {
	val := G.CueValue.LookupPath(cue.ParsePath("Out"))
	if val.Err() != nil {
		return []error{val.Err()}
	}

	Out := make([]*File, 0)
	err := val.Decode(&Out)
	if err != nil {
		return []error{err}
	}

	// need this extra work to load In into a cue.Value
	L, err := val.List()
	if err != nil {
		return []error{err}
	}

	G.Out = make([]*File, 0)
	i := 0
	allErrs := []error{}
	for L.Next() {
		v := L.Value()
		elem := Out[i]

		err := G.loadFile(elem, v)
		if err != nil {
			allErrs = append(allErrs, err)
		}

		i++
	}

	if len(allErrs) > 0 {
		return allErrs
	}

	return nil
}

func (G *Generator) loadFile(file *File, val cue.Value) error {

	// Only keep valid elements
	// Invalid include conditional elements in CUE Gen which are not "included"
	if file != nil && file.Filepath != "" {

		tcE := file.TemplateContent == ""
		tpE := file.TemplatePath == ""
		dfE := file.DatafileFormat == ""

		// check template fields (See TODO in schema/gen/file.cue)
		// error if none are set
		if tcE && tpE && dfE {
			err := fmt.Errorf("In %s (%s), at least one of [TemplateContent, TemplatePath, DatafileFormat] must be set, all are empty", G.Name, file.Filepath)
			file.Errors = append(file.Errors, err)
			return err
		}
		// more than one is set
		if !(tcE || tpE) || !(tcE || dfE) || !(tpE || dfE) {
			err := fmt.Errorf("In %s (%s), only one of [TemplateContent, TemplatePath, DatafileFormat] may be set, multiple are", G.Name, file.Filepath)
			file.Errors = append(file.Errors, err)
			return err
		}
		// only one is set

		// If datafile format
		if !dfE {
			val := val.LookupPath(cue.ParsePath("Val"))
			if val.Err() == nil && val.Exists() {
				file.Value = val
			} else {
				file.Value = G.Val
			}
		} else {
			// TODO< check if a tc looks like a tp, or vice-a-versa
			// perhaps look for a path or template indicators

			in := val.LookupPath(cue.ParsePath("In"))
			// manage In value
			// If In exists
			if in.Err() == nil {
				// merge with G.In
				for k, v := range G.In {
					// only copy in top-level elements which do not exist already
					if _, ok := file.In[k]; !ok {
						file.In[k] = v
					}
				}
			} else {
				// else, just use G.In
				file.In = G.In
			}

		}

		G.Out = append(G.Out, file)
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

func (G *Generator) loadSubgens() (errs []error) {

	val := G.CueValue.LookupPath(cue.ParsePath("Generators"))
	if val.Err() != nil {
		return []error{val.Err()}
	}

	iter, err := val.Fields()
	if err != nil {
		return []error{err}
	}

	for iter.Next() {
		name := iter.Selector().String()
		v := iter.Value()
		sg := NewGenerator(name, v)

		sgerrs := sg.LoadCue()
		if len(sgerrs) > 0 {
			errs = append(errs, sgerrs...)
		}

		G.Generators[name] = sg
	}

	return errs
}
