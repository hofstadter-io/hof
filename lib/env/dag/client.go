package dag

import (
	"context"
	"sync"

	"dagger.io/dagger"
	"github.com/hofstadter-io/hof/lib/env"
)

type Dag struct {
	mx  sync.RWMutex
	ctx context.Context
	dag *dagger.Client

	// should probably just pass in the flags...
	noCache bool

	// catalog allows us to consolidate references across CUE that might get duplicated
	// as well as each entry holding the value, config, and go types for the entire life-cycle
	cat catalog
	hdl stepHandlerMap
}

func NewClient(ctx context.Context, client *dagger.Client) (d *Dag, err error) {
	d = &Dag{
		ctx: ctx,
		dag: client,
		cat: newCatalog(),
	}

	d.hdl = d.makeStepHandlers()

	return d, nil
}

// probably want to expand this to...
// 1. track more than host stuff (containers, dirs, services, etc...)
// 2. track the config, object, id, usage?
// 3. look up by name or id? (need to fill back the id?) or can we create an interface Key() <- better probably
// 4. enough info to track, walk, schedule, and visualize the dag
// 5. an interface or real type here, instead of any
type catalog map[Keyer]any

func newCatalog() catalog {
	return make(map[Keyer]any)
}

type Keyer interface {
	Key() string
}

// preferences #hof: metadata: [memo|id|name]
func vegMemoKey(e *env.Env) string {
	meta := e.Hof.Metadata
	if meta.Memo != "" {
		return meta.Memo
	}
	if meta.ID != "" {
		return meta.ID
	}
	if meta.Name != "" {
		return meta.Name
	}
	return ""
}
