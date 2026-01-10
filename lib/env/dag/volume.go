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

func (idx *hashCacheIndex) Key() string {
	if idx.cfg == nil {
		return "#cache.nil"
	}
	mk := vegMemoKey(idx.node)
	if mk != "" {
		return fmt.Sprintf("#cache.%s", mk)
	}
	return fmt.Sprintf("#cache.%s", idx.cfg.Name)
}

func (d *Dag) hashCache(step cue.Value) (*dagger.CacheVolume, error) {
	d.mx.RLock()
	var cfg hashCacheConfig
	err := step.Decode(&cfg)
	d.mx.RUnlock()
	if err != nil {
		return nil, fmt.Errorf("while decoding hashCache: %w", err)
	}

	// index for query and create if not found
	idx := &hashCacheIndex{
		val: step,
		cfg: &cfg,
	}

	// lookup
	ia, ok := d.cat.Load(idx)
	if ok {
		ix := ia.(*hashCacheIndex)
		return ix.vol, nil
	}

	// load for realz
	idx.vol = d.dag.CacheVolume(cfg.Name)
	// memoize
	d.cat.Store(idx, idx)

	return idx.vol, nil
}

type stepTempConfig struct {
	Kind   string `json:"$kind"`
	Path   string `json:"path"`
	Size   int    `json:"size"`
	Expand bool   `json:"expand"`
}

func (d *Dag) stepTempHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	d.mx.RLock()
	var cfg stepTempConfig
	err := step.Decode(&cfg)
	d.mx.RUnlock()
	if err != nil {
		return nil, fmt.Errorf("while decoding stepTerm: %w", err)
	}

	c = c.WithMountedTemp(cfg.Path, dagger.ContainerWithMountedTempOpts{
		Size:   cfg.Size,
		Expand: cfg.Expand,
	})

	return c, nil
}

type stepMountConfig struct {
	Kind string `json:"$kind"`
	// args
	Path   string    `json:"path"`
	Source cue.Value `json:"source"`

	// opts (depending on source type?)
	Owner  string `json:"owner"`
	Expand bool   `json:"expand"`
	Mode   int    `json:"mode"`
}

func (d *Dag) stepMountHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	d.mx.RLock()
	var cfg stepMountConfig
	err := step.Decode(&cfg)
	d.mx.RUnlock()
	if err != nil {
		return c, err
	}

	// look for kind
	k := cfg.Source.LookupPath(cue.ParsePath("$kind"))
	if !k.Exists() {
		return c, fmt.Errorf("missing $kind in stepDir source: %v", step)
	}
	ks, _ := k.String()
	switch ks {
	case "#cache":
		cache, err := d.hashCache(cfg.Source)
		if err != nil {
			return nil, err
		}
		c = c.WithMountedCache(cfg.Path, cache, dagger.ContainerWithMountedCacheOpts{
			Owner:  cfg.Owner,
			Expand: cfg.Expand,
		})

	case "#secret":
		shh, err := d.hashSecret(cfg.Source)
		if err != nil {
			return nil, err
		}
		c = c.WithMountedSecret(cfg.Path, shh, dagger.ContainerWithMountedSecretOpts{
			Owner:  cfg.Owner,
			Expand: cfg.Expand,
			Mode:   cfg.Mode,
		})

	case "#file":
		file, _, err := d.hashFile(cfg.Source, false)
		if err != nil {
			return nil, err
		}
		c = c.WithMountedFile(cfg.Path, file, dagger.ContainerWithMountedFileOpts{
			Owner:  cfg.Owner,
			Expand: cfg.Expand,
		})

	case "#hostFile":
		file, _, err := d.HashHostFile(cfg.Source, false)
		if err != nil {
			return nil, err
		}
		c = c.WithMountedFile(cfg.Path, file, dagger.ContainerWithMountedFileOpts{
			Owner:  cfg.Owner,
			Expand: cfg.Expand,
		})

	case "#dir":
		dir, _, err := d.hashDir(cfg.Source, false)
		if err != nil {
			return nil, err
		}
		c = c.WithMountedDirectory(cfg.Path, dir, dagger.ContainerWithMountedDirectoryOpts{
			Owner:  cfg.Owner,
			Expand: cfg.Expand,
		})

	case "#hostDir":
		dir, _, err := d.HashHostDir(cfg.Source, false)
		if err != nil {
			return nil, err
		}
		c = c.WithMountedDirectory(cfg.Path, dir, dagger.ContainerWithMountedDirectoryOpts{
			Owner:  cfg.Owner,
			Expand: cfg.Expand,
		})

	case "#gitRepo":
		repo, rcfg, rerr := d.hashGitRepo(cfg.Source, false)
		if rerr == nil {
			var dir *dagger.Directory
			if rcfg != nil && rcfg.Ref != "" {
				dir = repo.Ref(rcfg.Ref).Tree()
			} else {
				dir = repo.Head().Tree()
			}
			c = c.WithMountedDirectory(cfg.Path, dir, dagger.ContainerWithMountedDirectoryOpts{
				Owner:  cfg.Owner,
				Expand: cfg.Expand,
			})
		} else {
			err = rerr
		}

	default:
		return c, fmt.Errorf("unsupported $kind in stepMount.source: %v", step)
	}

	return c, nil
}
