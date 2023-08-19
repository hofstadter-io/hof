package cuecmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/format"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/datamodel"
	"github.com/hofstadter-io/hof/lib/gen"
	"github.com/hofstadter-io/hof/lib/runtime"
)

func makeRuntime(args []string, rflags flags.RootPflagpole) (*runtime.Runtime, error) {

	R, err := runtime.New(args, rflags)
	if err != nil {
		return R, err
	}

	err = R.Load()
	if err != nil {
		return R, err
	}

	err = R.EnrichDatamodels(nil, EnrichDatamodel)
	if err != nil {
		return R, err
	}

	return R, nil
}

func EnrichDatamodel(R *runtime.Runtime, dm *datamodel.Datamodel) error {
	err := dm.LoadHistory()
	if err != nil {
		return err
	}
	err = dm.CalcDiffs()
	if err != nil {
		return err
	}

	return nil
}

func writeOutput(val cue.Value, opts []cue.Option, outtype, outfile string, exs []string) (err error) {
	// when not set, this makes it so our loop will iterate once and output everything
	if len(exs) == 0 {
		exs = append(exs, "")
	}

	out := os.Stdout
	if outfile != "" {
		out, err = os.OpenFile(outfile, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
	}

	hadError := false
	for _, ex := range exs {
		// if more than one output, prefix with name in commment
		if len(exs) > 1 {
			fmt.Fprintln(out, "//", ex)
		}

		var v cue.Value

		if ex == "" {
			// special case
			v = val
		} else {
			p := cue.ParsePath(ex)
			if p.Err() == nil {
				v = val.LookupPath(p)
			} else {
				ctx := val.Context()
				v = ctx.CompileString(
					ex,
					cue.Filename(ex),
					cue.InferBuiltins(true),
					cue.Scope(val),
				)
			}
		}

		if v.Err() != nil {
			hadError = true
			fmt.Fprint(out, v.Err())
			continue
		}

		// we have a value now, without basic errors

		// what encoding for output
		if outtype == "" {
			if outfile == "" {
				// our default when both are empty
				outtype = "cue"
			} else {
				ext := filepath.Ext(outfile)
				outtype = strings.TrimPrefix(ext, ".")
			}

		}


		switch outtype {
		case "cue":
			write := func(n ast.Node) error {
				b, err := format.Node(n)
				if err != nil {	
					fmt.Fprint(out, err)
					return err
				}
				fmt.Fprint(out, string(b))
				return nil
			}

			// get formatted value
			syn := v.Syntax(opts...)
			// hack to remove the extra {} around values when in some situations
			if s, ok := syn.(*ast.StructLit); ok {
				f := &ast.File{
					Decls: s.Elts,
				}
				syn = f
			}

			// eval / write the value
			err = write(syn)
			if err != nil {
				hadError = true
				continue
			}

		case "json":
			b, err := gen.FormatJson(v)
			if err != nil {
				hadError = true
				continue
			}
			fmt.Fprint(out, string(b))

		case "yaml":
			b, err := gen.FormatYaml(v)
			if err != nil {
				hadError = true
				continue
			}
			fmt.Fprint(out, string(b))

		case "xml":
			b, err := gen.FormatXml(v)
			if err != nil {
				hadError = true
				continue
			}
			fmt.Fprintln(out, string(b))

		case "toml":
			b, err := gen.FormatToml(v)
			if err != nil {
				hadError = true
				continue
			}
			fmt.Fprint(out, string(b))
			
		default:
			return fmt.Errorf("unknown output type %s", outtype)	
		}
	}

	if hadError {
		// messages already printed, we want an empty message
		return fmt.Errorf("")
	}

	return nil
}
