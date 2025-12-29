package dag

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"cuelang.org/go/cue"
	"dagger.io/dagger"
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
		"sync":        d.stepSyncHandler,
		"exec":        d.stepExecHandler,
		"user":        d.stepUserHandler,
		"workdir":     d.stepWorkdirHandler,
		"file":        d.stepFileHandler,
		"dir":         d.stepDirHandler,
		"mount":       d.stepMountHandler,
		"env":         d.stepEnvHandler,
		"envfile":     d.stepEnvfileHandler,
		"secret":      d.stepSecretHandler,
		"expose":      d.stepExposeHandler,
		"bindService": d.stepBindServiceHandler,
		"entrypoint":  d.stepEntrypointHandler,
		"args":        d.stepArgsHandler,
		"term":        d.stepTermHandler,
		"temp":        d.stepTempHandler,

		// // cataloged
		// "#hostImage":   d.hashHostImage,
		// "#hostFile":    d.hashHostFile,
		// "#hostDir":     d.hashHostDir,
		// "#hostService": d.hashHostService,
		// "#hostTunnel":  d.hashHostTunnel,
		// "#hostSocket":  d.hashHostSocket,
		// "#container":   d.hashContainer,
		// "#file":        d.hashFile,
		// "#dir":         d.hashDir,
		// "#service":     d.hashService,
		// "#temp":        d.hashTemp,
		// "#cache":       d.hashCache,
		// "#volume":      d.hashVolume,
		// "#secret":      d.hashSecret,
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

type stepFileConfig struct {
	Kind string `json:"$kind"`
	// args
	Path    string    `json:"path"`
	Content cue.Value `json:"content"`
	// opts
	Permissions int    `json:"permissions"`
	Owner       string `json:"owner"`
	Expand      bool   `json:"expand"`
}

func (d *Dag) stepFileHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	var cfg stepFileConfig
	err := step.Decode(&cfg)
	if err != nil {
		return c, err
	}

	var f *dagger.File
	switch ik := cfg.Content.IncompleteKind(); ik {
	case cue.StringKind:
		_, name := filepath.Split(cfg.Path)
		s, _ := cfg.Content.String()
		f = d.dag.File(name, s)
	case cue.StructKind:
		// look for kind
		k := cfg.Content.LookupPath(cue.ParsePath("$kind"))
		if !k.Exists() {
			return c, fmt.Errorf("missing $kind in struct file source: %v", step)
		}
		ks, _ := k.String()
		switch ks {
		case "#file":
			f, err = d.hashFile(cfg.Content)
		case "#hostFile":
			f, err = d.hashHostFile(cfg.Content)

		default:
			return c, fmt.Errorf("unsupported $kind in struct file source: %v", step)
		}
	}

	c = c.WithFile(cfg.Path, f)

	return c, nil
}

type stepDirConfig struct {
	Kind string `json:"$kind"`
	// args
	Path   string    `json:"path"`
	Source cue.Value `json:"source"`
	// opts
	Include   []string `json:"include"`
	Exclude   []string `json:"exclude"`
	Gitignore bool     `json:"gitignore"`
	Owner     string   `json:"owner"`
	Expand    bool     `json:"expand"`
}

func (d *Dag) stepDirHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	var cfg stepDirConfig
	err := step.Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("while decoding stepDir: %w", err)
	}

	var dir *dagger.Directory
	switch ik := cfg.Source.IncompleteKind(); ik {
	case cue.StructKind:
		// look for kind
		k := cfg.Source.LookupPath(cue.ParsePath("$kind"))
		if !k.Exists() {
			return c, fmt.Errorf("missing $kind in stepDir source: %v", step)
		}
		ks, _ := k.String()
		switch ks {
		case "#gitRepo":
			repo, err := d.hashGitRepo(cfg.Source)
			if err != nil {
				return nil, err
			}
			dir = repo.Head().Tree()

		case "#dir":
			dir, err = d.hashDir(cfg.Source)
			if err != nil {
				return nil, err
			}

		case "#hostDir":
			dir, err = d.hashHostDir(cfg.Source)
			if err != nil {
				return nil, err
			}

		case "#container":
			ctr, err := d.HashContainer(cfg.Source)
			if err != nil {
				return nil, err
			}
			dir = ctr.Directory(cfg.Path)

		case "#hostImage":
			ctr, err := d.HashHostImage(cfg.Source)
			if err != nil {
				return nil, err
			}
			dir = ctr.Directory(cfg.Path)

		default:
			return c, fmt.Errorf("unsupported $kind in stepDir source: %v", step)
		}

	default:
		return c, fmt.Errorf("unsupported stepDir value type: %v", step)
	}

	c = c.WithDirectory(cfg.Path, dir, dagger.ContainerWithDirectoryOpts{
		Include:   cfg.Include,
		Exclude:   cfg.Exclude,
		Gitignore: cfg.Gitignore,
		Owner:     cfg.Owner,
		Expand:    cfg.Expand,
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

	case "#file":
		file, err := d.hashFile(cfg.Source)
		if err != nil {
			return nil, err
		}
		c = c.WithMountedFile(cfg.Path, file, dagger.ContainerWithMountedFileOpts{
			Owner:  cfg.Owner,
			Expand: cfg.Expand,
		})

	case "#hostFile":
		file, err := d.hashHostFile(cfg.Source)
		if err != nil {
			return nil, err
		}
		c = c.WithMountedFile(cfg.Path, file, dagger.ContainerWithMountedFileOpts{
			Owner:  cfg.Owner,
			Expand: cfg.Expand,
		})

	case "#dir":
		dir, err := d.hashDir(cfg.Source)
		if err != nil {
			return nil, err
		}
		c = c.WithMountedDirectory(cfg.Path, dir, dagger.ContainerWithMountedDirectoryOpts{
			Owner:  cfg.Owner,
			Expand: cfg.Expand,
		})

	case "#hostDir":
		dir, err := d.hashHostDir(cfg.Source)
		if err != nil {
			return nil, err
		}
		c = c.WithMountedDirectory(cfg.Path, dir, dagger.ContainerWithMountedDirectoryOpts{
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

	ev := step.LookupPath(cue.ParsePath("$expand"))
	es, _ := ev.String()
	eb, _ := strconv.ParseBool(es)

	for k, v := range envs {
		if k != "$kind" && k != "$expand" {
			c = c.WithEnvVariable(k, v, dagger.ContainerWithEnvVariableOpts{
				Expand: eb,
			})
		}
	}
	return c, nil
}

func (d *Dag) stepEnvfileHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	return c, nil
}

func (d *Dag) stepSecretHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
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

func (d *Dag) stepArgsHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
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

func (d *Dag) stepTermHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
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
