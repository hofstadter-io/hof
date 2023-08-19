package cuecmd

import (
	"cuelang.org/go/cue"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/cuetils"
	"github.com/hofstadter-io/hof/lib/runtime"
)

func Eval(args []string, rflags flags.RootPflagpole, cflags flags.EvalFlagpole) error {

	R, err := runtime.New(args, rflags)
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

	// build options
	opts := []cue.Option{
		cue.Concrete(cflags.Concrete),
		cue.Docs(cflags.Comments),
		cue.Attributes(cflags.Attributes),
		cue.Optional(cflags.Optional),
		cue.Hidden(cflags.Hidden),
		cue.InlineImports(cflags.InlineImports),
		cue.ErrorsAsValues(rflags.IngoreErrors),
	}

	err = writeOutput(val, opts, cflags.Out, cflags.Outfile, cflags.Expression)
	if err != nil {
		return err
	}

	return nil
}

