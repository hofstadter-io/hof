package cuecmd

import (
	"fmt"
	"os"
	"time"

	"cuelang.org/go/cue"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/runtime"
)

func Vet(args []string, rflags flags.RootPflagpole, cflags flags.VetFlagpole) error {

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
		cue.Docs(cflags.Comments),
		cue.Attributes(cflags.Attributes),
		cue.Definitions(cflags.Definitions),
		cue.Optional(cflags.Optional),
		cue.ErrorsAsValues(rflags.IngoreErrors || rflags.AllErrors),
	}

	// these two have to be done specially
	// because there are three options [true, false, missing]
	if cflags.Concrete {
		opts = append(opts, cue.Concrete(true))
	}
	if cflags.Hidden {
		opts = append(opts, cue.Hidden(true))
	}

	out := os.Stdout
	exs := cflags.Expression
	if len(exs) == 0 {
		exs = []string{""}
	}

	hadError := false
	handleErr := func(ex string, err error) {
		if err == nil {
			return
		}
		hadError = true
		if len(exs) > 1 {
			fmt.Fprintln(out, "//", ex)
		}
		fmt.Fprint(out, err)
	}

	// TODO, need to find unplaced data files and validate them

	for _, ex := range exs {

		pkg := R.BuildInstances[0].ID()
		v := getValByEx(ex, pkg, val)
		if v.Err() != nil {
			handleErr(ex, v.Err())
			continue
		}
	
		err := v.Validate(append(opts, )...)
		handleErr(ex, err)
	}

	if hadError {
		// messages already printed, we want an empty message
		return fmt.Errorf("")
	}

	return nil
}
