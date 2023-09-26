package cuecmd

import (
	"fmt"
	"time"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/format"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/cuetils"
	"github.com/hofstadter-io/hof/lib/runtime"
)

func Export(args []string, rflags flags.RootPflagpole, cflags flags.ExportFlagpole) error {

	start := time.Now()
	R, err := runtime.New(args, rflags)

	defer func() {
		if R.Flags.Stats {
			fmt.Println(R.Stats)
			end := time.Now()
			fmt.Printf("\nTotal Elapsed Time: %s\n\n", end.Sub(start))
		}
	}()

	if err != nil {
		return err
	}

	err = R.Load()
	if err != nil {
		return cuetils.ExpandCueError(err)
	}

	val := R.Value
	if val.Err() != nil {
		return cuetils.ExpandCueError(val.Validate())
	}

	// build options
	opts := []cue.Option{
		cue.Concrete(true),
		cue.Final(),
		cue.Docs(cflags.Comments),
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
	err = writeOutput(val, pkg, opts, fopts, cflags.Out, cflags.Outfile, cflags.Expression, rflags.Schema, cflags.Escape, false, false)
	if err != nil {
		return err
	}

	return nil
}
