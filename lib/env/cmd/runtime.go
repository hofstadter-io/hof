package cmd

import (
	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/env"
	"github.com/hofstadter-io/hof/lib/runtime"
)

func prepRuntime(args []string, rflags flags.RootPflagpole) (*runtime.Runtime, error) {

	// create our core runtime
	r, err := runtime.New(args, rflags)
	if err != nil {
		return nil, err
	}

	err = r.Load()
	if err != nil {
		return nil, err
	}

	err = r.EnrichEnv(nil, EnrichEnv)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func EnrichEnv(R *runtime.Runtime, e *env.Env) error {

	// no-op
	return nil
}
