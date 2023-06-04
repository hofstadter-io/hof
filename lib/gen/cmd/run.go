package cmd

import (
	"fmt"
	"time"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
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

	if gflags.Stats {
		fmt.Printf("\nGrand Total Elapsed Time: %s\n\n", elapsed)
	}

	return nil
}

func run(args []string, rflags flags.RootPflagpole, gflags flags.GenFlagpole) error {

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

