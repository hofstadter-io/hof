package container

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/hofstadter-io/hof/lib/yagu"
)

type RuntimeBinary string

const (
	RuntimeBinaryNone    RuntimeBinary = "none"
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

type Runtime interface {
	Binary() string
	Version(context.Context) (RuntimeVersion, error)
	Images(context.Context, Ref) ([]Image, error)
	Containers(context.Context, Name) ([]Container, error)
	Run(context.Context, Ref, Params) error
	Remove(context.Context, Name) error
	Pull(context.Context, Ref) error
}

func newRuntime(bin RuntimeBinary) runtime {
	_, debug := os.LookupEnv("HOF_CONTAINER_DEBUG")

	return runtime{
		bin:   bin,
		debug: debug,
	}
}

var rt Runtime

type runtime struct {
	bin   RuntimeBinary
	debug bool
}

func (r runtime) Binary() string {
	return string(r.bin)
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

	if r.debug {
		fmt.Println(cmd.String())
		fmt.Println(stdout.String())
	}

	return &stdout, nil
}

func (r runtime) execJSON(ctx context.Context, resp any, args ...string) error {
	stdout, err := r.exec(ctx, args...)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	if err := json.NewDecoder(stdout).Decode(resp); err != nil {
		return fmt.Errorf("json decode: %w", err)
	}

	return nil
}

func (r runtime) Containers(ctx context.Context, name Name) ([]Container, error) {
	args := []string{
		"container", "ls", "-a",
		"--filter", fmt.Sprintf("name=%s", name),
		"--format", "{{ json . }}",
	}

	stdout, err := r.exec(ctx, args...)
	if err != nil {
		return nil, fmt.Errorf("exec: %w", err)
	}

	containers, err := ndjson[Container](stdout)
	if err != nil {
		return nil, fmt.Errorf("ndjson: %w", err)
	}

	// fmt.Println("CONTAINERS:", containers)

	// HACK: various fixes...
	for i, _ := range containers {
		c := containers[i]
		if c.State == "" && c.Status != "" {
			c.State = c.Status
		}
		// fix status string
		if strings.HasPrefix(c.State, "Up") {
			c.State = "running"
		}
		
		containers[i] = c
	}

	return containers, nil
}

func (r runtime) Images(ctx context.Context, ref Ref) ([]Image, error) {
	args := []string{
		"image", "ls",
		"--filter", fmt.Sprintf("reference=%s*", ref),
		"--format", "{{ json . }}",
	}

	stdout, err := r.exec(ctx, args...)
	if err != nil {
		return nil, fmt.Errorf("exec: %w", err)
	}

	imgs, err := ndjson[Image](stdout)
	if err != nil {
		return nil, fmt.Errorf("ndjson: %w", err)
	}

	// fmt.Println("IMAGES:", imgs)

	// need to process images here, merge Tag into RepoTags by Repository
	m := map[string]Image{}
	for _, img := range imgs {
		// HACK podman fix for different output
		if img.Repository == "" {
			if len(img.Names) > 0 {
				n := img.Names[0]
				p := strings.Index(n, ":")
				img.Repository = n[:p]
				img.Tag = n[p+1:]
			}
		}
		i, ok := m[img.Repository]
		if !ok { 
			i = img
		}
		if img.Tag != "" {
			i.RepoTags = append(i.RepoTags, img.Tag)
		}
		m[img.Repository] = i
	}

	// build up return list
	ret := []Image{}
	for _, i := range m {
		ret = append(ret, i)
	}

	return ret, nil
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

	port, err := yagu.GetFreePort()
	if err != nil {
		return fmt.Errorf("while getting a free port: %w", err)
	}

	args := []string{
		"run",
		"-p",
		fmt.Sprintf("%d:3000", port),
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

func ndjson[T any](r io.Reader) ([]T, error) {
	var (
		ts []T
		R  = bufio.NewReader(r)
	)

	bs, err := io.ReadAll(R)
	if err != nil {
		return nil, fmt.Errorf("readall []: %w", err)
	}
	// fmt.Println("bs:", len(bs), string(bs))

	//// some runtimes return an array
	if bytes.HasPrefix(bs, []byte{'['}) {
		if err := json.Unmarshal(bs, &ts); err != nil {
			return nil, fmt.Errorf("json unmarshal []: %w", err)
		}
	} else if len(bs) > 0 {
		// fmt.Println("GOT HERE")
		// other runtimes return an ndjson
		S  := bufio.NewScanner(bytes.NewReader(bs))
		for S.Scan() {
			var t T
			if err := json.Unmarshal(S.Bytes(), &t); err != nil {
				return nil, fmt.Errorf("json unmarshal: %w", err)
			}
			// fmt.Println("t:", t)
			ts = append(ts, t)
		}

		if err := S.Err(); err != nil {
			return nil, fmt.Errorf("scanner: %w", err)
		}
	}

	return ts, nil
}
