package dag

import (
	"fmt"
	"sort"
	// "time"

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

func (idx *hashContainerIndex) Key() string {
	// default if no config, shouldn't really get here
	if idx.cfg == nil {
		return "#container.nil"
	}
	mk := vegMemoKey(idx.node)
	if mk != "" {
		return fmt.Sprintf("#container.%s", mk)
	}
	// return the name on this config as a last resort
	return fmt.Sprintf("#container.%s", idx.cfg.Name)
}

// TODO, change this to take a context (for nested OTEL spans)
func (d *Dag) HashContainer(val cue.Value, noCache bool) (*dagger.Container, error) {
	var err error
	val, err = d.ResolveShouldi(val)
	if err != nil {
		return nil, err
	}
	if !val.Exists() {
		return nil, fmt.Errorf("HashContainer: resolved to empty value")
	}

	d.mx.RLock()
	var cfg hashContainerConfig
	err = val.Decode(&cfg)
	d.mx.RUnlock()
	if err != nil {
		return nil, fmt.Errorf("while decoding HashContainer: %w", err)
	}

	// index for query and create if not found
	idx := &hashContainerIndex{
		val: val,
		cfg: &cfg,
	}

	// lookup
	ia, ok := d.cat.Load(idx)
	if ok {
		ix := ia.(*hashContainerIndex)
		return ix.ctr, nil
	}

	//
	// build for realz
	//
	// possibly from scratch
	c := d.dag.Container()

	// obtain the from
	switch fk := cfg.From.IncompleteKind(); fk {
	case cue.StringKind:
		s, _ := cfg.From.String()
		switch s {
		case "":
			return nil, fmt.Errorf("empty from string in:", val)

		case "scratch":
			// no-op, do nothing

		default:
			c = c.From(s)
		}

	case cue.StructKind:
		kv := cfg.From.LookupPath(cue.ParsePath("$kind"))
		if !kv.Exists() {
			return nil, fmt.Errorf("missing $kind in from: %v", cfg.From)
		}
		k, _ := kv.String()
		switch k {
		case "#container":
			c, err = d.HashContainer(cfg.From, noCache)
			if err != nil {
				return c, err
			}
		case "#hostImage":
			c, err = d.HashHostImage(cfg.From, noCache)
			if err != nil {
				return c, err
			}
		case "#dockerBuild":
			c, err = d.HashDockerBuild(cfg.From, noCache)
			if err != nil {
				return c, err
			}
		case "#rootfs":
			var dir *dagger.Directory
			dir, err = d.hashRootFS(cfg.From, noCache)
			if err != nil {
				return nil, err
			}
			c = c.WithRootfs(dir)
		}

	default:
		return nil, fmt.Errorf("unsupported from kind: %v", fk)
	}

	// possibly bust cache
	if noCache {
		// c = c.WithEnvVariable("BUSTED_CACHE", time.Now().Local().String())
		c = c.WithEnvVariable("BUSTED_CACHE", "womp womp, why don't you go fix it then?")
	}

	// add env before the container goes (most common, determinstic order)
	envs := cfg.Envs
	keys := make([]string, 0, len(envs))
	for k := range envs {
		if k != "$kind" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := envs[k]
		c = c.WithEnvVariable(k, v, dagger.ContainerWithEnvVariableOpts{
			Expand: true,
		})
	}

	// apply our steps
	c, err = d.addSteps(c, cfg.Steps)
	if err != nil {
		return c, fmt.Errorf("while adding steps: %w", err)
	}

	// add labels as we finish up
	for k, v := range cfg.Labels {
		c = c.WithAnnotation(k, v)
	}

	// save
	idx.ctr = c

	// memoize
	d.cat.Store(idx, idx)

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

type hashDockerBuildConfig struct {
	Kind string `json:"$kind"`
	Name string `json:"name"`

	Source     cue.Value `json:"source"`
	Dockerfile string    `json:"dockerfile"`
	Platform   string    `json:"platform"`
	Target     string    `json:"target"`

	BuildArgs map[string]string `json:"buildArgs"`
	Secrets   []cue.Value       `json:"secrets"`
	NoInit    bool              `json:"noInit"`
}

type hashDockerBuildIndex struct {
	node *env.Env
	val  cue.Value
	cfg  *hashDockerBuildConfig
	ctr  *dagger.Container
}

func (idx *hashDockerBuildIndex) Key() string {
	if idx.cfg == nil {
		return "#dockerBuild.nil"
	}
	mk := vegMemoKey(idx.node)
	if mk != "" {
		return fmt.Sprintf("#dockerBuild.%s", mk)
	}

	return fmt.Sprintf("#dockerBuild.%s", idx.cfg.Name)
}

func (d *Dag) HashDockerBuild(step cue.Value, noCache bool) (*dagger.Container, error) {
	var err error
	step, err = d.ResolveShouldi(step)
	if err != nil {
		return nil, err
	}
	if !step.Exists() {
		return nil, fmt.Errorf("HashDockerBuild: resolved to empty value")
	}

	d.mx.RLock()
	var cfg hashDockerBuildConfig
	err = step.Decode(&cfg)
	d.mx.RUnlock()
	if err != nil {
		return nil, fmt.Errorf("while decoding hashDockerBuild: %w", err)
	}
	// fmt.Println("hashDockerBuild.config", cfg)

	// index for query and create if not found
	idx := &hashDockerBuildIndex{
		val: step,
		cfg: &cfg,
	}

	// lookup
	ia, ok := d.cat.Load(idx)
	if ok {
		ix := ia.(*hashDockerBuildIndex)
		return ix.ctr, nil
	}

	// load for realz

	// look for kind
	var dir *dagger.Directory
	k := cfg.Source.LookupPath(cue.ParsePath("$kind"))
	if !k.Exists() {
		return nil, fmt.Errorf("missing $kind in #service.source: %v", step)
	}
	ks, _ := k.String()
	switch ks {
	case "#dir":
		dir, _, err = d.hashDir(cfg.Source, noCache)
		if err != nil {
			return nil, err
		}
	case "#hostDir":
		dir, _, err = d.HashHostDir(cfg.Source, noCache)
		if err != nil {
			return nil, err
		}
	case "#rootfs":
		dir, err = d.hashRootFS(cfg.Source, noCache)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unupported service.source $kind")
	}

	// fmt.Println("hashDockerBuild.preparing")

	args := []dagger.BuildArg{}
	for k, v := range cfg.BuildArgs {
		args = append(args, dagger.BuildArg{Name: k, Value: v})
	}
	secrets := []*dagger.Secret{}

	idx.ctr = dir.DockerBuild(dagger.DirectoryDockerBuildOpts{
		Dockerfile: cfg.Dockerfile,
		Platform:   dagger.Platform(cfg.Platform),
		BuildArgs:  args,
		Target:     cfg.Target,
		Secrets:    secrets,
		NoInit:     cfg.NoInit,
	})

	// fmt.Println("hashService.done")

	// memoize
	d.cat.Store(idx, idx)

	return idx.ctr, nil
}
