package dag

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"dagger.io/dagger"
)

type Step map[string]any

type Container struct {
	From   any    `json:"from"`
	Steps  []Step `json:"steps"`
	Labels map[string]string
}

func Build(client *dagger.Client, ctx context.Context, c Container, noCache bool) (*dagger.Container, error) {

	r := client.Container()

	// handle FROM first, without cache busting
	// if you want that, cache bust the from image manually
	//   otherwise we end up from images N-1 times
	switch t := c.From.(type) {
	case string:
		r = r.From(t)

	case Container:
		b, err := Build(client, ctx, t, false)
		if err != nil {
			return b, fmt.Errorf("while building the from image: %v %v", c, t)
		}
		r = b

	case map[string]any:
		m, err := mapToContainer(t)
		if err != nil {
			return nil, fmt.Errorf("while parsing the from image: %v %v", c, t)
		}
		b, err := Build(client, ctx, m, false)
		if err != nil {
			return b, fmt.Errorf("while building the from image: %v %v", c, t)
		}
		r = b
	default:
		return nil, fmt.Errorf("uknown from kind %v", t)
	}

	// possibly bust cache
	if noCache {
		r = r.WithEnvVariable("BUSTED_CACHE", time.Now().String())
	}

	// apply our steps
	r, err := addSteps(client, ctx, r, c.Steps)
	if err != nil {
		return r, fmt.Errorf("while adding steps: %w", err)
	}

	for k, v := range c.Labels {
		r = r.WithAnnotation(k, v)
	}
	// todo, set labels

	return r, nil
}

func addSteps(client *dagger.Client, ctx context.Context, c *dagger.Container, steps []Step) (*dagger.Container, error) {
	for i, s := range steps {
		var err error
		// todo, if step is an []any, assume nested steps, this should make the CUE simpler

		c, err = addStep(client, ctx, c, s)
		if err != nil {
			return c, fmt.Errorf("while adding step %d: %w", i, err)
		}
	}

	return c, nil
}

func addStep(client *dagger.Client, ctx context.Context, c *dagger.Container, s Step) (*dagger.Container, error) {

	// fmt.Printf("    %#+v\n", pretty.Formatter(s))

	kind, ok := s["$kind"]
	if !ok {
		return nil, fmt.Errorf("missing kind in step")
	}

	switch kind {
	case "sync":
		var err error
		c, err = c.Sync(ctx)
		if err != nil {
			return c, err
		}

	case "exec":
		args, ok := s["args"]
		if !ok {
			return c, fmt.Errorf("missing args in Exec")
		}
		as := make([]string, 0, len(args.([]any)))
		for _, a := range args.([]any) {
			as = append(as, a.(string))
		}
		c = c.WithExec(as, dagger.ContainerWithExecOpts{
			Expand: true,
		})

	case "user":
		name, ok := s["name"]
		if !ok {
			return c, fmt.Errorf("missing name in User")
		}
		c = c.WithUser(name.(string))

	case "workdir":
		path, ok := s["path"]
		if !ok {
			return c, fmt.Errorf("missing path in Workdir")
		}
		p, ok := path.(string)
		if !ok {
			return c, fmt.Errorf("path in Workdir is not a string")
		}
		c = c.WithWorkdir(p)

	case "file":
		path, ok := s["path"]
		if !ok {
			return c, fmt.Errorf("missing path in File")
		}
		content, ok := s["content"]
		if !ok {
			return c, fmt.Errorf("missing content in File")
		}

		p := path.(string)

		_, name := filepath.Split(p)

		f := client.File(name, content.(string))

		c = c.WithFile(p, f)

	case "env":
		for k, v := range s {
			if k != "$kind" {
				c = c.WithEnvVariable(k, v.(string), dagger.ContainerWithEnvVariableOpts{Expand: true})
			}
		}

	case "entrypoint":
		args, ok := s["args"]
		if !ok {
			return c, fmt.Errorf("missing args in Entrypoint")
		}
		as := make([]string, 0, len(args.([]any)))
		for _, a := range args.([]any) {
			as = append(as, a.(string))
		}
		c = c.WithEntrypoint(as)

	case "args":
		args, ok := s["args"]
		if !ok {
			return c, fmt.Errorf("missing args in Args")
		}
		as := make([]string, 0, len(args.([]any)))
		for _, a := range args.([]any) {
			as = append(as, a.(string))
		}
		c = c.WithDefaultArgs(as)

	case "term":
		args, ok := s["args"]
		if !ok {
			return c, fmt.Errorf("missing args in Term")
		}
		as := make([]string, 0, len(args.([]any)))
		for _, a := range args.([]any) {
			as = append(as, a.(string))
		}
		c = c.WithDefaultTerminalCmd(as)

	default:
		return c, fmt.Errorf("unknown kind %q", kind)
	}

	return c, nil
}

func mapToContainer(m map[string]any) (Container, error) {
	var c Container
	c.From = m["from"]

	steps := m["steps"].([]any)
	c.Steps = make([]Step, 0, len(steps))
	for _, s := range steps {
		c.Steps = append(c.Steps, Step(s.(map[string]any)))
	}

	labels := m["labels"].(map[string]any)
	c.Labels = make(map[string]string)
	for k, v := range labels {
		c.Labels[k] = v.(string)
	}

	return c, nil
}
