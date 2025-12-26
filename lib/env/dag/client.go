package dag

import (
	"context"
	"sync"

	"dagger.io/dagger"
)

type Dag struct {
	mx  sync.RWMutex
	ctx context.Context
	dag *dagger.Client
}

func NewClient(ctx context.Context, client *dagger.Client) (d *Dag, err error) {
	return &Dag{
		ctx: ctx,
		dag: client,
	}, nil
}
