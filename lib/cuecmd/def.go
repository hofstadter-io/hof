package cuecmd

import (
	"fmt"
	"time"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/format"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/runtime"
)

func Def(args []string, rflags flags.RootPflagpole, cflags flags.DefFlagpole) error {

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
		fmt.Println("Def:", err)
		return err
	}

	// build options
	opts := []cue.Option{
		cue.Docs(cflags.Comments),
		cue.Attributes(cflags.Attributes),
		cue.InlineImports(cflags.InlineImports),
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
	err = writeOutput(R.Value, pkg, opts, fopts, cflags.Out, cflags.Outfile, cflags.Expression, rflags.Schema, false, false, true)
	if err != nil {
		return err
	}

	return nil
}
