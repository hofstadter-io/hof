package dag

import (
	"fmt"
	"strings"

	"cuelang.org/go/cue"
	"dagger.io/dagger"
	"github.com/hashicorp/go-envparse"
)

type StepKind struct {
	Kind  string `json:"$kind"`
	Value cue.Value
}

type Step map[string]any

type stepHandler func(c *dagger.Container, step cue.Value) (*dagger.Container, error)
type stepHandlerMap map[string]stepHandler

func (d *Dag) makeStepHandlers() stepHandlerMap {
	return stepHandlerMap{
		// not cataloged
		"sync":    d.stepSyncHandler,
		"exec":    d.stepExecHandler,
		"user":    d.stepUserHandler,
		"workdir": d.stepWorkdirHandler,
		"file":    d.stepFileHandler,
		"dir":     d.stepDirHandler,
		"mount":   d.stepMountHandler,
		"env":     d.stepEnvHandler,
		"envfile": d.stepEnvfileHandler,
		"secret":  d.stepSecretHandler,
		// "secretvars":  d.stepSecretFileHandler,
		"expose":      d.stepExposeHandler,
		"bindService": d.stepBindServiceHandler,
		"entrypoint":  d.stepEntrypointHandler,
		"args":        d.stepDefaultArgsHandler,
		"term":        d.stepDefaultTermHandler,
		"terminal":    d.stepTerminalHandler,
		"temp":        d.stepTempHandler,
	}
}

func (d *Dag) stepSyncHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	return c.Sync(d.ctx)
}

type stepExecConfig struct {
	Kind string   `json:"$kind"`
	Args []string `json:"args"`

	// opts
	UseEntrypoint  bool   `json:"useEntrypoint"`
	Stdin          string `json:"stdin"`
	RedirectStdin  string `json:"redirectStdin"`
	RedirectStdout string `json:"redirectStdout"`
	RedirectStderr string `json:"redirectStderr"`
	Expect         string `json:"string"`

	ExperimentalPrivilegedNesting bool `json:"experimentalPrivilegedNesting"`
	InsecureRootCapabilities      bool `json:"insecureRootCapabilities"`
	Expand                        bool `json:"expand"`
	NoInit                        bool `json:"noInit"`
}

func (d *Dag) stepExecHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	var cfg stepExecConfig
	err := step.Decode(&cfg)
	if err != nil {
		return c, err
	}

	c = c.WithExec(cfg.Args, dagger.ContainerWithExecOpts{
		UseEntrypoint:  cfg.UseEntrypoint,
		Stdin:          cfg.Stdin,
		RedirectStdin:  cfg.RedirectStdin,
		RedirectStdout: cfg.RedirectStdout,
		RedirectStderr: cfg.RedirectStderr,
		Expect:         dagger.ReturnType(cfg.Expect),

		ExperimentalPrivilegedNesting: cfg.ExperimentalPrivilegedNesting,
		InsecureRootCapabilities:      cfg.InsecureRootCapabilities,
		Expand:                        cfg.Expand,
		NoInit:                        cfg.NoInit,
	})
	return c, nil
}

type stepUserConfig struct {
	Kind string `json:"$kind"`
	Name string `json:"name"`
}

func (d *Dag) stepUserHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	var cfg stepUserConfig
	err := step.Decode(&cfg)
	if err != nil {
		return c, err
	}

	c = c.WithUser(cfg.Name)
	return c, nil
}

type stepWorkdirConfig struct {
	Kind string `json:"$kind"`
	Path string `json:"path"`
}

func (d *Dag) stepWorkdirHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	var cfg stepWorkdirConfig
	err := step.Decode(&cfg)
	if err != nil {
		return c, err
	}

	c = c.WithWorkdir(cfg.Path)
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
	var cfg stepMountConfig
	err := step.Decode(&cfg)
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
		file, _, err := d.hashFile(cfg.Source)
		if err != nil {
			return nil, err
		}
		c = c.WithMountedFile(cfg.Path, file, dagger.ContainerWithMountedFileOpts{
			Owner:  cfg.Owner,
			Expand: cfg.Expand,
		})

	case "#hostFile":
		file, _, err := d.hashHostFile(cfg.Source)
		if err != nil {
			return nil, err
		}
		c = c.WithMountedFile(cfg.Path, file, dagger.ContainerWithMountedFileOpts{
			Owner:  cfg.Owner,
			Expand: cfg.Expand,
		})

	case "#dir":
		dir, _, err := d.hashDir(cfg.Source)
		if err != nil {
			return nil, err
		}
		c = c.WithMountedDirectory(cfg.Path, dir, dagger.ContainerWithMountedDirectoryOpts{
			Owner:  cfg.Owner,
			Expand: cfg.Expand,
		})

	case "#hostDir":
		dir, _, err := d.hashHostDir(cfg.Source)
		if err != nil {
			return nil, err
		}
		c = c.WithMountedDirectory(cfg.Path, dir, dagger.ContainerWithMountedDirectoryOpts{
			Owner:  cfg.Owner,
			Expand: cfg.Expand,
		})

	default:
		return c, fmt.Errorf("unsupported $kind in stepMount.source: %v", step)
	}

	return c, nil
}

