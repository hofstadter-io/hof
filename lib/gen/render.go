package gen

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"cuelang.org/go/cue"
	"github.com/fatih/color"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/cuetils"
	"github.com/hofstadter-io/hof/lib/templates"
)

type RenderConfig struct {
	// What's loaded by CUE
	RootValue        cue.Value

	// Template configuration
	TemplateConfigs  []RenderTemplateConfig
	Partials         []string

	// internal working fields
	G *Generator
}

// parsed version of the --template flag
// semicolon separated: <filepath>:<?cuepath>@<schema>;[]<?outpath>
// each extra section is 
type RenderTemplateConfig struct {
	// Template filepath
	Filepath string
	
	// CUE path to input value within global value
	Cuepath  string

	// CUE path to schema value within global value
	Schema string

	// Filepath to write results, possibly templated
	Outpath  string

	// Is this a repeated template
	Repeated bool
}

func Render(args []string, rootflags flags.RootPflagpole, cmdflags flags.GenFlagpole) error {

	var RC RenderConfig

	// set partials
	RC.Partials = cmdflags.Partial

	// parse template flags
	for _, tf := range cmdflags.Template {
		cfg, err := parseTemplateFlag(tf)
		if err != nil {
			return err
		}
		RC.TemplateConfigs = append(RC.TemplateConfigs, cfg)
	}

	// load CUE
	crt, err := cuetils.CueRuntimeFromEntrypointsAndFlags(args)
	if err != nil {
		return err
	}

	RC.RootValue = crt.CueValue

	RC.G = NewGenerator("HofRenderCmd", RC.RootValue)
	RC.G.UseDiff3 = cmdflags.Diff3

	// everything loaded, now process
	err = RC.setupGenerator()
	if err != nil {
		return err
	}

	// run this generator
	errs := RC.G.GenerateFiles()
	if len(errs) > 0 {
		fmt.Println(errs)
		return errs[0]
	}

	return nil
}

// deconstructs the flag into struct
// semicolon separated: <filepath>:<?cuepath>@<schema>;<?outpath>
func parseTemplateFlag(tf string) (cfg RenderTemplateConfig, err error) {
	// look for ;
	parts := strings.Split(tf, ";")
	if len(parts) > 1 {
		tf = parts[0]
		cfg.Outpath = parts[1]
	}
	// repeated template?
	if strings.HasPrefix(cfg.Outpath, "[]") {
		cfg.Outpath = strings.TrimPrefix(cfg.Outpath, "[]")
		cfg.Repeated = true
	}

	// look for @
	parts = strings.Split(tf, "@")
	if len(parts) > 1 {
		tf = parts[0]
		cfg.Schema = parts[1]
	}

	// look for :
	parts = strings.Split(tf, ":")
	if len(parts) > 1 {
		tf = parts[0]
		cfg.Cuepath = parts[1]
		// '.' is an alias for "" or the root value
		if cfg.Cuepath == "." {
			cfg.Cuepath = ""
		}
	}

	// should only have template path left
	cfg.Filepath = tf

	return cfg, nil
}

func (RC *RenderConfig) setupGenerator() (err error) {
	G := RC.G

	// setup partials in Generator
	G.Partials = []*TemplateGlobs{ &TemplateGlobs{ Globs: RC.Partials }}

	// setup templates in Generator
	tgs := []string{}
	for _, cfg := range RC.TemplateConfigs {
		tgs = append(tgs, cfg.Filepath)
	}
	G.Templates = []*TemplateGlobs{ &TemplateGlobs{ Globs: tgs }}

	// possibly load shadow dir
	if G.UseDiff3 {
		G.Shadow, err = LoadShadow(G.Name, false)
		if err != nil {
			return err
		}
	}

	// setup Out fields in Generator from our template configs
	err = RC.setupTemplateConfigs()
	if err != nil {
		return err
	}

	// everything should be set to init now
	errs := G.Initialize()
	if len(errs) > 0 {
		fmt.Println(errs)
		return fmt.Errorf("while initializing generator")
	}

	// now run the generator
	errs = G.GenerateFiles()
	if len(errs) > 0 {
		fmt.Println(errs)
		return fmt.Errorf("while rendering templates")
	}

	errs = RC.writeOutput()
	if len(errs) > 0 {
		fmt.Println(errs)
		return fmt.Errorf("while writing output")
	}

	if G.UseDiff3 {
		for _, F := range G.Files {
			if F.IsConflicted > 0 {
				msg := fmt.Sprint("MERGE CONFLICT in:", F.Filepath)
				color.Red(msg)
			}
		}
	}

	return nil
}

