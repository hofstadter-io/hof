package cuecmd

import (
	"fmt"
	"time"

	"cuelang.org/go/cue"

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

	if err != nil {
		return cuetils.ExpandCueError(err)
	}

	err = R.Load()
	if err != nil {
		return cuetils.ExpandCueError(err)
	}

	val := R.Value
	if val.Err() != nil {
		return cuetils.ExpandCueError(val.Err())
	}

	fmt.Println("Defs?", cflags.Definitions)

	// build options
	opts := []cue.Option{
		cue.Concrete(cflags.Concrete),
		cue.Docs(cflags.Comments),
		cue.Definitions(cflags.Definitions),
		cue.Attributes(cflags.Attributes),
		cue.Optional(cflags.Optional),
		cue.Hidden(cflags.Hidden),
		cue.InlineImports(cflags.InlineImports),
		cue.ErrorsAsValues(rflags.IngoreErrors),
	}

	if cflags.Final {
		opts = append(opts, cue.Final())
	}

	err = writeOutput(val, opts, cflags.Out, cflags.Outfile, cflags.Expression)
	if err != nil {
		return err
	}

	return nil
}

