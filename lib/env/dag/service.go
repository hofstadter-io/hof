package dag

import (
	"fmt"

	"cuelang.org/go/cue"
	"dagger.io/dagger"
	"github.com/hofstadter-io/hof/lib/env"
)

type hashServiceConfig struct {
	Kind string `json:"$kind"`

	Name     string             `json:"name"`
	Hostname string             `json:"hostname"`
	Ports    []stepExposeConfig `json:"ports"`

	Source cue.Value `json:"source"`

	Args          []string `json:"args"`
	UseEntrypoint bool     `json:"useEntrypoint"`

	ExperimentalPrivilegedNesting bool `json:"experimentalPrivilegedNesting"`
	InsecureRootCapabilities      bool `json:"insecureRootCapabilities"`
	Expand                        bool `json:"expand"`
	NoInit                        bool `json:"noINit"`
}

type hashServiceIndex struct {
	node *env.Env
	val  cue.Value
	cfg  *hashServiceConfig
	svc  *dagger.Service
}

func (idx *hashServiceIndex) Key() string {
	if idx.cfg == nil {
		return "service.nil"
	}
	return fmt.Sprintf("service.%s", idx.cfg.Name)
}

func (d *Dag) hashService(step cue.Value) (*dagger.Service, error) {
	var cfg hashServiceConfig
	err := step.Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("while decoding hashService: %w", err)
	}

	// index for query and create if not found
	idx := &hashServiceIndex{
		val: step,
		cfg: &cfg,
	}

	// lookup
	ia, ok := d.cat[idx]
	if ok {
		ix := ia.(*hostServiceIndex)
		return ix.svc, nil
	}

	// load for realz

	// look for kind
	var c *dagger.Container
	k := cfg.Source.LookupPath(cue.ParsePath("$kind"))
	if !k.Exists() {
		return nil, fmt.Errorf("missing $kind in #service.source: %v", step)
	}
	ks, _ := k.String()
	switch ks {
	case "#container":
		c, err = d.hashContainer(cfg.Source)
		if err != nil {
			return nil, err
		}
	case "#hostImage":
		c, err = d.hashHostImage(cfg.Source)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unupported service.source $kind")
	}

	// prepare as-service inputs
	for _, p := range cfg.Ports {
		c = c.WithExposedPort(p.Port, dagger.ContainerWithExposedPortOpts{
			Description:                 p.Name,
			Protocol:                    dagger.NetworkProtocol(p.Protocol),
			ExperimentalSkipHealthcheck: p.ExperimentalSkipHealthchecks,
		})
	}
	idx.svc = c.AsService(dagger.ContainerAsServiceOpts{
		Args:                          cfg.Args,
		UseEntrypoint:                 cfg.UseEntrypoint,
		ExperimentalPrivilegedNesting: cfg.ExperimentalPrivilegedNesting,
		InsecureRootCapabilities:      cfg.InsecureRootCapabilities,
		Expand:                        cfg.Expand,
		NoInit:                        cfg.NoInit,
	})

	// memoize
	d.cat[idx] = idx

	return idx.svc, nil
}
