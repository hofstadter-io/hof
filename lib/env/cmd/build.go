package cmd

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
	"github.com/codemodus/kace"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/env"
	"github.com/hofstadter-io/hof/lib/env/dag"
	"github.com/hofstadter-io/hof/lib/env/incept"
	"github.com/hofstadter-io/hof/lib/runtime"
)

func Build(args []string, rflags flags.RootPflagpole) error {
	dst := os.Getenv("DAGGER_SESSION_TOKEN")

	// check the runtime first before starting dagger
	R, err := prepRuntime(args, rflags)
	if err != nil {
		return err
	}

	// incept if we are not in dagger
	if dst == "" {
		// Prepare arguments for the inner command (hof env build [args...])
		inceptArgs := append([]string{"hof", "env", "build"}, args...)

		// Run incept
		err := incept.Incept(context.Background(), inceptArgs, &incept.InceptOptions{
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

	// do normal build stuff

	// this should be on the runtime probable?
	ctx := context.Background()
	os.Setenv("_EXPERIMENTAL_DAGGER_RUNNER_HOST", DAGGER_HOST)

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return fmt.Errorf("while connecting to dagger: %w", err)
	}

	for _, e := range R.Envs {
		// only building containers right now
		if e.Hof.Env.Kind != "container" {
			continue
		}
		// fmt.Println("-:", e.Hof.Env.Name, e.Hof.Env.Kind)
		// we just try to "build" everything unless there are args
		do := true
		if len(args) > 0 {
			do = false
			for _, a := range args {
				if a == e.Hof.Env.Name {
					do = true
					break
				}
			}
		}
		if do {
			err = build(R, client, ctx, e)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func build(R *runtime.Runtime, client *dagger.Client, ctx context.Context, e *env.Env) error {
	id := e.Hof.Metadata.ID
	if id == "" {
		id = kace.Snake(e.Hof.Metadata.Name) + " (auto)"
	}

	var c dag.Container
	err := e.Value.Decode(&c)
	if err != nil {
		return err
	}

	i, err := dag.Build(client, ctx, c)
	if err != nil {
		return err
	}

	i, err = i.Sync(ctx)
	if err != nil {
		return err
	}

	return nil
}

func inceptVeg(args []string, rflags flags.RootPflagpole) error {

	return nil
}
