package gen

import (
	"fmt"
	"path/filepath"
	"time"

	"cuelang.org/go/cue"

	"github.com/hofstadter-io/hof/lib/templates"
)

const CUE_VENDOR_DIR = "./cue.mod/pkg/"

type TemplateGlobs struct {
	// Globs to load
	Globs []string
	// Prefix to trim
	TrimPrefix string
	// Custom delims
	Delims *templates.Delims
}

type StaticGlobs struct {
	// Globs to load
	Globs []string
	// Prefix to trim
	TrimPrefix string
	// Prefix to add before output
	OutPrefix string
}

type TemplateContent struct {
	Content string
	Delims  *templates.Delims
}

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
	Out []*File

	//
	// Generator configuration set in Cue code
	//

	Templates []*TemplateGlobs
	Partials  []*TemplateGlobs

	// Filepath globs for static files to load
	Statics []*StaticGlobs

	// The following will be automatically added to the template context
	// under its name for reference in GenFiles  and partials in templates
	EmbeddedTemplates map[string]*TemplateContent
	EmbeddedPartials  map[string]*TemplateContent

	// Static files are available for pure cue generators that want to have static files
	// These should be named by their filepath, but be the content of the file
	EmbeddedStatics map[string]string

	// Subgenerators for composition
	Generators map[string]*Generator

	// Used for indexing into the vendor directory...
	PackageName string

	//
	// Hof internal usage
	//

	// Disabled? we do this when looking at expressions and optimizing
	// TODO, make this field available in cuelang?
	Disabled bool

	// Template System Cache
	TemplateMap templates.TemplateMap
	PartialsMap templates.TemplateMap

	// Files and the shadow dir for doing neat things
	Files  map[string]*File
	Shadow map[string]*File

	// Print extra information
	Debug bool

	// Status for this generator and processing
	Stats *GeneratorStats

	// Cuelang related, also set externally
	CueValue cue.Value
}

func NewGenerator(label string, value cue.Value) *Generator {
	return &Generator{
		Name:        label,
		CueValue:    value,
		PartialsMap: templates.NewTemplateMap(),
		TemplateMap: templates.NewTemplateMap(),
		Generators:  make(map[string]*Generator),
		Files:       make(map[string]*File),
		Shadow:      make(map[string]*File),
		Stats:       &GeneratorStats{},
	}
}

func (G *Generator) GenerateFiles() []error {
	errs := []error{}

	start := time.Now()

	for _, F := range G.Files {

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

func (G *Generator) Initialize() []error {
	var errs []error

	// First do partials, so available to all templates
	errs = G.initPartials()
	if len(errs) > 0 {
		return errs
	}

	// Then do templates, will be needed for files
	errs = G.initTemplates()
	if len(errs) > 0 {
		return errs
	}

	// Then do files, we should be ready to gen/write now
	errs = G.initFileGens()
	if len(errs) > 0 {
		return errs
	}

	return errs
}

func (G *Generator) initPartials() []error {
	var errs []error

	// First named
	for path, tc := range G.EmbeddedPartials {
		T, err := templates.CreateFromString(path, tc.Content, tc.Delims)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		G.PartialsMap[path] = T
	}

	for _, tg := range G.Partials {
		for _, glob := range tg.Globs {
			// setup vars
			prefix := filepath.Clean(tg.TrimPrefix)
			if G.PackageName != "" {
				glob = filepath.Join(CUE_VENDOR_DIR, G.PackageName, glob)
				prefix = filepath.Join(CUE_VENDOR_DIR, G.PackageName, prefix)
			}

			pMap, err := templates.CreateTemplateMapFromFolder(glob, prefix, tg.Delims)
			if err != nil {
				errs = append(errs, err)
				continue
			}

			for k, T := range pMap {
				_, ok := G.PartialsMap[k]
				if !ok {
					// TODO, do we also want to namespace with the template module name?
					G.PartialsMap[k] = T
				} else {
					errs = append(errs, fmt.Errorf("duplicate partial %q", k))
				}
			}
		}
	}

	// register all partials with partials
	for _, P := range G.PartialsMap {
		G.registerPartials(P)
	}
	return errs
}

func (G *Generator) initTemplates() []error {
	var errs []error

	// First named
	for path, tc := range G.EmbeddedTemplates {
		T, err := templates.CreateFromString(path, tc.Content, tc.Delims)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		G.TemplateMap[path] = T
	}

	for _, tg := range G.Templates {
		for _, glob := range tg.Globs {
			// setup vars
			prefix := filepath.Clean(tg.TrimPrefix)
			if G.PackageName != "" {
				glob = filepath.Join(CUE_VENDOR_DIR, G.PackageName, glob)
				prefix = filepath.Join(CUE_VENDOR_DIR, G.PackageName, prefix)
			}

			pMap, err := templates.CreateTemplateMapFromFolder(glob, prefix, tg.Delims)
			if err != nil {
				errs = append(errs, err)
				continue
			}

			for k, T := range pMap {
				_, ok := G.TemplateMap[k]
				if !ok {
					// TODO, do we also want to namespace with the template module name?
					G.TemplateMap[k] = T
				} else {
					errs = append(errs, fmt.Errorf("duplicate partial %q", k))
				}
			}
		}
	}

	// Register partials with all templates
	for _, T := range G.TemplateMap {
		G.registerPartials(T)
	}

	return errs
}

func (G *Generator) initFileGens() []error {
	var errs []error

	for _, F := range G.Out {
		G.Files[F.Filepath] = F
	}

	for _, F := range G.Files {
		err := G.ResolveFile(F)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func (G *Generator) ResolveFile(F *File) error {

	// Inline template content
	if F.TemplateContent != "" {

		T, err := templates.CreateFromString(F.Filepath /* or "inline"? */, F.TemplateContent, F.TemplateDelims)
		if err != nil {
			return err
		}

		// Now register partials with all templates
		G.registerPartials(T)

		F.TemplateInstance = T
	}

	// Template is embedded or loaded from FS
	if F.TemplatePath != "" {
		T, ok := G.TemplateMap[F.TemplatePath]
		if !ok {
			// TODO, do we need to do check for a namespaced prefix?
			err := fmt.Errorf("Named template %q not found for %s %s", F.TemplatePath, G.Name, F.Filepath)
			F.IsErr = 1
			F.Errors = append(F.Errors, err)
			return err
		}

		F.TemplateInstance = T
	}

	return nil
}

func (G *Generator) registerPartials(T *templates.Template) error {
	if T.T == nil {
		return fmt.Errorf("T template is not initialized %q", T.Name)
	}

	for k, P := range G.PartialsMap {
		t := T.T.New(k)

		// todo, do we need to do this twice?, has it already been done?
		// maybe? because of how text/template contexts work
		templates.AddGolangHelpers(t)
		t.Parse(P.Source)

		// T.T = t
	}

	return nil
}
