package dag

import (
	"fmt"
	"os/exec"
	"strings"

	"cuelang.org/go/cue"
	"dagger.io/dagger"
	"github.com/hofstadter-io/hof/lib/env"
)

type hashHostExecConfig struct {
	Kind string   `json:"$kind"`
	Args []string `json:"args"`

	// opts
	UseEntrypoint  bool   `json:"useEntrypoint"`
	Stdin          string `json:"stdin"`
	RedirectStdin  string `json:"redirectStdin"`
	RedirectStdout string `json:"redirectStdout"`
	RedirectStderr string `json:"redirectStderr"`
	Expect         string `json:"string"`
}

type hashHostExecIndex struct {
	node *env.Env
	val  cue.Value
	cfg  *hashHostExecConfig
	cmd  *exec.Cmd
}

func (idx *hashHostExecIndex) Key() string {
	if idx.cfg == nil {
		return "#service.nil"
	}
	mk := vegMemoKey(idx.node)
	if mk != "" {
		return fmt.Sprintf("#service.%s", mk)
	}
	return fmt.Sprintf("#service.%s", strings.Join(idx.cfg.Args, " "))
}

func (d *Dag) hashExecHandler(step cue.Value) (*exec.Cmd, error) {
	var cfg hashHostExecConfig
	err := step.Decode(&cfg)
	if err != nil {
		return nil, err
	}
	// fmt.Println("hashHostDec.config", cfg)

	// index for query and create if not found
	idx := &hashHostExecIndex{
		val: step,
		cfg: &cfg,
	}

	// lookup
	ia, ok := d.cat[idx]
	if ok {
		ix := ia.(*hashHostExecIndex)
		return ix.cmd, nil
	}

	return idx.cmd, nil
}

type stepExecConfig struct {
	Kind string   `json:"$kind"`
	Args []string `json:"args"`

	// opts
	Workdir        string `json:"workdir"`
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

func (d *Dag) stepSyncHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	return c.Sync(d.ctx)
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

type stepDefaultArgsConfig struct {
	Kind string   `json:"$kind"`
	Args []string `json:"args"`
}

func (d *Dag) stepDefaultArgsHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	var cfg stepDefaultArgsConfig
	err := step.Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("while decoding stepArgs: %w", err)
	}

	c = c.WithDefaultArgs(cfg.Args)
	return c, nil
}

type stepDefaultTermConfig struct {
	Kind string   `json:"$kind"`
	Args []string `json:"args"`

	ExperimentalPrivilegedNesting bool `json:"experimentalPrivilegedNesting"`
	InsecureRootCapabilities      bool `json:"insecureRootCapabilities"`
}

func (d *Dag) stepDefaultTermHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	var cfg stepDefaultTermConfig
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
