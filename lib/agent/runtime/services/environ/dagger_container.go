package environ

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

type daggerContainer struct {
	ctx context.Context
	cfg *Config

	dag *dagger.Client
	c   *dagger.Container
}

type Config struct {
	From string
	Env  map[string]string
}

func NewContainer(ctx context.Context, cfg *Config, dag *dagger.Client) (dc *daggerContainer, err error) {
	// internal wrapper / representation
	dc = &daggerContainer{
		ctx: ctx,
		cfg: cfg,
		dag: dag,
	}

	// start dagger environment
	// TODO, check if it starts with a sha256 and look up
	if cfg.From != "" {
		dc.c = dag.Container().From(cfg.From)
	} else {
		// blank, we can swap out later and copy directories as needed
		// copying files like multi-stage docker...? are we starting to proxy too much API? (can our code gen and ai help with this)
		dc.c = dag.Container()
	}

	// TODO attach env vars
	// TODO attach secrets

	// sync so we have have checked things are good at this point?
	dc.c, err = dc.c.Sync(dc.ctx)
	if err != nil {
		return dc, fmt.Errorf("while initializing container: %w", err)
	}

	return dc, nil
}

func (dc *daggerContainer) ID() (string, error) {
	id, err := dc.c.ID(dc.ctx)
	return string(id), err
}

func (dc *daggerContainer) Load(id string) (*dagger.Container, error) {
	return nil, nil
}
