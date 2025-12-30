package cmd

import (
	"context"
	"fmt"
	"os"
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
			Verbose:     rflags.Verbosity,
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

	fmt.Println("building:")
	for _, e := range R.Envs {
		// fmt.Println("-:", e.Hof.Env.Name, e.Hof.Env.Kind)
		// only building containers right now
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
		if !do {
			continue
		}

		switch e.Hof.Env.Kind {
		case "container", "hostImage", "dockerBuile":
			fmt.Printf(" - %s (%s)\n", e.Hof.Env.Name, e.Hof.Env.Kind)
			i, err := d.Container(e, eflags.NoCache)
			if err != nil {
				fmt.Println("error:", err)
				return err
			}

			i, err = i.Sync(ctx)
			if err != nil {
				return err
			}
		case "file", "hostFile":
			fmt.Printf(" - %s (%s)\n", e.Hof.Env.Name, e.Hof.Env.Kind)
			i, _, err := d.File(e, eflags.NoCache)
			if err != nil {
				fmt.Println("error:", err)
				return err
			}

			i, err = i.Sync(ctx)
			if err != nil {
				return err
			}
		case "dir", "hostDir", "gitRepo":
			fmt.Printf(" - %s (%s)\n", e.Hof.Env.Name, e.Hof.Env.Kind)
			i, _, err := d.Dir(e, eflags.NoCache)
			if err != nil {
				fmt.Println("error:", err)
				return err
			}

			i, err = i.Sync(ctx)
			if err != nil {
				return err
			}
		}

	}

	return err
}
