package cmd

import (
	"fmt"
)

// Clears and reloads a runtime, rereading inputs and reprocessing everything
// fast determines if the CUE code is reloaded and evaluated or not (fast will skip CUE).
func (R *Runtime) Reload(fast bool) (err error) {
	R.Lock()
	defer R.Unlock()

	if R.Flags.Verbosity > 1 {
		fmt.Printf("Runtime.Reload(%v)\n", fast)
	}

	R.Clear()

	if !fast {
		err = R.Load()
		if err != nil {
			return err
		}
	}

	err = R.localLoad()
	if err != nil {
		return err
	}

	return nil
}

func (R *Runtime) localLoad() error {
	err := R.EnrichDatamodels(nil, EnrichDatamodelBuilder(R))
	if err != nil {
		return err
	}

	// fmt.Println(R.Value)

	// inject history & diffs
	//err = R.()
	//if err != nil {
	//  return err
	//}

	// the generators to load up
	gens := R.GenFlags.Generator
	// we want to skip any generators
	// if we are in adhoc mode and haven't set -G
	if len(R.GenFlags.Template) > 0 && len(gens) == 0 {
		// this value should not match user data
		// so we effectively omit all generators, besides adhoc
		gens = []string{"HOF_ADHOC_OMIT_GENERATORS"}
	}

	err = R.EnrichGenerators(gens, EnrichGeneratorBuilder(R))
	if err != nil {
		return err
	}

	err = R.Initialize()
	if err != nil {
		return err
	}

	return nil
}

func (R *Runtime) Initialize() error {
	if R.Flags.Verbosity > 1 {
		fmt.Printf("Runtime.Initialize()\n")
	}

	err := R.CreateAdhocGenerator()
	if err != nil {
		return err
	}

	/*
	for _, G := range R.Generators {
		errs := G.Initialize()
		if len(errs) != 0 {
			var emsg string
			for _, err := range errs {
				emsg += fmt.Sprintf("%s\n", err.Error())
			}
			return fmt.Errorf("while initializing %s:\n%s", G.Hof.Path, emsg)
		}
	}
	*/

	return nil
}

