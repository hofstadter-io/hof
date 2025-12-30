package cmd

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strings"

	"dagger.io/dagger"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/env/dag"
	"github.com/hofstadter-io/hof/lib/env/incept"
)

func Build(args []string, rflags flags.RootPflagpole, eflags flags.EnvPflagpole) error {
	args, cueargs := splitArgs(args)

	// check the runtime first before starting dagger
	R, err := prepRuntime(cueargs, rflags)
	if err != nil {
		return err
	}

	// incept if we are not in dagger
	dst := os.Getenv("DAGGER_SESSION_TOKEN")
	if dst == "" {
		err := incept.Incept(context.Background(), os.Args, &incept.InceptOptions{
			Progress:    eflags.Progress,
			Interactive: eflags.OnFailure,
			NoExit:      eflags.NoExit,
			Stdout:      os.Stdout,
			Stderr:      os.Stderr,
			Stdin:       os.Stdin,
		})
		if err != nil {
			return err
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
	d, _ := dag.NewClient(ctx, client)

	valid := []string{"container", "hostImage", "dockerBuild"}
	fmt.Println("building:")
	for _, e := range R.Envs {
		// fmt.Println("-:", e.Hof.Env.Name, e.Hof.Env.Kind)
		// only building containers right now
		if !slices.Contains(valid, e.Hof.Env.Kind) {
			continue
		}
		// we just try to "build" everything unless there are args
		do := true
		if len(args) > 0 {
			do = false
			for _, a := range args {
				if strings.HasPrefix(e.Hof.Env.Name, a) {
					do = true
					break
				}
			}
		}
		if do {
			fmt.Println(" -", e.Hof.Env.Name)

			i, err := d.Build(e, eflags.NoCache)
			if err != nil {
				fmt.Println("error:", err)
				return err
			}

			i, err = i.Sync(ctx)
			if err != nil {
				return err
			}

		}
		if err != nil {
			return err
		}
	}

	return nil
}
