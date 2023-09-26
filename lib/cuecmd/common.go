package cuecmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/format"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/cuetils"
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

func writeOutput(
	val cue.Value,
	pkg string,
	opts []cue.Option,
	fopts []format.Option,
	outtype, outfile string,
	exs, schemas []string,
	escape, defaults, wantErrors bool,
) (err error) {
	// fmt.Println("writeOutput", pkg, exs)
	// when not set, this makes it so our loop will iterate once and output everything
	if len(exs) == 0 {
		exs = append(exs, "")
	}

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

	// ensure concrete if not cue output
	if outtype != "cue" {
		opts = append(opts, cue.Concrete(true))
	}

	// setup output writer
	out := os.Stdout
	if outfile != "" && outfile != "-" {
		out, err = os.OpenFile(outfile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		defer out.Close()
	}

	// error handling, so we can still process everything
	hadError := false
	handleErr := func(err error, ex string) {
		hadError = true
		if len(exs) > 1 {
			fmt.Fprintln(os.Stderr, "//", ex)
		}
		fmt.Fprint(os.Stderr, cuetils.ExpandCueError(err))
	}

	handleStuff := func(err error, s, ex string) {
		if err != nil {
			handleErr(err, ex)
			return
		}
		if len(exs) > 1 {
			fmt.Fprintln(out, "//", ex)
		}
		fmt.Fprint(out, s)
	}

	// range of expressions the user desires
	for _, ex := range exs {
		// if more than one output, prefix with name in commment
		v := getValByEx(ex, pkg, val)
		if !v.Exists() {
			handleErr(v.Err(), ex)
			continue
		}
		//if !wantErrors && v.Err() != nil {
		//  handleErr(v.Err(), ex)
		//  continue
		//}

		if defaults {
			v, _ = v.Default()
		}

		for _, schema := range schemas {
			s := getValByEx(schema, pkg, val)
			// we don't ignore here because we want to actually have these schemas in the value to use
			// and ignore any errors the data may have against them
			if s.Err() != nil {
				return fmt.Errorf("unable to find schema in value: %w\n", s.Err())
			}

			// keep unifying with all schemas
			v = v.Unify(s)
		}

		if !wantErrors {
			err = v.Validate(opts...)
			if err != nil {
				handleErr(err, ex)
				continue
			}
		}

		// we have a good(ish) value now, without basic errors

		switch outtype {
		case "cue":
			// get formatted value
			syn := v.Syntax(opts...)
			// hack to remove the extra {} around values when in some situations
			syn = cuetils.ToFile(syn)

			b, err := format.Node(syn, fopts...)
			handleStuff(err, string(b), ex)

		case "json":
			b, err := gen.FormatJson(v, escape)
			handleStuff(err, string(b), ex)

		case "yaml":
			b, err := gen.FormatYaml(v)
			handleStuff(err, string(b), ex)

		case "xml":
			b, err := gen.FormatXml(v)
			handleStuff(err, string(b), ex)

		case "toml":
			b, err := gen.FormatToml(v)
			handleStuff(err, string(b), ex)
			
		case "text":
			s, err := v.String()
			handleStuff(err, s, ex)
			
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

func getValByEx(ex, pkg string, val cue.Value) cue.Value {
	if ex == "" || ex == "." {
		return val
	} else {
		p := exToPath(ex, pkg)
		if p.Err() == nil {
			return val.LookupPath(p)
		} else {
			ctx := val.Context()
			return ctx.CompileString(
				ex,
				cue.Filename("--expression:"+ex),
				cue.InferBuiltins(true),
				cue.Scope(val),
			)
		}
	}
}

func exToPath(ex, pkg string) (cue.Path) {
	if pkg == "" {
		pkg = "_"
	}
	var sels []cue.Selector
	// assume we can split on dots
	parts := strings.Split(ex, ".")
	for _, part := range parts {
		if strings.HasPrefix(part, "_") {
			sels = append(sels, cue.Hid(part, pkg))
			// fmt.Println("SELS", pkg, sels)
		} else {
			p := cue.ParsePath(part)
			sels = append(sels, p.Selectors()...)
			// fmt.Printf("P: %#+v %v\n", p.Selectors(), p.Err())
		}
	}

	return cue.MakePath(sels...)
}
