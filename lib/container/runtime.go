package container

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
)

type RuntimeBinary string

const (
	RuntimeBinaryNerdctl RuntimeBinary = "nerdctl"
	RuntimeBinaryPodman  RuntimeBinary = "podman"
	RuntimeBinaryDocker  RuntimeBinary = "docker"
)

type (
	Ref  string
	Name string
)

type Params struct {
	Name    string
	Env     []string
	Replace bool
}

type Runtime interface {
	GetImages(context.Context, Ref) ([]Image, error)
	GetContainers(context.Context, Name) ([]Container, error)
	Start(context.Context, Ref, Params) error
	Stop(context.Context, string) error
	Pull(context.Context, Ref) error
}

func newRuntime(bin RuntimeBinary) runtime {
	return runtime{
		bin: bin,
	}
}

type runtime struct {
	bin RuntimeBinary
}

func (r runtime) exec(ctx context.Context, resp any, args ...string) error {
	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	args = append(args, []string{"--format", "json"}...)

	cmd := exec.CommandContext(ctx, string(r.bin), args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("cmd run: %w\n%s", err, stderr.String())
	}

	if err := json.NewDecoder(&stdout).Decode(resp); err != nil {
		return fmt.Errorf("json decode: %w", err)
	}

	return nil
}

func (r runtime) GetContainers(ctx context.Context, name Name) ([]Container, error) {
	var (
		arg        = fmt.Sprintf("name=%s", name)
		containers []Container
	)

	if err := r.exec(ctx, &containers, "container", "ls", arg); err != nil {
		return nil, fmt.Errorf("exec: %w", err)
	}

	return containers, nil
}

func (r runtime) GetImages(ctx context.Context, ref Ref) ([]Image, error) {
	var (
		arg    = fmt.Sprintf("reference=%s*", ref)
		images []Image
	)

	if err := r.exec(ctx, &images, "image", "ls", arg); err != nil {
		return nil, fmt.Errorf("exec: %w", err)
	}

	return images, nil
}

func (runtime) Pull(Ref) error {
	panic("unimplemented")
}

func (runtime) Start(Ref, Params) error {
	panic("unimplemented")
}

func (runtime) Stop(string) error {
	panic("unimplemented")
}
