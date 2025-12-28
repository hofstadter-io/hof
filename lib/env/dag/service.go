package dag

import (
	"fmt"
	"strings"

	"cuelang.org/go/cue"
	"dagger.io/dagger"
	"github.com/hofstadter-io/hof/lib/env"
)

type stepBindServiceConfig struct {
	Kind  string `json:"$kind"`
	Alias string `json:"alias"`

	Service cue.Value `json:"service"`
}

func (d *Dag) stepBindServiceHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	var cfg stepBindServiceConfig
	err := step.Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("while decoding stepEntrypoint: %w", err)
	}
	// fmt.Println("bindService.config", cfg)

	s, err := d.hashService(cfg.Service)
	if err != nil {
		return nil, err
	}
	hn, err := s.Hostname(d.ctx)
	if err != nil {
		fmt.Println("hn.error", err)
	}
	if hn == "" {
		s = s.WithHostname(cfg.Alias)
		hn, _ = s.Hostname(d.ctx)
		if err != nil {
			fmt.Println("hn.error.2", err)
		}
	}

	// fmt.Printf("buildService.attach: %q %q\n", cfg.Alias, hn)

	c = c.WithServiceBinding(cfg.Alias, s)
	return c, nil
}

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
	// fmt.Println("hashService.config", cfg)

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

	// fmt.Println("hashService.preparing")

	// prepare as-service inputs
	for _, p := range cfg.Ports {
		c = c.WithExposedPort(p.Port, dagger.ContainerWithExposedPortOpts{
			Description:                 p.Name,
			Protocol:                    dagger.NetworkProtocol(strings.ToUpper(p.Protocol)),
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

	// fmt.Println("hashService.done")

	// memoize
	d.cat[idx] = idx

	return idx.svc, nil
}