func (d *Dag) stepEnvHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	var envs map[string]string
	err := step.Decode(&envs)
	if err != nil {
		return nil, fmt.Errorf("while decoding stepEnv: %w", err)
	}

	for k, v := range envs {
		if k != "$kind" {
			c = c.WithEnvVariable(k, v, dagger.ContainerWithEnvVariableOpts{
				Expand: true,
			})
		}
	}
	return c, nil
}

type stepEnvfileConfig struct {
	Kind string    `json:"$kind"`
	File cue.Value `json:"file"`
}

// needed for checking below, loop copied from module source
var envpair envparse.Pair

func (d *Dag) stepEnvfileHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {

	// DEV HACK
	// return c, nil

	var cfg stepEnvfileConfig
	err := step.Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("while decoding stepEnvfile: %w", err)
	}
	k := cfg.File.LookupPath(cue.ParsePath("$kind"))
	if !k.Exists() {
		return c, fmt.Errorf("missing $kind in stepEnvfile source: %v, got %v", step, k)
	}

	var file *dagger.File
	ks, _ := k.String()
	switch ks {
	case "#file":
		file, _, err = d.hashFile(cfg.File)

	case "#hostFile":
		file, _, err = d.hashHostFile(cfg.File)

	default:
		return c, fmt.Errorf("unsupported $kind in envfile.file: %v", step)
	}

	if err != nil {
		return nil, err
	}

	contents, err := file.Contents(d.ctx)
	if err != nil {
		return nil, err
	}

	r := strings.NewReader(contents)
	parser := envparse.New(r)
	for {
		kv, err := parser.Next()
		if err != nil {
			return nil, err
		}

		if kv == envpair {
			break
		}

		c = c.WithEnvVariable(kv.Key, kv.Val)
	}

	return c, nil
}

type stepExposeConfig struct {
	Kind     string `json:"$kind"`
	Name     string `json:"name"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`

	ExperimentalSkipHealthchecks bool `json:"experimentalSkipHealthchecks"`
}

func (d *Dag) stepExposeHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	var cfg stepExposeConfig
	err := step.Decode(&cfg)
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

type stepEntrypointConfig struct {
	Kind string   `json:"$kind"`
	Args []string `json:"args"`

	KeepDefaultArgs bool `json:"keepDefaultArgs"`
}

func (d *Dag) stepEntrypointHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	var cfg stepEntrypointConfig
	err := step.Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("while decoding stepEntrypoint: %w", err)
	}

	c = c.WithEntrypoint(cfg.Args, dagger.ContainerWithEntrypointOpts{
		KeepDefaultArgs: cfg.KeepDefaultArgs,
	})
	return c, nil
}

type stepArgsConfig struct {
	Kind string   `json:"$kind"`
	Args []string `json:"args"`
}

func (d *Dag) stepDefaultArgsHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	var cfg stepArgsConfig
	err := step.Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("while decoding stepArgs: %w", err)
	}

	c = c.WithDefaultArgs(cfg.Args)
	return c, nil
}

type stepTermConfig struct {
	Kind string   `json:"$kind"`
	Args []string `json:"args"`

	ExperimentalPrivilegedNesting bool `json:"experimentalPrivilegedNesting"`
	InsecureRootCapabilities      bool `json:"insecureRootCapabilities"`
}

func (d *Dag) stepDefaultTermHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	var cfg stepTermConfig
	err := step.Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("while decoding stepTerm: %w", err)
	}

	c = c.WithDefaultTerminalCmd(cfg.Args, dagger.ContainerWithDefaultTerminalCmdOpts{
		ExperimentalPrivilegedNesting: cfg.ExperimentalPrivilegedNesting,
		InsecureRootCapabilities:      cfg.InsecureRootCapabilities,
	})
	return c, nil
}

type stepTerminalConfig struct {
	Kind string   `json:"$kind"`
	Args []string `json:"args"`

	ExperimentalPrivilegedNesting bool `json:"experimentalPrivilegedNesting"`
	InsecureRootCapabilities      bool `json:"insecureRootCapabilities"`
}

func (d *Dag) stepTerminalHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	var cfg stepTerminalConfig
	err := step.Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("while decoding stepTerm: %w", err)
	}

	fmt.Println("Terminal:", cfg, step)

	c = c.Terminal(dagger.ContainerTerminalOpts{
		Cmd: cfg.Args,
	})

	return c, nil
}

type stepTempConfig struct {
	Kind   string `json:"$kind"`
	Path   string `json:"path"`
	Size   int    `json:"size"`
	Expand bool   `json:"expand"`
}

func (d *Dag) stepTempHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	var cfg stepTempConfig
	err := step.Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("while decoding stepTerm: %w", err)
	}

	c = c.WithMountedTemp(cfg.Path, dagger.ContainerWithMountedTempOpts{
		Size:   cfg.Size,
		Expand: cfg.Expand,
	})

	return c, nil
}
