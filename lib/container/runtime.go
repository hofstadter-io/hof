package container

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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
	Name    Name
	Env     []string
	Replace bool
}

type RuntimeVersion struct {
	Client struct {
		Version    string
		APIVersion string
	}
	Server struct {
		Version       string
		APIVersion    string
		MinAPIVersion string
	}
}

type Runtime interface {
	Version(context.Context) (RuntimeVersion, error)
	Images(context.Context, Ref) ([]Image, error)
	Containers(context.Context, Name) ([]Container, error)
	Run(context.Context, Ref, Params) error
	Remove(context.Context, Name) error
	Pull(context.Context, Ref) error
}

func newRuntime(bin RuntimeBinary) runtime {
	return runtime{
		bin: bin,
	}
}

var _ Runtime = runtime{}

type runtime struct {
	bin RuntimeBinary
}

func (r runtime) exec(ctx context.Context, args ...string) (io.Reader, error) {
	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	cmd := exec.CommandContext(ctx, string(r.bin), args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("cmd run: %w\n%s", err, stderr.String())
	}

	return &stdout, nil
}

func (r runtime) execJSON(ctx context.Context, resp any, args ...string) error {
	args = append(args, []string{"--format", "json"}...)

	stdout, err := r.exec(ctx, args...)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	if err := json.NewDecoder(stdout).Decode(resp); err != nil {
		return fmt.Errorf("json decode: %w", err)
	}

	return nil
}

func (r runtime) Version(ctx context.Context) (RuntimeVersion, error) {
	var rv RuntimeVersion
	if err := r.execJSON(ctx, &rv, "version"); err != nil {
		return rv, fmt.Errorf("exec json: %w", err)
	}

	return rv, nil
}

func (r runtime) Containers(ctx context.Context, name Name) ([]Container, error) {
	var (
		arg        = fmt.Sprintf("name=%s", name)
		containers []Container
	)

	if err := r.execJSON(ctx, &containers, "container", "ls", arg); err != nil {
		return nil, fmt.Errorf("exec json: %w", err)
	}

	return containers, nil
}

func (r runtime) Images(ctx context.Context, ref Ref) ([]Image, error) {
	var (
		arg    = fmt.Sprintf("reference=%s*", ref)
		images []Image
	)

	if err := r.execJSON(ctx, &images, "image", "ls", arg); err != nil {
		return nil, fmt.Errorf("exec json: %w", err)
	}

	return images, nil
}

func (r runtime) Pull(ctx context.Context, ref Ref) error {
	if _, err := r.exec(ctx, "pull", string(ref)); err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	return nil
}

func (r runtime) Run(ctx context.Context, ref Ref, p Params) error {
	if p.Replace {
		if err := r.Remove(ctx, p.Name); err != nil {
			return fmt.Errorf("remove: %w", err)
		}
	}

	args := []string{
		"run",
		"--detach",
		"--name", string(p.Name),
	}

	for _, e := range p.Env {
		args = append(args, []string{"--env", e}...)
	}

	args = append(args, string(ref))

	if _, err := r.exec(ctx, args...); err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	return nil
}

func (r runtime) Remove(ctx context.Context, name Name) error {
	if _, err := r.exec(ctx, "rm", "--force", string(name)); err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	return nil
}
