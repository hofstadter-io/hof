package dag

import (
	"context"
	"sync"

	"dagger.io/dagger"
)

type localEnv struct {
	mx  sync.RWMutex
	ctx context.Context
	dag *dagger.Client
}

var LE *localEnv

// TODO, we should have one or more of these on the Runtime, useful in other commands / packages
// ... but, we only want to connect when we need it to avoid painful delays from establish said connection

// const DAGGER_HOST = "container://veg-dagger-engine"

func Initialize(ctx context.Context) (err error) {
	// os.Setenv("_EXPERIMENTAL_DAGGER_RUNNER_HOST", DAGGER_HOST)
	LE = &localEnv{
		ctx: ctx,
	}
	LE.dag, err = dagger.Connect(ctx)
	return err
}

func Client() *localEnv {
	if LE == nil {
		panic("LocalEnv has not been initialized, you must do so manually as a developer of veg")
	}
	return LE
}
