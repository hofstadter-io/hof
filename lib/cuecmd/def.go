package cuecmd

import (
	"fmt"
	"os"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/format"

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

	// get formatted value
	syn := val.Syntax(opts...)
	b, err := format.Node(syn)
	if err != nil {
		return err
	}

	// output results
	if cflags.Outfile != "" {
		err = os.WriteFile(cflags.Outfile, b, 0o644)
		if err != nil {
			return err
		}
	} else {
		fmt.Println(string(b))
	}

	return nil
}
