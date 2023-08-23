package cuecmd

import (
	"fmt"
	"time"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/format"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
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
		return err
	}

	val := R.Value
	if val.Err() != nil {
		return val.Err()
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

	pkg := R.BuildInstances[0].PkgName
	err = writeOutput(val, pkg, opts, fopts, cflags.Out, cflags.Outfile, cflags.Expression, cflags.Escape, false)
	if err != nil {
		return err
	}

	return nil
}
