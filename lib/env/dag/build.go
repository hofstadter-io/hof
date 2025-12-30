package dag

import (
	"fmt"

	"dagger.io/dagger"
	"github.com/hofstadter-io/hof/lib/env"
)

type kinder struct {
	Kind string `json:"$kind"`
}

func (d *Dag) Build(e *env.Env, noCache bool) (*dagger.Container, error) {
	d.noCache = noCache

	// it's probably wrong to assume this in general
	var k kinder
	err := e.Value.Decode(&k)
	if err != nil {
		return nil, err
	}

	switch k.Kind {
	case "#container":
		return d.HashContainer(e.Value)
	case "#hostImage":
		return d.HashHostImage(e.Value)
	case "#dockerBuild":
		return d.HashDockerBuild(e.Value)
	default:
		return nil, fmt.Errorf("unsupported build target: %v", k.Kind, e.Value)
	}
}

func (d *Dag) Service(e *env.Env, noCache bool) (*dagger.Service, *hashServiceConfig, error) {
	d.noCache = noCache

	// it's probably wrong to assume this in general
	var k kinder
	err := e.Value.Decode(&k)
	if err != nil {
		return nil, nil, err
	}

	switch k.Kind {
	case "#service":
		s, cfg, err := d.hashService(e.Value)

		return s, cfg, err
	default:
		return nil, nil, fmt.Errorf("unsupported build target: %v", k.Kind, e.Value)
	}
}
