package cmd

import (
	"fmt"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/datamodel"
	"github.com/hofstadter-io/hof/lib/runtime"
)

func Run(cmd string, args []string, rflags flags.RootPflagpole, dflags flags.DatamodelPflagpole) error {
	// fmt.Printf("lib/datamodel.Run.%s %v %v %v\n", cmd, args, rflags, dflags)

	R, err := runtime.New(args, rflags)
	if err != nil {
		return err
	}

	err = R.Load()
	if err != nil {
		return err
	}

	err = R.EnrichDatamodels(dflags.Datamodels, EnrichDatamodel)
	if err != nil {
		return err
	}

	// Now, with our datamodles at hand, run the command
	switch cmd {
	case "list":
		err = list(R, dflags)

	case "tree":
		err = tree(R, dflags)

	case "checkpoint":
		err = checkpoint(R, dflags, flags.Datamodel__CheckpointFlags)

	case "diff":
		err = diff(R, dflags)

	case "log":
		err = log(R, dflags)

	default:
		err = fmt.Errorf("%s command not implemented yet", cmd)
	}

	return err
}

func EnrichDatamodel(R *runtime.Runtime, dm *datamodel.Datamodel) error {
	err := dm.LoadHistory()
	if err != nil {
		return err
	}
	err = dm.CalcDiffs()
	if err != nil {
		return err
	}

	return nil
}
