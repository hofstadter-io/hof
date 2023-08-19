package cuecmd

import (
	"fmt"
	"os"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/format"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/datamodel"
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

func writeOutput(val cue.Value, opts []cue.Option, outfile string, exs []string) (err error) {
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

	for i, ex := range exs {
		if i > 0 {
			fmt.Fprint(out, "// ---\n")
		}

		var v cue.Value

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

		if v.Err() != nil {
			return v.Err()
		}
		
		// get formatted value
		syn := v.Syntax(opts...)
		b, err := format.Node(syn)
		if err != nil {
			return err
		}

		fmt.Fprint(out, string(b))
	}

	return nil
}
