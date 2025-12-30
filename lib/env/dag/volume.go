package dag

import (
	"fmt"

	"cuelang.org/go/cue"
	"dagger.io/dagger"
	"github.com/hofstadter-io/hof/lib/env"
)

type hashCacheConfig struct {
	Kind string `json:"$kind"`
	Name string `json:"name"`
}

type hashCacheIndex struct {
	node *env.Env
	val  cue.Value
	cfg  *hashCacheConfig
	vol  *dagger.CacheVolume
}

func (h *hashCacheIndex) Key() string {
	if h.cfg == nil {
		return "#cache.nil"
	}
	return fmt.Sprintf("#cache.%s", h.cfg.Name)
}

func (d *Dag) hashCache(step cue.Value) (*dagger.CacheVolume, error) {
	var cfg hashCacheConfig
	err := step.Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("while decoding hashCache: %w", err)
	}

	// index for query and create if not found
	idx := &hashCacheIndex{
		val: step,
		cfg: &cfg,
	}

	// lookup
	ia, ok := d.cat[idx]
	if ok {
		ix := ia.(*hashCacheIndex)
		return ix.vol, nil
	}

	// load for realz
	idx.vol = d.dag.CacheVolume(cfg.Name)
	// memoize
	d.cat[idx] = idx

	return idx.vol, nil
}