// setup Out []*Files in Generator
func (RC *RenderConfig) setupTemplateConfigs() (err error) {
	stdout := 0

	for _, cfg := range RC.TemplateConfigs {

		Val := RC.RootValue

		if cfg.Cuepath != "" {
			Val = Val.LookupPath(cue.ParsePath(cfg.Cuepath))

			// unify with a schema
			if cfg.Schema != "" {
				s := RC.RootValue.LookupPath(cue.ParsePath(cfg.Schema))
				Val = Val.Unify(s)
			}
		}

		addFile := func(val cue.Value) (err error) {
			f := new(File)
			f.TemplatePath = cfg.Filepath

			err = val.Decode(&f.In)
			if err != nil {
				return err
			}

			// Set output filepath, always render as a template
			op := cfg.Outpath
			if op == "" {
				op = fmt.Sprintf("hof-stdout-%d", stdout)
				stdout += 1
			}
			ft, err := templates.CreateFromString("outpath", op, nil)
			if err != nil {
				return err
			}
			bs, err := ft.Render(f.In)
			if err != nil {
				return err
			}
			f.Filepath = string(bs)

			RC.G.Out = append(RC.G.Out, f)
			return nil
		}

		// check if val is a list
		if cfg.Repeated {
			if iter, ierr := Val.List(); ierr == nil {
				for iter.Next() {
					val := iter.Value()
					err := addFile(val)
					if err != nil {
						return err
					}
				}
			} else if iter, ierr := Val.Fields(); ierr == nil {
				for iter.Next() {
					val := iter.Value()
					err := addFile(val)
					if err != nil {
						return err
					}
				}
			} else {
				return fmt.Errorf("repeated template value is not iterable")
			}
		} else {
			err := addFile(Val)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func (RC *RenderConfig) writeOutput() (errs []error) {
	G := RC.G

	writestart := time.Now()

	// Finally write the generator files
	for _, F := range G.OrderedFiles {
		// Write the actual output
		if F.DoWrite && len(F.Errors) == 0 {
			err := F.WriteOutput()
			if err != nil {
				errs = append(errs, err)
				return errs
			}
		}

		if G.UseDiff3 {
			// Write the shadow too, or if it doesn't exist
			if F.DoWrite || (F.IsSame > 0 && F.ShadowFile == nil) {
				err := F.WriteShadow(path.Join(SHADOW_DIR, G.Name))
				if err != nil {
					errs = append(errs, err)
					return errs
				}
			}

			// remove from shadows map so we can cleanup what remains
			delete(G.Shadow, path.Join(G.Name, F.Filepath))
		}
	}

	// Cleanup File & Shadow
	// fmt.Println("Clean Shadow", G.Name)
	if G.UseDiff3 {
		for f, _ := range G.Shadow {
			genFilename := strings.TrimPrefix(f, G.Name+"/")
			shadowFilename := path.Join(SHADOW_DIR, f)
			fmt.Println("  -", G.Name, f, genFilename, shadowFilename)

			err := os.Remove(genFilename)
			if err != nil {
				if strings.Contains(err.Error(), "no such file or directory") {
					continue
				}
				errs = append(errs, err)
				return errs
			}

			err = os.Remove(shadowFilename)
			if err != nil {
				if strings.Contains(err.Error(), "no such file or directory") {
					continue
				}
				errs = append(errs, err)
				return errs
			}

			G.Stats.NumDeleted += 1
		}
	}

	writeend := time.Now()
	G.Stats.WritingTime = writeend.Sub(writestart).Round(time.Millisecond)

	return errs
}
