package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"dagger.io/dagger"
	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/env/incept"
)

func Export(args []string, rflags flags.RootPflagpole, cflags flags.EnvPflagpole) error {

	// check the runtime first before starting dagger
	R, err := prepRuntime(nil, rflags)
	if err != nil {
		return err
	}

	// incept if we are not in dagger
	dst := os.Getenv("DAGGER_SESSION_TOKEN")
	if dst == "" {
		err := incept.Incept(context.Background(), os.Args, &incept.InceptOptions{
			Progress:    cflags.Progress,
			Interactive: cflags.Interactive,
			Stdout:      os.Stdout,
			Stderr:      os.Stderr,
			Stdin:       os.Stdin,
		})
		if err != nil {
			return fmt.Errorf("while running incept: %w", err)
		}

		return nil
	}

	// do normal stuff now that we are incepted

	// this should be on the runtime probable?
	ctx := context.Background()
	os.Setenv("_EXPERIMENTAL_DAGGER_RUNNER_HOST", DAGGER_HOST)

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return fmt.Errorf("while connecting to dagger: %w", err)
	}

	fmt.Println("exporting:")
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
				if strings.HasPrefix(e.Hof.Env.Name, a) {
					do = true
					break
				}
			}
		}
		if do {
			fmt.Println(" -", e.Hof.Env.Name)
			i, err := build(R, client, ctx, e, cflags.NoCache)
			if err != nil {
				return fmt.Errorf("while build'n image: %w", err)
			}

			err = i.ExportImage(ctx, fmt.Sprintf("%s:%s", e.Hof.Env.Name, "local"))

		}
	}

	return nil
}
