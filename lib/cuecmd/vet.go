package cuecmd

import (
	"fmt"
	"os"
	"time"

	"cuelang.org/go/cue"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/cuetils"
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

	wantErrors := rflags.IngoreErrors || rflags.AllErrors

	// this is a bit hacky (?), but we use this so we can vet the orphaned files rather than add them to the value
	R.DontPlaceOrphanedFiles = true
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
		cue.ErrorsAsValues(wantErrors),
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
		err = cuetils.ExpandCueError(err)
		hadError = true
		if len(exs) > 1 {
			fmt.Fprintln(out, "//", ex)
		}
		fmt.Fprint(out, err)
	}

	// TODO, how do we think about the cross-product of { files } x { -e } x { -l }
	// maybe -l doesn't make sense here? (or only files that can be placed)

	// setup our bi and other stuff
	bi := R.BuildInstances[0]
	if R.Flags.Verbosity > 1 {
		fmt.Println("ID:", bi.ID(), bi.PkgName, bi.Module)
	}
	pkg := bi.PkgName
	if bi.Module == "" {
		pkg = bi.ID()
	}

	// vet the orphaned files
	if len(bi.OrphanedFiles) > 0 {
		for i, f := range bi.OrphanedFiles {
			// fmt.Println("vet:", f.Filename)
			F, err := R.LoadOrphanedFile(f, pkg, bi.Root, bi.Dir, i, len(bi.OrphanedFiles))
			if err != nil {
				handleErr("during load", err)
				continue
			}
			fv := R.CueContext.BuildFile(F, cue.Filename(f.Filename))

			// vet the value with each expression
			for _, ex := range exs {

				v := getValByEx(ex, pkg, val)
				if v.Err() != nil {
					handleErr(ex, v.Err())
					continue
				}
				
				v = v.Unify(fv)
			
				// we want to ensure concrete when validating data (orphaned files)
				opts = append(opts, cue.Concrete(true))
				err := v.Validate(append(opts, )...)
				handleErr(ex, err)
			}

		}
	} else {
		// vet the root value at each expression
		// often this will default to [""] which is just the whole value
		for _, ex := range exs {

			v := getValByEx(ex, pkg, val)
			if v.Err() != nil {
				handleErr(ex, v.Err())
				continue
			}
		
			err := v.Validate(append(opts, )...)
			handleErr(ex, err)
		}
	}



	if hadError {
		// messages already printed, we want an empty message
		return fmt.Errorf("")
	}

	return nil
}
