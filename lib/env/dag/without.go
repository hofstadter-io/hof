package dag

import (
	"fmt"

	"cuelang.org/go/cue"
	"dagger.io/dagger"
)

func (d *Dag) stepWithoutDefaultArgsHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	return c.WithoutDefaultArgs(), nil
}

type stepWithoutDirectoryConfig struct {
	Path   string `json:"path"`
	Expand bool   `json:"expand"`
}

func (d *Dag) stepWithoutDirectoryHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	var cfg stepWithoutDirectoryConfig
	if err := step.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("while decoding stepWithoutDirectory: %w", err)
	}
	return c.WithoutDirectory(cfg.Path, dagger.ContainerWithoutDirectoryOpts{
		Expand: cfg.Expand,
	}), nil
}

type stepWithoutEntrypointConfig struct {
	KeepDefaultArgs bool `json:"keepDefaultArgs"`
}

func (d *Dag) stepWithoutEntrypointHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	var cfg stepWithoutEntrypointConfig
	if err := step.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("while decoding stepWithoutEntrypoint: %w", err)
	}
	return c.WithoutEntrypoint(dagger.ContainerWithoutEntrypointOpts{
		KeepDefaultArgs: cfg.KeepDefaultArgs,
	}), nil
}

type stepWithoutEnvVariableConfig struct {
	Name string `json:"name"`
}

func (d *Dag) stepWithoutEnvVariableHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	var cfg stepWithoutEnvVariableConfig
	if err := step.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("while decoding stepWithoutEnvVariable: %w", err)
	}
	return c.WithoutEnvVariable(cfg.Name), nil
}

type stepWithoutExposedPortConfig struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
}

func (d *Dag) stepWithoutExposedPortHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	var cfg stepWithoutExposedPortConfig
	if err := step.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("while decoding stepWithoutExposedPort: %w", err)
	}
	return c.WithoutExposedPort(cfg.Port, dagger.ContainerWithoutExposedPortOpts{
		Protocol: dagger.NetworkProtocol(cfg.Protocol),
	}), nil
}

type stepWithoutFileConfig struct {
	Path   string `json:"path"`
	Expand bool   `json:"expand"`
}

func (d *Dag) stepWithoutFileHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	var cfg stepWithoutFileConfig
	if err := step.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("while decoding stepWithoutFile: %w", err)
	}
	return c.WithoutFile(cfg.Path, dagger.ContainerWithoutFileOpts{
		Expand: cfg.Expand,
	}), nil
}

type stepWithoutFilesConfig struct {
	Paths  []string `json:"paths"`
	Expand bool     `json:"expand"`
}

func (d *Dag) stepWithoutFilesHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	var cfg stepWithoutFilesConfig
	if err := step.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("while decoding stepWithoutFiles: %w", err)
	}
	return c.WithoutFiles(cfg.Paths, dagger.ContainerWithoutFilesOpts{
		Expand: cfg.Expand,
	}), nil
}

type stepWithoutLabelConfig struct {
	Name string `json:"name"`
}

func (d *Dag) stepWithoutLabelHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	var cfg stepWithoutLabelConfig
	if err := step.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("while decoding stepWithoutLabel: %w", err)
	}
	return c.WithoutLabel(cfg.Name), nil
}

type stepWithoutMountConfig struct {
	Path   string `json:"path"`
	Expand bool   `json:"expand"`
}

func (d *Dag) stepWithoutMountHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	var cfg stepWithoutMountConfig
	if err := step.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("while decoding stepWithoutMount: %w", err)
	}
	return c.WithoutMount(cfg.Path, dagger.ContainerWithoutMountOpts{
		Expand: cfg.Expand,
	}), nil
}

type stepWithoutRegistryAuthConfig struct {
	Address string `json:"address"`
}

func (d *Dag) stepWithoutRegistryAuthHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	var cfg stepWithoutRegistryAuthConfig
	if err := step.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("while decoding stepWithoutRegistryAuth: %w", err)
	}
	return c.WithoutRegistryAuth(cfg.Address), nil
}

type stepWithoutSecretVariableConfig struct {
	Name string `json:"name"`
}

func (d *Dag) stepWithoutSecretVariableHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	var cfg stepWithoutSecretVariableConfig
	if err := step.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("while decoding stepWithoutSecretVariable: %w", err)
	}
	return c.WithoutSecretVariable(cfg.Name), nil
}

type stepWithoutUnixSocketConfig struct {
	Path   string `json:"path"`
	Expand bool   `json:"expand"`
}

func (d *Dag) stepWithoutUnixSocketHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	var cfg stepWithoutUnixSocketConfig
	if err := step.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("while decoding stepWithoutUnixSocket: %w", err)
	}
	return c.WithoutUnixSocket(cfg.Path, dagger.ContainerWithoutUnixSocketOpts{
		Expand: cfg.Expand,
	}), nil
}

func (d *Dag) stepWithoutUserHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	return c.WithoutUser(), nil
}

func (d *Dag) stepWithoutWorkdirHandler(c *dagger.Container, step cue.Value) (*dagger.Container, error) {
	return c.WithoutWorkdir(), nil
}
