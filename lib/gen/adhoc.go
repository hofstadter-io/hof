package gen

import (
	"fmt"
	"strings"

	"cuelang.org/go/cue"

	"github.com/hofstadter-io/hof/lib/templates"
)

// parsed version of the --template flag
// semicolon separated: <filepath>:<?cuepath>@<schema>;[]<?outpath>
// each extra section is 
type AdhocTemplateConfig struct {
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

func (R *Runtime) CreateAdhocGenerator() error {
	if len(R.Flagpole.Template) == 0 {
		if R.Verbosity > 1 {
			fmt.Println("Skipping Adhoc Generator")
		}
		return nil
	}

	if R.Verbosity > 1 {
		fmt.Println("Creating Adhoc Generator")
	}
	// parse template flags
	tcfgs := []AdhocTemplateConfig{}
	globs := make([]string,0)
	for _, tf := range R.Flagpole.Template {
		cfg, err := parseTemplateFlag(tf)
		if err != nil {
			return err
		}
		if R.Verbosity > 2 {
			fmt.Printf("%s -> %#v\n", tf, cfg)
		}
		tcfgs = append(tcfgs, cfg)
		globs = append(globs, cfg.Filepath)
	}

	G := NewGenerator("AdhocGen", R.CueRuntime.CueValue)
	G.UseDiff3 = R.Flagpole.Diff3
	G.Outdir = ""

	G.Templates = []*TemplateGlobs{ &TemplateGlobs{Globs: globs} }
	G.Partials  = []*TemplateGlobs{ &TemplateGlobs{Globs: R.Flagpole.Partial} }

	Val := R.CueRuntime.CueValue


	stdout := 0
	for _, cfg := range tcfgs {
		val := Val
		if cfg.Cuepath != "" {
			val = val.LookupPath(cue.ParsePath(cfg.Cuepath))
		}
		if cfg.Schema != "" {
			schema := Val.LookupPath(cue.ParsePath(cfg.Schema))
			val = val.Unify(schema)
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
				// when -O is set, -T will replicate the template name to the outpath name
				if R.Flagpole.Outdir != "" {
					op = cfg.Filepath
				} else {
					op = fmt.Sprintf("hof-stdout-%d", stdout)
					stdout += 1
				}
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

			G.Out = append(G.Out, f)
			return nil
		}

		if cfg.Repeated {
			// check if val is a list
			if iter, ierr := val.List(); ierr == nil {
				for iter.Next() {
					v := iter.Value()
					err := addFile(v)
					if err != nil {
						return err
					}
				}
			// check if val is a struct
			} else if iter, ierr := val.Fields(); ierr == nil {
				for iter.Next() {
					v := iter.Value()
					err := addFile(v)
					if err != nil {
						return err
					}
				}
			} else {
				return fmt.Errorf("repeated template value is not iterable")
			}
		} else {
			err := addFile(val)
			if err != nil {
				return err
			}
		}

	}

	errs := G.Initialize()
	if len(errs) > 0 {
		fmt.Println(errs)
		return fmt.Errorf("while initializing adhoc generator")
	}

	if R.Verbosity > 2 {
		fmt.Printf("G: %#v\n", G)
	}

	R.Generators["AdhocGen"] = G
	return nil
}

// deconstructs the flag into struct
// semicolon separated: <filepath>:<?cuepath>@<schema>;<?outpath>
func parseTemplateFlag(tf string) (cfg AdhocTemplateConfig, err error) {
	// look for ;
	parts := strings.Split(tf, "=")
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
