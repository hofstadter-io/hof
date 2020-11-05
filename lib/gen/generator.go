package gen

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"cuelang.org/go/cue"

	"github.com/hofstadter-io/hof/lib/templates"
)


// A generator pulled from the cue instances
type Generator struct {
	//
	// Set by Hof via cuelang extraction
	// Label in Cuelang
	Name string

	// Base directory for output
	Outdir string

  // "Global" input, merged with out replacing onto the files
	In map[string]interface{}

  // The list fo files for hof to generate, in cue values
	Out []map[string]interface{}

	//
	// Generator configuration set in Cue code
	//

  // Subgenerators for composition
  Generators []*Generator

  // Template delimiters
	TemplateConfig *templates.Config

  // The following will be automatically added to the template context
  // under its name for reference in GenFiles  and partials in templates
  NamedTemplates map[string]string
  NamedPartials  map[string]string

  // Static files are available for pure cue generators that want to have static files
  // These should be named by their filepath, but be the content of the file
  StaticFiles map[string]string

  //
  // For file based generators
  //   files here will be automatically added to the template context
  //   under its filepath for reference in GenFiles and partials in templates

  // Used for indexing into the vendor directory...
  PackageName string

  // Base directory of entrypoint templates to load
  TemplatesDir string
	TemplatesDirConfig map[string]*templates.Config

  // Base directory of partial templatess to load
  PartialsDir string
	PartialsDirConfig map[string]*templates.Config

  // Filepath globs for static files to load
  StaticGlobs []string


	//
	// Hof internal usage
	//

	// Disabled? we do this when looking at expressions and optimizing
	// TODO, make this field available in cuelang?
	Disabled bool

	// Template System Cache
	PartialsMap templates.TemplateMap
	TemplateMap templates.TemplateMap

	// Files and the shadow dir for doing neat things
	Files map[string]*File
	Shadow map[string]*File

	// Status for this generator and processing
	Stats *GeneratorStats

	// Cuelang related, also set externally
	CueValue         cue.Value
}

func NewGenerator(label string, value cue.Value) *Generator{
	return &Generator {
		Name: label,
		CueValue: value,
		PartialsMap: templates.NewMap(),
		TemplateMap: templates.NewMap(),
		Files: make(map[string]*File),
		Shadow: make(map[string]*File),
		Stats: &GeneratorStats{},
	}
}

func (G *Generator) GenerateFiles() []error {
	errs := []error{}

	start := time.Now()

	for _, F := range G.Files {

		// fmt.Printf("GenerateFile: %s\n%#+v\n\n", F.Filepath, F)
		if F.Filepath == "" {
			F.IsSkipped = 1
			continue
		}
		shadowFN := filepath.Join(G.Name, F.Filepath)
		F.ShadowFile = G.Shadow[shadowFN]
		err := F.Render(filepath.Join(SHADOW_DIR, shadowFN))
		if err != nil {
			F.IsErr = 1
			F.Errors = append(F.Errors, err)
			errs = append(errs, err)
			continue
		}
	}

	elapsed := time.Now().Sub(start).Round(time.Millisecond)
	G.Stats.RenderingTime = elapsed

	return errs
}


func (G *Generator) Initialize() ([]error) {
	var errs []error
	// fmt.Println("Intitializing Generator: ", G.Name)

	// First do partials, so available to all templates
	errs = G.initPartials()
	if len(errs) > 0 {
		// fmt.Printf("initPartials Errors:\n%v\n", errs)
		return errs
	}
	// fmt.Println("  Partials:", G.PartialsMap )

	errs = G.initTemplates()
	if len(errs) > 0 {
		// fmt.Printf("initTemplates Errors:\n%v\n", errs)
		return errs
	}
	// fmt.Println("  Templates:", G.TemplateMap )

	errs = G.initFileGens()
	if len(errs) > 0 {
		// fmt.Printf("initFileGens Errors:\n%v\n", errs)
		return errs
	}

	// fmt.Println("Intitialized Generator: ", G.Name)
	// fmt.Printf("%# v\n", pretty.Formatter(G))

	return errs
}

const CUE_VENDOR_DIR = "./cue.mod/pkg/"

