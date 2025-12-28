package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"dagger.io/dagger"
	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/env"
	"github.com/hofstadter-io/hof/lib/env/dag"
	"github.com/hofstadter-io/hof/lib/env/incept"
)

func Run(name string, rflags flags.RootPflagpole, eflags flags.EnvPflagpole, cflags flags.Env__RunFlagpole) error {
	R, err := prepRuntime([]string{}, rflags)
	if err != nil {
		return err
	}

	// incept if we are not in dagger
	dst := os.Getenv("DAGGER_SESSION_TOKEN")
	if dst == "" {
		// Run incept
		err := incept.Incept(context.Background(), os.Args, &incept.InceptOptions{
			Progress:    eflags.Progress,
			Interactive: true,
			Stdout:      os.Stdout,
			Stderr:      os.Stderr,
			Stdin:       os.Stdin,
		})
		if err != nil {
			return err
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

	i, err := d.Build(e, eflags.NoCache)
	if err != nil {
		return err
	}

	// this terminal thing is ignoring what the container may have set if not built by dagger in this engine
	// do we have a manually set command?
	var cmd []string
	if cflags.Command != "" {
		cmd = strings.Fields(cflags.Command)
	}
	// default args / cmd?
	if len(cmd) == 0 {
		args, _ := i.DefaultArgs(ctx)
		if len(args) > 0 {
			cmd = args
		}
	}
	// entrypoint?
	if len(cmd) == 0 {
		entry, _ := i.Entrypoint(ctx)
		if len(entry) > 0 {
			cmd = entry
		}
	}

	i, err = i.Terminal(dagger.ContainerTerminalOpts{
		Cmd: cmd,
	}).Sync(ctx)
	if err != nil {
		return err
	}

	return nil
}
