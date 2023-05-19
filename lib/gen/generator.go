package gen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cuelang.org/go/cue"
	"github.com/mattn/go-zglob"

	"github.com/hofstadter-io/hof/lib/hof"
	"github.com/hofstadter-io/hof/lib/templates"
)

const CUE_VENDOR_DIR = "cue.mod/pkg"

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
	*hof.Node[Generator]
	//
	// Set by Hof via cuelang extraction
	// Label in Cuelang
	Name string

	// Base directory for output
	Outdir string

	// Other important dirs when loading templates (auto set)
	CueModuleRoot string
	WorkingDir    string
	CwdToRoot     string  // module root <- working dir (../..)

	// "Global" input, merged with out replacing onto the files
	In  map[string]any
	Val cue.Value

	// File globs to watch and trigger regen on change
	WatchFull []string
	WatchFast  []string

	// Formatting
	FormattingDisabled bool
	FormatData         bool
	FormattingConfigs  map[string]FmtConfig

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

	// backpointers, if a subgen
	parent  *Generator

	// Used for indexing into the vendor directory...
	PackageName string

	// Use Diff3 & Shadow
	Diff3FlagSet bool // set by flag
	UseDiff3 bool
	NoFormat bool

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
	OrderedFiles    []*File
	Files  map[string]*File
	Shadow map[string]*File

	// Print extra information
	Debug bool
	Verbosity int

	// Status for this generator and processing
	Stats *GeneratorStats

	// Cuelang related, also set externally
	CueValue cue.Value
}

func NewGenerator(node *hof.Node[Generator]) *Generator {
	// TODO, only transfer what is needed

	return &Generator{
		Node: node,

		// generator specific vals
		Name:          node.Hof.Label,
		CueValue:      node.Value,

		// initialize containers
		PartialsMap:   templates.NewTemplateMap(),
		TemplateMap:   templates.NewTemplateMap(),
		Generators:    make(map[string]*Generator),
		Files:         make(map[string]*File),
		Shadow:        make(map[string]*File),
		Stats:         &GeneratorStats{},
	}
}

// Returns Generators name path, including parents
// as a path like string
func (G *Generator) NamePath() string {
	p := G.Name
	if G.parent != nil {
		p = filepath.Join(G.parent.NamePath(), p)
	}
	return p
}

// Returns Generators contribution to the output path,
// including parents contributions if a subgen.
// Each gen in the path is [parent]/G.Outdir
func (G *Generator) OutputPath() string {
	p := G.Outdir
	if G.parent != nil {
		p = filepath.Join(G.parent.OutputPath(), p)
	}
	return p
}

// Returns Generators contribution to the shadow path,
// including parents contributions if a subgen.
// Each gen in the path is [parent]/G.Name/G.Outdir
func (G *Generator) ShadowPath() string {
	p := filepath.Join(G.Name, G.Outdir)
	if G.parent != nil {
		p = filepath.Join(G.parent.ShadowPath(), p)
	}
	return p
}

