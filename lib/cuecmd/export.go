package cuecmd

import (
	"fmt"
	"time"

	"cuelang.org/go/cue"

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

	err = writeOutput(val, opts, cflags.Out, cflags.Outfile, cflags.Expression)
	if err != nil {
		return err
	}

	return nil
}
