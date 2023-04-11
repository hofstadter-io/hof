package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"cuelang.org/go/cue"

	"github.com/hofstadter-io/hof/lib/gen"
	"github.com/hofstadter-io/hof/lib/hof"
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

	// Is this a data file? What type?
	DataFormat string
}

func (R *Runtime) CreateAdhocGenerator() error {
	if len(R.GenFlags.Template) == 0 {
		if R.Flags.Verbosity > 1 {
			fmt.Println("Skipping Ad-hoc Generator")
		}
		return nil
	}

	// parse template flags
	tcfgs := []AdhocTemplateConfig{}
	globs := make([]string,0)
	for _, tf := range R.GenFlags.Template {
		cfg, err := ParseTemplateFlag(tf)
		if err != nil {
			return err
		}
		if R.Flags.Verbosity > 2 {
			fmt.Printf("%s -> %#v\n", tf, cfg)
		}
		tcfgs = append(tcfgs, cfg)
		if cfg.Filepath != "" {
			globs = append(globs, cfg.Filepath)
		}
	}

	var h hof.Hof
	h.Label = "AdhocGen"
	h.Path = R.Value.Path().String()
	h.Gen.Root = true
	h.Gen.Name = "AdhocGen"
	node := &hof.Node[gen.Generator]{
		Value: R.Value,
		Hof: h,
	}
	G := gen.NewGenerator(node)
	// reset some vals for ad-hoc
	G.CwdToRoot = ""
	G.Outdir = ""

	G.Templates = []*gen.TemplateGlobs{ &gen.TemplateGlobs{Globs: globs} }
	G.Partials  = []*gen.TemplateGlobs{ &gen.TemplateGlobs{Globs: R.GenFlags.Partial} }

	Val := R.Value


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
			f := new(gen.File)

			// we need this for rendering the output
			// and/or setting the input for the file
			var V any
			err = val.Decode(&V)
			if err != nil {
				return err
			}

			// data or template file
			if cfg.DataFormat != "" {
				// fmt.Println("data file:", cfg.DataFormat)
				f.DatafileFormat = cfg.DataFormat
				f.Value = val
			} else {
				f.TemplatePath = cfg.Filepath
				f.In = V
			}

			// Set output filepath, always render as a template
			op := cfg.Outpath
			if op == "" {
				// when -O is set, -T will replicate the template name to the outpath name
				if R.GenFlags.Outdir != "" {
					op = cfg.Filepath
				} else {
					op = fmt.Sprintf("hof-stdout-%d", stdout)
					stdout += 1
				}
			}

			// 
			ft, err := templates.CreateFromString("outpath", op, nil)
			if err != nil {
				return err
			}
			bs, err := ft.Render(V)
			if err != nil {
				return err
			}
			f.Filepath = string(bs)

			/*
			if cfg.DataFormat != "" {
				fmt.Println(*f)
			}
			*/

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
		return fmt.Errorf("while initializing ad-hoc generator")
	}

	if R.Flags.Verbosity > 2 {
		fmt.Printf("AdhocGen: %#v\n", G)
	}

	R.Generators = append(R.Generators, G)
	return nil
}

// deconstructs the flag into struct
// semicolon separated: <filepath>:<?cuepath>@<schema>=<?outpath>
func ParseTemplateFlag(tf string) (cfg AdhocTemplateConfig, err error) {
	// We work our way from end to start of the string, 

	// look for =
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
	// if no template, then data file format
	// infer from Outpath ext
	if tf == "" {
		cfg.DataFormat = filepath.Ext(cfg.Outpath)[1:]  // trim '.' from ext
	} else {
		cfg.Filepath = tf
	}

	return cfg, nil
}
