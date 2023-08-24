package cuecmd

import (
	"fmt"
	"time"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/format"
	"github.com/kr/pretty"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/cuetils"
	"github.com/hofstadter-io/hof/lib/runtime"
)

func Eval(args []string, rflags flags.RootPflagpole, cflags flags.EvalFlagpole) error {
	start := time.Now()
	R, err := runtime.New(args, rflags)

	defer func() {
		if R.Flags.Stats {
			fmt.Println(R.Stats)
			end := time.Now()
			fmt.Printf("\nTotal Elapsed Time: %s\n\n", end.Sub(start))
		}
	}()

	wantErrors := rflags.IngoreErrors || rflags.AllErrors

	if err != nil {
		return cuetils.ExpandCueError(err)
	}

	err = R.Load()
	if err != nil && !wantErrors {
		// fmt.Println("load.Err", err)
		return cuetils.ExpandCueError(err)
	}

	val := R.Value
	// fmt.Println("val:", val.Err())
	if val.Err() != nil && !wantErrors {
		// fmt.Println("val.Err", val.Err())
		return cuetils.ExpandCueError(val.Err())
	}

	if R.Flags.Verbosity > 1 {
		fmt.Printf("%# v\n", pretty.Formatter(R.Flags))
		fmt.Printf("%# v\n", pretty.Formatter(cflags))
	}

	// build options
	opts := []cue.Option{
		cue.Docs(cflags.Comments),
		cue.Attributes(cflags.Attributes),
		cue.Definitions(cflags.Definitions),
		cue.Optional(cflags.Optional || cflags.All),
		cue.InlineImports(cflags.InlineImports),
		cue.ErrorsAsValues(wantErrors),
		cue.ResolveReferences(cflags.Resolve),
	}

	// these two have to be done specially
	// because there are three options [true, false, missing]
	if cflags.Concrete {
		opts = append(opts, cue.Concrete(true))
	}
	if cflags.Hidden || cflags.All {
		opts = append(opts, cue.Hidden(true))
	}

	if cflags.Final {
		// prepend final, so others still apply
		opts = append([]cue.Option{cue.Final()}, opts...)
	}

	fopts := []format.Option{}
	if cflags.Simplify {
		fopts = append(fopts, format.Simplify())
	}

	bi := R.BuildInstances[0]
	if R.Flags.Verbosity > 1 {
		fmt.Println("ID:", bi.ID(), bi.PkgName, bi.Module)
	}
	pkg := bi.PkgName
	if bi.Module == "" {
		pkg = bi.ID()
	}
	err = writeOutput(val, pkg, opts, fopts, cflags.Out, cflags.Outfile, cflags.Expression, rflags.Schema, cflags.Escape, cflags.Defaults, wantErrors)
	if err != nil {
		return err
	}

	return nil
}

