package gen

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/errors"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/hofstadter-io/hof/lib/templates"
)

func getLang() language.Tag {
	loc := os.Getenv("LC_ALL")
	if loc == "" {
		loc = os.Getenv("LANG")
	}
	loc = strings.Split(loc, ".")[0]
	return language.Make(loc)
}

func (G *Generator) LoadCue() ([]error) {
	// fmt.Println("Gen Load:", G.Name)

	var gen map[string]interface{}
	start := time.Now()

	// Decode the value into a temporary "generator" with timing
	err := G.CueValue.Decode(&gen)
	if err != nil {

		p := message.NewPrinter(getLang())
		format := func(w io.Writer, format string, args ...interface{}) {
			p.Fprintf(w, format, args...)
		}
		cwd, _ := os.Getwd()
		w := &bytes.Buffer{}
		errors.Print(w, err, &errors.Config{
			Format:  format,
			Cwd:     cwd,
			ToSlash: false,
		})
		s := w.String()
		fmt.Println(s)

		return []error{err}
	}

	// finalize load timing stats
	cueDecodeTime := time.Now()
	G.Stats.CueLoadingTime = cueDecodeTime.Sub(start)

	return G.decodeGenerator(gen)
}

func (G *Generator) decodeGenerator(gen map[string]interface{}) ([]error) {
	var errs []error

	// Get Out, or the files we want to render, required
	Out, ok := gen["Out"].([]interface{})
	if !ok {
		return []error{fmt.Errorf("Generator: %q is missing 'Out' field.", G.Name)}
	}

	// Get the Generator Input (if it has one)
	In, ok := gen["In"].(map[string]interface{})
	if ok {
		G.In = In
	}

	G.Outdir = gen["Outdir"].(string)

	//
	// From common
	//

	// deleimters
	configI, ok := gen["TemplateConfig"]
	if ok {
		config := configI.(map[string]interface{})
		G.TemplateConfig = &templates.Config{}

		G.TemplateConfig.TemplateSystem = config["TemplateSystem"].(string)
		G.TemplateConfig.AltDelims  = config["AltDelims"].(bool)
		G.TemplateConfig.SwapDelims = config["SwapDelims"].(bool)

		G.TemplateConfig.LHS2_D = config["LHS2_D"].(string)
		G.TemplateConfig.RHS2_D = config["RHS2_D"].(string)
		G.TemplateConfig.LHS3_D = config["LHS3_D"].(string)
		G.TemplateConfig.RHS3_D = config["RHS3_D"].(string)

		G.TemplateConfig.LHS2_S = config["LHS2_S"].(string)
		G.TemplateConfig.RHS2_S = config["RHS2_S"].(string)
		G.TemplateConfig.LHS3_S = config["LHS3_S"].(string)
		G.TemplateConfig.RHS3_S = config["RHS3_S"].(string)

		G.TemplateConfig.LHS2_T = config["LHS2_T"].(string)
		G.TemplateConfig.RHS2_T = config["RHS2_T"].(string)
		G.TemplateConfig.LHS3_T = config["LHS3_T"].(string)
		G.TemplateConfig.RHS3_T = config["RHS3_T"].(string)
	}

	G.PackageName, _  = gen["PackageName"].(string)

	// In cue code
	G.NamedTemplates = make(map[string]string)
	nt, ok := gen["NamedTemplates"].(map[string]interface{})
	if !ok {
		return []error{fmt.Errorf("Generator: %q field 'NamedTemplates' is not an object.", G.Name)}
	}
	for k, t := range nt {
		G.NamedTemplates[k] = t.(string)
	}

	G.NamedPartials = make(map[string]string)
	np, ok := gen["NamedPartials"].(map[string]interface{})
	if !ok {
		return []error{fmt.Errorf("Generator: %q field 'NamedParitals' is not an object.", G.Name)}
	}
	for k, p := range np {
		G.NamedPartials[k] = p.(string)
	}

	G.StaticFiles = make(map[string]string)
	sf, ok := gen["StaticFiles"].(map[string]interface{})
	if !ok {
		return []error{fmt.Errorf("Generator: %q field 'StaticFiles' is not an object.", G.Name)}
	}
	for k, s := range sf {
		G.StaticFiles[k] = s.(string)
	}

	// Eventually loaded from disk
	G.TemplatesDir = gen["TemplatesDir"].(string)
	G.PartialsDir  = gen["PartialsDir"].(string)
	G.StaticGlobs = make([]string, 0)
	sg, ok := gen["StaticGlobs"].([]interface{})
	if !ok {
		return []error{fmt.Errorf("Generator: %q field 'StaticGlobs' is not a list.", G.Name)}
	}
	for _, s := range sg {
		G.StaticGlobs = append(G.StaticGlobs, s.(string))
	}

	// TODO, load subgenerators

	// Decode generator files
	// Turn G.Out elements into G.Files
	for i, O := range Out {
		file := O.(map[string]interface{})

		F, err := G.decodeFile(i, file)
		if err != nil {
			errs = append(errs, err)
		}

		G.Files[F.Filepath] = F

	}

	// TODO, should we erase the CueValue here so we release the memory?
	//       for now, yes we will
	G.CueValue = cue.Value{}

	return errs
}

func (G *Generator) decodeFile(i int, file map[string]interface{}) (*File, error) {

	// Is this output missing a filename? then skip it
	if _, ok := file["Filepath"]; !ok {
		mockname := fmt.Sprintf("noname-%d", i)
		F := &File {
			FileStats: FileStats{
				IsSkipped: 1,
			},
			FinalContent: []byte(mockname),
		}
		// We skip files this way, probably want to continue to do that as convention
		return F, nil
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

		// deleimters
	configI, ok := file["TemplateConfig"]
	if ok {
		config := configI.(map[string]interface{})
		F.TemplateConfig = &templates.Config{}

		F.TemplateConfig.TemplateSystem = config["TemplateSystem"].(string)
		F.TemplateConfig.AltDelims  = config["AltDelims"].(bool)
		F.TemplateConfig.SwapDelims = config["SwapDelims"].(bool)

		F.TemplateConfig.LHS2_D = config["LHS2_D"].(string)
		F.TemplateConfig.RHS2_D = config["RHS2_D"].(string)
		F.TemplateConfig.LHS3_D = config["LHS3_D"].(string)
		F.TemplateConfig.RHS3_D = config["RHS3_D"].(string)

		F.TemplateConfig.LHS2_S = config["LHS2_S"].(string)
		F.TemplateConfig.RHS2_S = config["RHS2_S"].(string)
		F.TemplateConfig.LHS3_S = config["LHS3_S"].(string)
		F.TemplateConfig.RHS3_S = config["RHS3_S"].(string)

		F.TemplateConfig.LHS2_T = config["LHS2_T"].(string)
		F.TemplateConfig.RHS2_T = config["RHS2_T"].(string)
		F.TemplateConfig.LHS3_T = config["LHS3_T"].(string)
		F.TemplateConfig.RHS3_T = config["RHS3_T"].(string)
	}

	return F, nil
}