func (G *Generator) Initialize() []error {
	var errs []error
	if G.Verbosity > 1 {
		fmt.Println("initializing:", G.NamePath())
	}

	// zero, read static files
	errs = G.initStaticFiles()
	if len(errs) > 0 {
		return errs
	}

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

/* TODO, that the order of embedded vs disk files is inconsistent, we should clean this up and ensure consistent semantics (which may be the case?)
	- statics: disk -> embed
	- partials: embed -> disk
	- templates: embed -> disk
*/

func (G *Generator) initStaticFiles() []error {
	var errs []error

	bdir := G.CueModuleRoot
	// lookup in vendor directory, this will need to change once CUE uses a shared cache in the user homedir
	if G.PackageName != "" {
		bdir = filepath.Join(G.CueModuleRoot, CUE_VENDOR_DIR, G.PackageName)
	}

	// Start with static file globs
	for _, Static := range G.Statics {
		// we need to check if the base directory exists, becuase we have defaults in the schema
		_, err := os.Stat(filepath.Join(bdir, Static.TrimPrefix))
		if err != nil {
			fmt.Printf("warning: %s not found for %s, if you do not intend to use this directory, set 'Statics: []'\n", Static.TrimPrefix, G.PackageName)
			continue
		}

		for _, Glob := range Static.Globs {

			// get list of static files
			matches, err := zglob.Glob(filepath.Join(bdir, Glob))
			if err != nil {
				err = fmt.Errorf("while globbing %s / %s\n%w\n", bdir, Glob, err)
				errs = append(errs, err)
				return errs
			}
			if G.Verbosity > 1 {
				fmt.Printf("%s:%s:%s has %d static matches\n", G.NamePath(), bdir, Glob, len(matches))
			}

			// for each static file, calc some dirs and write output & shadow
			for _, match := range matches {
				info, err := os.Stat(match)
				if err != nil {
					fmt.Printf("warning: error while loading statics %s: %s\n", match, err)
					continue
				}
				if info.IsDir() {
					continue
				}
				// read the file
				content, err := os.ReadFile(match)
				if err != nil {
					errs = append(errs, err)
					continue
				}

				// remove and add prefixes, per the configuration
				mo := strings.TrimPrefix(match, filepath.Join(bdir, Static.TrimPrefix))
				// because Join removes?
				mo = strings.TrimPrefix(mo, "/")
				fp := filepath.Join(Static.OutPrefix, mo)

				if G.Verbosity > 2 {
					fmt.Println("static FN:", match, filepath.Join(bdir, Static.TrimPrefix), mo)
					fmt.Println("    ", fp, filepath.Clean(fp))
				}

				// create a file
				F := &File{
					Filepath:     filepath.Clean(fp),
					RenderContent: []byte(content),
					StaticFile:   true,
				}

				// check for collisions
				if _,ok := G.Files[F.Filepath]; ok {
					errs = append(errs, fmt.Errorf("duplicate static file %q in %q", F.Filepath, G.NamePath()))
					continue
				}

				if G.Verbosity > 1 {
					fmt.Printf(" +s %s:%s\n", G.NamePath(), F.Filepath)
				}

				G.Files[F.Filepath] = F
				G.OrderedFiles = append(G.OrderedFiles, F)
			}
		}
	}

	// Then the static files in cue
	for p, content := range G.EmbeddedStatics {
		F := &File{
			Filepath:     filepath.Clean(p),
			RenderContent: []byte(content),
			StaticFile:   true,
		}

		// check for collisions
		if _,ok := G.Files[F.Filepath]; ok {
			errs = append(errs, fmt.Errorf("duplicate static file %q in %q", F.Filepath, G.NamePath()))
			continue
		}

		if G.Verbosity > 1 {
			fmt.Printf(" +s %s:%s\n", G.NamePath(), F.Filepath)
		}

		G.Files[F.Filepath] = F
		G.OrderedFiles = append(G.OrderedFiles, F)
	}


	return errs
}

func (G *Generator) initPartials() []error {
	var errs []error

	// First named / embedded partials
	for path, tc := range G.EmbeddedPartials {
		T, err := templates.CreateFromString(path, tc.Content, tc.Delims)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		// check for collisions
		_, ok := G.PartialsMap[path]
		if !ok {
			if G.Verbosity > 1 {
				fmt.Printf(" +p %s:%s\n", G.NamePath(), path)
			}
			// TODO, do we also want to namespace with the template module name?
			G.PartialsMap[path] = T
		} else {
			errs = append(errs, fmt.Errorf("duplicate partial %s:%s", G.NamePath(), path))
		}
	}

	// then partials from disk via globs
	for _, tg := range G.Partials {
		prefix := filepath.Clean(tg.TrimPrefix)
		if G.PackageName != "" {
			prefix = filepath.Join(CUE_VENDOR_DIR, G.PackageName, prefix)
		}
		prefix = filepath.Join(G.CwdToRoot, prefix)

		// we need to check if the base directory exists, becuase we have defaults in the schema
		_, err := os.Stat(prefix)
		if err != nil {
			fmt.Printf("warning: %s not found for %s, if you do not intend to use this directory, set 'Partials: []'\n", tg.TrimPrefix, G.PackageName)
			continue
		}

		for _, glob := range tg.Globs {
			// setup vars
			glob = filepath.Clean(glob)

			if !strings.HasPrefix(glob, "/") {

				if G.PackageName != "" {
					glob = filepath.Join(CUE_VENDOR_DIR, G.PackageName, glob)
				}

				// this is how we deal with running generators in the same module
				// they are defined in, while keeping the path spec for them simple
				glob = filepath.Join(G.CwdToRoot, glob)
			}

			pMap, err := templates.CreateTemplateMapFromFolder(glob, prefix, tg.Delims)
			if G.Verbosity > 1 {
				fmt.Printf("%s:%s has %d partial matches\n", G.NamePath(), glob, len(pMap))
			}

			if err != nil {
				errs = append(errs, err)
				continue
			}

			for k, T := range pMap {
				_, ok := G.PartialsMap[k]
				if !ok {
					if G.Verbosity > 1 {
						fmt.Printf(" +p %s:%s\n", G.NamePath(), k)
					}
					// TODO, do we also want to namespace with the template module name?
					G.PartialsMap[k] = T
				} else {
					errs = append(errs, fmt.Errorf("duplicate partial %s:%s", G.NamePath(), k))
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

		_, ok := G.TemplateMap[path]
		if !ok {
			if G.Verbosity > 1 {
				fmt.Printf(" +t %s:%s\n", G.NamePath(), path)
			}

			// TODO, do we also want to namespace with the template module name?
			G.TemplateMap[path] = T
		} else {
			errs = append(errs, fmt.Errorf("duplicate template %s:%s", G.NamePath(), path))
		}
	}

	for _, tg := range G.Templates {
		prefix := filepath.Clean(tg.TrimPrefix)
		if G.PackageName != "" {
			prefix = filepath.Join(CUE_VENDOR_DIR, G.PackageName, prefix)
		}
		prefix = filepath.Join(G.CwdToRoot, prefix)

		// we need to check if the base directory exists, becuase we have defaults in the schema
		_, err := os.Stat(prefix)
		if err != nil {
			fmt.Printf("warning: %s not found in %s, if you do not intend to use this directory, set 'Templates: []'\n", tg.TrimPrefix, G.PackageName)
			continue
		}

		for _, glob := range tg.Globs {
			// setup vars
			glob = filepath.Clean(glob)

			// if not an absolute path
			if !strings.HasPrefix(glob, "/") {
				if G.PackageName != "" {
					glob = filepath.Join(CUE_VENDOR_DIR, G.PackageName, glob)
				}

				// this is how we deal with running generators in the same module
				// they are defined in, while keeping the path spec for them simple
				// note, these will be no-ops when there is no cue.mod
				glob = filepath.Join(G.CwdToRoot, glob)
				prefix = filepath.Join(G.CwdToRoot, prefix)
			}

			pMap, err := templates.CreateTemplateMapFromFolder(glob, prefix, tg.Delims)
			if G.Verbosity > 1 {
				fmt.Printf("%s:%s has %d template matches\n", G.NamePath(), glob, len(pMap))
			}

			if err != nil {
				errs = append(errs, err)
				continue
			}

			for k, T := range pMap {
				_, ok := G.TemplateMap[k]
				if !ok {
					if G.Verbosity > 1 {
						fmt.Printf(" +t %s:%s\n", G.NamePath(), k)
					}

					// TODO, do we also want to namespace with the template module name?
					G.TemplateMap[k] = T
				} else {
					errs = append(errs, fmt.Errorf("duplicate template %s:%s", G.NamePath(), k))
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
		// support text/template in output file path
		if strings.Contains(F.Filepath, "{{") {
			ft, err := templates.CreateFromString("filepath", F.Filepath, nil)
			if err != nil {
				errs = append(errs, err)
			}
			bs, err := ft.Render(F.In)
			if err != nil {
				errs = append(errs, err)
			}
			F.Filepath = string(bs)
		}

		F.Filepath = filepath.Clean(F.Filepath)

		// check for collisions
		if old,ok := G.Files[F.Filepath]; ok {
			fmt.Printf("WARN: duplicate generated file %q in %q & %q\n", F.Filepath, G.NamePath(), old.parent.NamePath())
			// errs = append(errs, fmt.Errorf("duplicate generated file %q in %q", F.Filepath, G.NamePath()))
			continue
		}

		if G.Verbosity > 1 {
			fmt.Printf(" +f %s:%s\n", G.NamePath(), F.Filepath)
		}

		F.parent = G

		G.Files[F.Filepath] = F
		G.OrderedFiles = append(G.OrderedFiles, F)
	}

	for _, F := range G.OrderedFiles {
		err := G.ResolveFile(F)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		ts := make([]string, 0, len(G.TemplateMap))
		for k,_ := range G.TemplateMap {
			ts = append(ts, k)
		}
		errs = append(errs, fmt.Errorf("%s templates: %v", G.NamePath(), ts))
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
		T.AddGolangHelpers()
		t.Parse(P.Source)
	}

	return nil
}
