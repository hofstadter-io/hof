package cmd

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/cuetils"
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
		return nil, cuetils.ExpandCueError(err)
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

const DAGGER_HOST = "container://veg-dagger-engine"

func daggerClient(ctx context.Context) (*dagger.Client, error) {
	os.Setenv("_EXPERIMENTAL_DAGGER_RUNNER_HOST", DAGGER_HOST)
	client, err := dagger.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("while connecting to dagger: %w", err)
	}
	return client, nil
}

func splitArgs(orig []string) (args, cueargs []string) {
	args = orig
	for i, a := range orig {
		if a == "%" {
			args = orig[:i]
			if i+1 < len(orig) {
				cueargs = orig[i+1:]
			}
			break
		}
	}
	return args, cueargs
}
