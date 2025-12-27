package dag

import (
	"fmt"

	"cuelang.org/go/cue"
	"dagger.io/dagger"
	"github.com/hofstadter-io/hof/lib/env"
)

type hashServiceConfig struct {
	Kind string `json:"$kind"`

	Name     string        `json:"name"`
	Hostname string        `json:"hostname"`
	Ports    []portForward `json:"ports"`

	Source cue.Value `json:"source"`

	Args          []string `json:"args"`
	UseEntrypoint bool     `json:"useEntrypoint"`
	Expand        bool     `json:"expand"`
	NoInit        bool     `json:"noINit"`
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
	ports := []dagger.PortForward{}
	for _, p := range cfg.Ports {
		ports = append(ports, dagger.PortForward{
			Protocol: dagger.NetworkProtocol(p.Protocol),
			Frontend: p.Frontend,
			Backend:  p.Backend,
		})
	}

	// idx.svc = d.dag.Host().Service(ports, dagger.HostServiceOpts{
	// 	Host: cfg.Host,
	// })

	// memoize
	d.cat[idx] = idx

	return idx.svc, nil
}
