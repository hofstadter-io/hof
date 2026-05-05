package dag

import (
	"fmt"
	"strings"

	"cuelang.org/go/cue"
	"dagger.io/dagger"
	"github.com/hofstadter-io/hof/lib/env"
)

type hashServiceConfig struct {
	Kind string `json:"$kind"`

	Name   string    `json:"name"`
	Source cue.Value `json:"source"`

	// various fields set when services are used
	Hostname string        `json:"hostname"`
	Ports    []portForward `json:"ports"`

	Args          []string `json:"args"`
	UseEntrypoint bool     `json:"useEntrypoint"`

	ExperimentalPrivilegedNesting bool `json:"experimentalPrivilegedNesting"`
	InsecureRootCapabilities      bool `json:"insecureRootCapabilities"`
	Expand                        bool `json:"expand"`
	NoInit                        bool `json:"noInit"`
}

type hashServiceIndex struct {
	node *env.Env
	val  cue.Value
	cfg  *hashServiceConfig
	svc  *dagger.Service
}

func (idx *hashServiceIndex) Key() string {
	if idx.cfg == nil {
		return "#service.nil"
	}
	mk := vegMemoKey(idx.node)
	if mk != "" {
		return fmt.Sprintf("#service.%s", mk)
	}
	return fmt.Sprintf("#service.%s", idx.cfg.Name)
}

func (d *Dag) HashService(step cue.Value, noCache bool) (*dagger.Service, *hashServiceConfig, error) {
	var err error
	step, err = d.ResolveShouldi(step)
	if err != nil {
		return nil, nil, err
	}
	if !step.Exists() {
		return nil, nil, fmt.Errorf("HashService: resolved to empty value")
	}

	d.mx.RLock()
	var cfg hashServiceConfig
	err = step.Decode(&cfg)
	d.mx.RUnlock()
	if err != nil {
		return nil, nil, fmt.Errorf("while decoding hashService: %w", err)
	}
	// fmt.Println("hashService.config", cfg)

	// index for query and create if not found
	idx := &hashServiceIndex{
		val: step,
		cfg: &cfg,
	}

	// lookup
	ia, ok := d.cat.Load(idx)
	if ok {
		ix := ia.(*hashServiceIndex)
		return ix.svc, idx.cfg, nil
	}

	// load for realz

	// look for kind
	var c *dagger.Container
	k := cfg.Source.LookupPath(cue.ParsePath("$kind"))
	if !k.Exists() {
		return nil, nil, fmt.Errorf("missing $kind in #service.source: %v", step)
	}
	ks, _ := k.String()
	switch ks {
	case "#container":
		c, err = d.HashContainer(cfg.Source, noCache)
		if err != nil {
			return nil, nil, err
		}
	case "#hostImage":
		c, err = d.HashHostImage(cfg.Source, noCache)
		if err != nil {
			return nil, nil, err
		}
	case "#dockerBuild":
		c, err = d.HashDockerBuild(cfg.Source, noCache)
		if err != nil {
			return nil, nil, err
		}
	case "#rootfs":
		dir, err := d.hashRootFS(cfg.Source, noCache)
		if err != nil {
			return nil, nil, err
		}
		c = d.dag.Container().WithRootfs(dir)
	default:
		return nil, nil, fmt.Errorf("unupported service.source $kind")
	}

	// fmt.Println("hashService.preparing")

	// prepare as-service inputs
	for _, p := range cfg.Ports {
		c = c.WithExposedPort(p.Backend, dagger.ContainerWithExposedPortOpts{
			Description: p.Name,
			Protocol:    dagger.NetworkProtocol(strings.ToUpper(p.Protocol)),
			// ExperimentalSkipHealthcheck: p.ExperimentalSkipHealthchecks,
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
	d.cat.Store(idx, idx)

	return idx.svc, idx.cfg, nil
}

type stepExposeConfig struct {
	Kind     string `json:"$kind"`
	Name     string `json:"name"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`

	ExperimentalSkipHealthchecks bool `json:"experimentalSkipHealthchecks"`
}

func (d *Dag) stepExposeHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	d.mx.RLock()
	var cfg stepExposeConfig
	err := step.Decode(&cfg)
	d.mx.RUnlock()
	if err != nil {
		return nil, fmt.Errorf("while decoding stepExpose: %w", err)
	}
	// fmt.Println("stepExpose.config", cfg)

	// c = c.WithExposedPort(cfg.Port)
	c = c.WithExposedPort(cfg.Port, dagger.ContainerWithExposedPortOpts{
		Description:                 cfg.Name,
		Protocol:                    dagger.NetworkProtocol(strings.ToUpper(cfg.Protocol)),
		ExperimentalSkipHealthcheck: cfg.ExperimentalSkipHealthchecks,
	})
	return c, nil
}

type stepBindServiceConfig struct {
	Kind  string `json:"$kind"`
	Alias string `json:"alias"`

	Service cue.Value `json:"service"`
}

func (d *Dag) stepBindServiceHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	d.mx.RLock()
	var cfg stepBindServiceConfig
	err := step.Decode(&cfg)
	d.mx.RUnlock()
	if err != nil {
		return nil, fmt.Errorf("while decoding stepEntrypoint: %w", err)
	}
	// fmt.Println("bindService.config", cfg)

	s, scfg, err := d.HashService(cfg.Service, false) // hmm, should we respect noCache here? probably not for a step binding
	if err != nil {
		return nil, err
	}
	if cfg.Alias == "" {
		cfg.Alias = scfg.Hostname
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
