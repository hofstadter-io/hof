package cuecmd

import (
	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/datamodel"
	"github.com/hofstadter-io/hof/lib/runtime"
)

func makeRuntime(cmd string, args []string, rflags flags.RootPflagpole) (*runtime.Runtime, error) {
	// fmt.Printf("lib/datamodel.Run.%s %v %v %v\n", cmd, args, rflags, dflags)

	R, err := runtime.New(args, rflags)
	if err != nil {
		return R, err
	}

	err = R.Load()
	if err != nil {
		return R, err
	}

	err = R.EnrichDatamodels(nil, EnrichDatamodel)
	if err != nil {
		return R, err
	}

	return R, nil
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
