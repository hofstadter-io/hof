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

func Run(name string, rflags flags.RootPflagpole) error {
	dst := os.Getenv("DAGGER_SESSION_TOKEN")

	// incept if we are not in dagger
	if dst == "" {
		// Run incept
		err := incept.Incept(context.Background(), []string{"hof", "env", "run", name}, &incept.InceptOptions{
			Progress: "tty",
			Stdout:   os.Stdout,
			Stderr:   os.Stderr,
			Stdin:    os.Stdin,
		})
		if err != nil {
			return fmt.Errorf("while running incept: %w", err)
		}

		return nil
	}
	R, err := prepRuntime([]string{}, rflags)
	if err != nil {
		return err
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

	var c dag.Container
	err = e.Value.Decode(&c)
	if err != nil {
		return err
	}
	// fmt.Println(pretty.Formatter(c))

	ctx := context.Background()
	client, err := dagger.Connect(ctx)
	if err != nil {
		return fmt.Errorf("while connecting to dagger in build: %w", err)
	}

	i, err := dag.Build(client, ctx, c)
	if err != nil {
		return err
	}

	i, err = i.Terminal(dagger.ContainerTerminalOpts{}).Sync(ctx)
	if err != nil {
		return err
	}

	return nil
}
