package container

import (
	"context"
	"fmt"

	"github.com/hofstadter-io/hof/lib/yagu"
)

type Params struct {
	Name    Name
	Env     []string
	Replace bool

	// TODO, add more docker stuff here
}

func (r runtime) Run(ctx context.Context, ref Ref, p *Params) error {
	if p.Replace {
		if err := r.Remove(ctx, p.Name); err != nil {
			return fmt.Errorf("remove: %w", err)
		}
	}

	port, err := yagu.GetFreePort()
	if err != nil {
		return fmt.Errorf("while getting a free port: %w", err)
	}

	// TODO, build up args better based on Params

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
