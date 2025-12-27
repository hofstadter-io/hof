package dag

import (
	"fmt"
	"time"

	"cuelang.org/go/cue"
	"dagger.io/dagger"
	"github.com/hofstadter-io/hof/lib/env"
)

type hashContainerConfig struct {
	Name   string            `json:"name"`
	From   cue.Value         `json:"from"`
	Envs   map[string]string `json:"envs"`
	Steps  []cue.Value       `json:"steps"`
	Labels map[string]string `json:"labels"`
}

type hashContainerIndex struct {
	node *env.Env
	val  cue.Value
	cfg  *hashContainerConfig
	ctr  *dagger.Container
}

func (h *hashContainerIndex) Key() string {
	if h.cfg == nil {
		return "#container.nil"
	}
	return fmt.Sprintf("#container.%s", h.cfg.Name)
}

func (d *Dag) hashContainer(step cue.Value) (*dagger.Container, error) {
	var cfg hashContainerConfig
	err := step.Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("while decoding hashHostDir: %w", err)
	}

	// index for query and create if not found
	idx := &hashContainerIndex{
		val: step,
		cfg: &cfg,
	}

	// lookup
	ia, ok := d.cat[idx]
	if ok {
		ix := ia.(*hashContainerIndex)
		return ix.ctr, nil
	}

	//
	// build for realz
	//
	c := d.dag.Container()

	// from
	switch fk := cfg.From.IncompleteKind(); fk {
	case cue.StringKind:
		s, _ := cfg.From.String()
		c = c.From(s)

	case cue.StructKind:
		kv := cfg.From.LookupPath(cue.ParsePath("$kind"))
		if !kv.Exists() {
			return nil, fmt.Errorf("missing $kind in from: %v", cfg.From)
		}
		k, _ := kv.String()
		switch k {
		case "#container":
			c, err = d.hashContainer(cfg.From)
			if err != nil {
				return c, err
			}
		case "#hostImage":
			c, err = d.hashHostImage(cfg.From)
			if err != nil {
				return c, err
			}
		}

	default:
		return nil, fmt.Errorf("unsupported from kind: %v", fk)
	}

	// possibly bust cache
	if d.noCache {
		c = c.WithEnvVariable("BUSTED_CACHE", time.Now().Local().String())
	}

	// apply our steps
	c, err = d.addSteps(c, cfg.Steps)
	if err != nil {
		return c, fmt.Errorf("while adding steps: %w", err)
	}

	for k, v := range cfg.Labels {
		c = c.WithAnnotation(k, v)
	}

	// save
	idx.ctr = c

	// memoize
	d.cat[idx] = idx

	return idx.ctr, nil
}

func (d *Dag) addSteps(c *dagger.Container, steps []cue.Value) (*dagger.Container, error) {
	var err error
	for i, s := range steps {
		switch ik := s.IncompleteKind(); ik {
		case cue.ListKind:
			it, _ := s.List()
			l, _ := s.Len().Int64()
			vals := make([]cue.Value, 0, l)
			for it.Next() {
				vals = append(vals, it.Value())
			}
			c, err = d.addSteps(c, vals)
			if err != nil {
				return c, fmt.Errorf("during step(%d)[%s]: %w", i, ik, err)
			}

		case cue.StructKind:
			kv := s.LookupPath(cue.ParsePath("$kind"))
			// fmt.Println("   -", i, kv)
			if !kv.Exists() {
				return c, fmt.Errorf("missing $kind on step: %v %v", s.Path(), s)
			}

			k, err := kv.String()
			if err != nil {
				return c, fmt.Errorf("$kind should be a string, we should never get here unless you are not using the schemas, got: %v", s)
			}

			h, ok := d.hdl[k]
			if !ok {
				return c, fmt.Errorf("unknown step(%d)[%s]: %v", i, k, s)
			}

			c, err = h(c, s)
			if err != nil {
				return c, fmt.Errorf("while adding step(%d)[%s@%v]: %w", i, k, s.Path(), err)
			}

		}
	}

	return c, nil
}
