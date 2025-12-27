package cmd

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/env"
	"github.com/hofstadter-io/hof/lib/env/dag"
	"github.com/hofstadter-io/hof/lib/env/incept"
)

func Run(name string, rflags flags.RootPflagpole, cflags flags.EnvPflagpole) error {
	R, err := prepRuntime([]string{}, rflags)
	if err != nil {
		return err
	}

	// incept if we are not in dagger
	dst := os.Getenv("DAGGER_SESSION_TOKEN")
	if dst == "" {
		// Run incept
		err := incept.Incept(context.Background(), os.Args, &incept.InceptOptions{
			Progress:    cflags.Progress,
			Interactive: true,
			Stdout:      os.Stdout,
			Stderr:      os.Stderr,
			Stdin:       os.Stdin,
		})
		if err != nil {
			return fmt.Errorf("while running incept: %w", err)
		}

		return nil
	}

	var e *env.Env
	for _, ee := range R.Envs {
		// only building containers right now
		if ee.Hof.Env.Kind == "container" && name == ee.Hof.Env.Name {
			e = ee
			break
		}
	}

	if e == nil {
		return fmt.Errorf("failed to find env %q", name)
	}

	ctx := context.Background()
	client, err := dagger.Connect(ctx)
	if err != nil {
		return fmt.Errorf("while connecting to dagger in build: %w", err)
	}
	d, _ := dag.NewClient(ctx, client)

	i, err := d.Build(e, cflags.NoCache)
	if err != nil {
		return err
	}

	i, err = i.Terminal(dagger.ContainerTerminalOpts{}).Sync(ctx)
	if err != nil {
		return err
	}

	return nil
}
