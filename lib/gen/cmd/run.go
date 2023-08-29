package cmd

import (
	"fmt"
	"os"
	"time"

	"cuelang.org/go/cue"

	flowcontext "github.com/hofstadter-io/hof/flow/context"
	"github.com/hofstadter-io/hof/flow/middleware"
	"github.com/hofstadter-io/hof/flow/tasks"
	"github.com/hofstadter-io/hof/flow/flow"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	hfmt "github.com/hofstadter-io/hof/lib/fmt"
)

func Run(args []string, rflags flags.RootPflagpole, gflags flags.GenFlagpole) error {
	// return GenLast(args, rootflags, cmdflags)
	verystart := time.Now()

	err := run(args, rflags, gflags)
	if err != nil {
		return err
	}

	// final timing
	veryend := time.Now()
	elapsed := veryend.Sub(verystart).Round(time.Millisecond)

	if rflags.Stats {
		fmt.Printf("\nGrand Total Elapsed Time: %s\n\n", elapsed)
	}

	return nil
}

func run(args []string, rflags flags.RootPflagpole, gflags flags.GenFlagpole) error {

	err := hfmt.UpdateFormatterStatus()
	if err != nil {
		return fmt.Errorf("update formatter status: %w", err)
	}

	R, err := prepRuntime(args, rflags, gflags)
	if err != nil {
		return err
	}

	// we need generators loaded at this point
	if R.GenFlags.AsModule != "" {
		return R.adhocAsModule()
	}

	// run pre-exec here
	for _, G := range R.Generators {

		// maybe run pre-flow per generator
		preFlow := G.CueValue.LookupPath(cue.ParsePath("PreFlow"))
		if preFlow.Exists() {
			if R.Flags.Verbosity > 0 {
				fmt.Println("running pre-flow:", preFlow)
			}
			if !R.GenFlags.Exec {
				fmt.Println("skipping pre-flow, use --exec to run")
			} else {
				ctx := flowcontext.New()
				ctx.RootValue = preFlow
				ctx.Stdin = os.Stdin
				ctx.Stdout = os.Stdout
				ctx.Stderr = os.Stderr
				ctx.Verbosity = R.Flags.Verbosity

				// how to inject tags into original value
				// fill / return value
				middleware.UseDefaults(ctx, R.Flags, flags.FlowPflags)
				tasks.RegisterDefaults(ctx)

				p, err := flow.OldFlow(ctx, preFlow)
				if err != nil {
					return err
				}

				err = p.Start()
				if err != nil {
					return err
				}

				G.CueValue = G.CueValue.FillPath(cue.ParsePath("PreFlow"), preFlow)
				if G.CueValue.Err() != nil {
					return err
				}
			}

		} else if !preFlow.Exists() {
			if G.Verbosity > 0 {
				fmt.Println("pre-exec not found")
			}
		} else if preFlow.Err() != nil {
			return preFlow.Err()
		}
	}

	return R.runGen()
}

