package cuecmd

import (
	"cuelang.org/go/cue"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/runtime"
)

func Def(args []string, rflags flags.RootPflagpole, cflags flags.DefFlagpole) error {

	R, err := runtime.New(args, rflags)
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
		cue.Docs(true),
	}
	if cflags.Attributes {
		opts = append(opts, cue.Attributes(true))
	}
	if cflags.InlineImports {
		opts = append(opts, cue.InlineImports(true))
	}

	err = writeOutput(val, opts, cflags.Outfile, cflags.Expression)
	if err != nil {
		return err
	}

	return nil
}