func (G *Generator) initPartials() []error {
	var errs []error

	// First named
	for k, content := range G.NamedPartials {
		T, err := templates.CreateFromString(k, content, G.TemplateConfig.TemplateSystem, G.TemplateConfig)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		G.PartialsMap[k] = T
	}

	// Then file based partials, but don't overwrite
	pDir := G.PartialsDir
	if G.PackageName != "" {
		pDir = filepath.Join(CUE_VENDOR_DIR, G.PackageName, G.PartialsDir)
	}
	pMap, err := templates.CreateTemplateMapFromFolder(pDir, G.TemplateConfig.TemplateSystem, G.TemplateConfig, G.PartialsDirConfig)
	if err != nil {
		return append(errs, err)
	}
	// fmt.Println("pFileMap", pDir, pMap)

	for k, T := range pMap {
		if strings.HasPrefix(k, "/") {
			k = k[1:]
		}
		_, ok := G.PartialsMap[k]
		if !ok {
			G.PartialsMap[k] = T
		}

		// add second copy without the partials prefix
		// seems to be an edge case when using a generator from within it's own directory
		if strings.HasPrefix(k, "partials/") {
			k = strings.TrimPrefix(k, "partials/")
			G.PartialsMap[k] = T
		}
	}

	return errs
}

func (G *Generator) initTemplates() []error {
	var errs []error

	// First named
	for k, content := range G.NamedTemplates {
		T, err := templates.CreateFromString(k, content, G.TemplateConfig.TemplateSystem, G.TemplateConfig)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		G.TemplateMap[k] = T
	}

	// Then file based template, but don't overwrite
	tDir := G.TemplatesDir
	if G.PackageName != "" {
		tDir = filepath.Join(CUE_VENDOR_DIR, G.PackageName, G.TemplatesDir)
	}
	tMap, err := templates.CreateTemplateMapFromFolder(tDir, G.TemplateConfig.TemplateSystem, G.TemplateConfig, G.TemplatesDirConfig)
	if err != nil {
		return append(errs, err)
	}

	for k, T := range tMap {
		if strings.HasPrefix(k, "/") {
			k = k[1:]
		}
		_, ok := G.TemplateMap[k]
		if !ok {
			G.TemplateMap[k] = T
		}
	}

	// Now register partials with all patrials and templates
	for _, P := range G.PartialsMap {
		G.registerPartials(P)
	}
	for _, T := range G.TemplateMap {
		G.registerPartials(T)
	}

	return errs
}

func (G *Generator) initFileGens() []error {
	var errs []error

	for _, F := range G.Files {
		err := G.ResolveFile(F)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func (G *Generator) ResolveFile(F *File) error {

	// Override delims
	if F.TemplateConfig == nil {
		// Just use gen's if nil
		F.TemplateConfig = G.TemplateConfig
	} else {
		// Override "default" '.'
		F.TemplateConfig.OverrideDotDefaults(G.TemplateConfig)
	}

	// both valued?
	if F.Template != "" && F.TemplateName != "" {
		err := fmt.Errorf("Cannot specify both Template and TemplateName in Gen: %q File: %q TName: %q\n", G.Name, F.Filepath, F.TemplateName)
		F.IsErr = 1
		F.Errors = append(F.Errors, err)
		return err
	}
	// both emtpy?
	if F.Template == "" && F.TemplateName == "" {
		err := fmt.Errorf("Must specify one of Template and TemplateName in Gen: %q File: %q TName: %q\n", G.Name, F.Filepath, F.TemplateName)
		F.IsErr = 1
		F.Errors = append(F.Errors, err)
		return err
	}

	// Named or File Template
	if F.Template == "" && F.TemplateName != "" {
		T, ok := G.TemplateMap[F.TemplateName]
		if !ok {
			// Try adding the generators template dir as a prefix when PackageName is empty
			if G.PackageName == "" {
				T, ok = G.TemplateMap[filepath.Join(G.TemplatesDir, F.TemplateName)]
			}

			// check if we have not found the template
			if !ok {
				err := fmt.Errorf("Named template %q not found for %s %s\n", F.TemplateName, G.Name, F.Filepath)
				F.IsErr = 1
				F.Errors = append(F.Errors, err)
				return err
			}
		}

		F.TemplateInstance = T
		return nil
	}

	if F.Template != "" {

		T, err := templates.CreateFromString(F.Filepath, F.Template, F.TemplateConfig.TemplateSystem, F.TemplateConfig)
		if err != nil {
			return err
		}

		// Now register partials with all templates
		G.registerPartials(T)

		F.TemplateInstance = T
	}

	// fmt.Println("    TI:", F.TemplateInstance)

	return nil
}

func (G *Generator) registerPartials(T *templates.Template) {
	if T.R != nil {
		for k, P := range G.PartialsMap {
			if T.Config.TemplateSystem == P.Config.TemplateSystem {
				T.R.RegisterPartialTemplate(k, P.R)
			}
		}
	}

	if T.T != nil {
		for k, P := range G.PartialsMap {
			// fmt.Println("Partial - Golang -", k)
			if T.Config.TemplateSystem == P.Config.TemplateSystem {
				t := T.T.New(k)
				// TODO Delims again here?
				templates.AddGolangHelpers(t)

				t.Parse(P.Source)

				// T.T = t
			}
		}
	}

}
