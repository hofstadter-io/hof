package cmd

import (
	"fmt"
	"time"

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

	return R.runGen()
}

